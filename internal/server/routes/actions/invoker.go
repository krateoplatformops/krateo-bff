package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type invoker struct {
	rc            *rest.Config
	authnNS       string
	kubeServerURL string
	kubeProxyURL  string
}

func (inv *invoker) do(req *http.Request) (map[string]any, error) {
	opts := optionsFromRequest(req)

	log := zerolog.Ctx(req.Context()).
		With().
		Str("sub", opts.subject).
		Str("orgs", opts.orgs).
		Logger()

	x := core.API{
		Path: ptr.To(buildApiPath(req.Method, opts)),
		Verb: ptr.To(req.Method),
		EndpointRef: &core.Reference{
			Name:      fmt.Sprintf("%s-clientconfig", opts.subject),
			Namespace: inv.authnNS,
		},
	}

	dat, err := inv.buildPayload(req, opts)
	if err != nil {
		log.Err(err).Msg("building payload for API request")
		return nil, err
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
		return nil, err
	}
	ep.Debug = opts.verbose

	if len(inv.kubeProxyURL) > 0 {
		ep.ProxyURL = inv.kubeProxyURL
	}

	if len(inv.kubeServerURL) > 0 {
		ep.ServerURL = inv.kubeServerURL
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

	return rt, nil
}

func (inv *invoker) buildPayload(req *http.Request, opts options) ([]byte, error) {
	if req.Method == http.MethodGet || req.Method == http.MethodDelete {
		return nil, nil
	}

	buf, err := io.ReadAll(io.LimitReader(req.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, fmt.Errorf("expected payload")
	}

	rv, err := inv.findResourceVersion(req)
	if err != nil {
		return buf, err
	}

	m := map[string]any{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	tmp := unstructured.Unstructured{}
	tmp.SetUnstructuredContent(map[string]any{
		"spec": m,
	})
	tmp.SetAPIVersion(fmt.Sprintf("%s/%s", opts.group, opts.version))
	tmp.SetKind(opts.kind)
	tmp.SetName(opts.name)
	tmp.SetNamespace(opts.namespace)

	if len(rv) > 0 {
		tmp.SetResourceVersion(rv)
	}

	return json.Marshal(tmp.UnstructuredContent())
}

func (inv *invoker) findResourceVersion(req *http.Request) (rv string, err error) {
	if req.Method != http.MethodPut {
		return "", nil
	}

	getter := &getter{
		rc:            inv.rc,
		authnNS:       inv.authnNS,
		kubeServerURL: inv.kubeServerURL,
		kubeProxyURL:  inv.kubeProxyURL,
	}

	uns, err := getter.Get(req)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return "", err
		}
	}
	if uns != nil {
		rv = uns.GetResourceVersion()
	}

	return rv, nil
}
