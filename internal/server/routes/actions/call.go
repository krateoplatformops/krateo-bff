package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/krateoplatformops/krateo-bff/internal/server/decode"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

const (
	callerPath = "/apis/action"
)

func newCaller(rc *rest.Config, authnNS string) (string, http.HandlerFunc) {
	handler := &caller{
		rc:      rc,
		authnNS: authnNS,
	}
	return callerPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*caller)(nil)

type caller struct {
	rc      *rest.Config
	authnNS string
}

func (r *caller) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	qs := req.URL.Query()
	verbose := qs.Has("v")

	sub := qs.Get("sub")
	if len(sub) == 0 {
		encode.BadRequest(wri, fmt.Errorf("missing required request param: 'sub'"))
		return
	}

	log := zerolog.Ctx(req.Context()).With().Logger()

	x := core.API{}
	err := decode.JSONBody(wri, req, &x)
	if err != nil {
		log.Err(err).Msg("decoding JSON API data")
		mr := &decode.MalformedRequest{}
		if errors.As(err, &mr) {
			encode.BadRequest(wri, err)
		} else {
			encode.InternalError(wri, err)
		}
		return
	}

	ref := x.EndpointRef
	if ptr.Deref(x.KrateoGateway, false) {
		ref = &core.Reference{
			Name:      fmt.Sprintf("%s-clientconfig", sub),
			Namespace: r.authnNS,
		}
	}

	ep, err := endpoints.Resolve(context.TODO(), r.rc, ref)
	if err != nil {
		log.Err(err).Msg("resolving endpoint reference")
		encode.InternalError(wri, err)
		return
	}
	ep.Debug = verbose

	hc, err := api.HTTPClientForEndpoint(ep)
	if err != nil {
		log.Err(err).
			Str("endpoint-name", x.EndpointRef.Name).
			Str("endpoint-ref", x.EndpointRef.Namespace).
			Msg("unable to create HTTP client for endpoint")
		encode.InternalError(wri, err)
		return
	}

	rt, err := api.Call(context.TODO(), hc, api.CallOptions{
		API:      &x,
		Endpoint: ep,
	})
	if err != nil {
		log.Err(err).
			Str("endpoint-name", x.EndpointRef.Name).
			Str("endpoint-ref", x.EndpointRef.Namespace).
			Msg("unable to call endpoint")
		encode.InternalError(wri, err)
		return
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rt); err != nil {
		log.Err(err).Msg("unable to serve api call response")
	}
}
