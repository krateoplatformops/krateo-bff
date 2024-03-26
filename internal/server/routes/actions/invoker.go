package actions

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/rs/zerolog"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type invokeOptions struct {
	name      string
	namespace string
	subject   string
	orgs      string
	plural    string
	group     string
	version   string
	verbose   bool
}

type invoker struct {
	rc      *rest.Config
	authnNS string
}

func (inv *invoker) do(req *http.Request) map[string]any {
	opts := invokeOptionsFromRequest(req)

	log := zerolog.Ctx(req.Context()).
		With().
		Str("sub", opts.subject).
		Str("orgs", opts.orgs).
		Logger()

	x := core.API{
		Path: ptr.To(
			path.Join("/apis", opts.group, opts.version,
				"namespaces", opts.namespace,
				opts.plural, opts.name),
		),
		Verb: ptr.To(req.Method),
		EndpointRef: &core.Reference{
			Name:      fmt.Sprintf("%s-clientconfig", opts.subject),
			Namespace: inv.authnNS,
		},
	}

	var dat []byte
	switch req.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		var err error
		dat, err = io.ReadAll(io.LimitReader(req.Body, 1<<20))
		if err != nil {
			log.Warn().Msg(err.Error())
		}
	default:
		dat = nil
	}

	if dat != nil {
		x.Payload = ptr.To(string(dat))
	}

	ep, err := endpoints.Resolve(context.Background(), inv.rc, x.EndpointRef)
	if err != nil {
		log.Err(err).
			Str("endpoint.name", x.EndpointRef.Name).
			Str("endpoint.namespace", x.EndpointRef.Namespace).
			Msg("resolving endpoint reference")
		return nil
	}
	ep.Debug = opts.verbose

	hc, err := api.HTTPClientForEndpoint(ep)
	if err != nil {
		log.Err(err).
			Str("endpoint.name", x.EndpointRef.Name).
			Str("endpoint.namespace", x.EndpointRef.Namespace).
			Msg("unable to create HTTP client for endpoint")
		return nil
	}

	rt, err := api.Call(context.Background(), hc, api.CallOptions{
		API:      &x,
		Endpoint: ep,
	})
	if err != nil {
		log.Err(err).
			Str("endpoint.name", x.EndpointRef.Name).
			Str("endpoint.namespace", x.EndpointRef.Namespace).
			Msg("unable to call endpoint")
		return nil
	}

	return rt
}

func invokeOptionsFromRequest(req *http.Request) (opts invokeOptions) {
	qs := req.URL.Query()

	opts.verbose, _ = strconv.ParseBool(qs.Get("verbose"))
	opts.version = qs.Get("version")
	opts.group = qs.Get("group")
	opts.plural = qs.Get("plural")
	opts.name = qs.Get("name")
	opts.namespace = qs.Get("namespace")
	opts.subject = qs.Get("sub")
	opts.orgs = qs.Get("orgs")

	return
}
