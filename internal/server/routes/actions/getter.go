package actions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/rs/zerolog"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type getter struct {
	rc            *rest.Config
	authnNS       string
	kubeServerURL string
	kubeProxyURL  string
}

func (g *getter) Get(req *http.Request) (*unstructured.Unstructured, error) {
	opts := optionsFromRequest(req)

	log := zerolog.Ctx(req.Context()).
		With().
		Str("sub", opts.subject).
		Str("orgs", opts.orgs).
		Logger()

	x := core.API{
		Path: ptr.To(buildApiPath(http.MethodGet, opts)),
		Verb: ptr.To(http.MethodGet),
		EndpointRef: &core.Reference{
			Name:      fmt.Sprintf("%s-clientconfig", opts.subject),
			Namespace: g.authnNS,
		},
	}

	ep, err := endpoints.Resolve(context.Background(), g.rc, x.EndpointRef)
	if err != nil {
		log.Err(err).
			Str("endpoint.name", x.EndpointRef.Name).
			Str("endpoint.namespace", x.EndpointRef.Namespace).
			Msg("resolving endpoint reference")
		return nil, err
	}
	ep.Debug = opts.verbose

	if len(g.kubeProxyURL) > 0 {
		ep.ProxyURL = g.kubeProxyURL
	}
	if len(g.kubeServerURL) > 0 {
		ep.ServerURL = g.kubeServerURL
	}

	hc, err := api.HTTPClientForEndpoint(ep)
	if err != nil {
		log.Err(err).
			Str("endpoint.name", x.EndpointRef.Name).
			Str("endpoint.namespace", x.EndpointRef.Namespace).
			Msg("unable to create HTTP client for endpoint")
		return nil, err
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
		return nil, err
	}

	uns := &unstructured.Unstructured{}
	uns.SetUnstructuredContent(rt)
	return uns, nil
}
