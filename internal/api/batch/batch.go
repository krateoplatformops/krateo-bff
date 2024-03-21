package batch

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type CallOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Subject    string
	Tpl        tmpl.JQTemplate
	ApiList    []*core.API
}

func Call(ctx context.Context, opts CallOptions) (map[string]any, error) {
	if opts.Tpl == nil {
		tpl, err := tmpl.New("${", "}")
		if err != nil {
			return nil, err
		}
		opts.Tpl = tpl
	}

	apiMap := map[string]*core.API{}
	for _, x := range opts.ApiList {
		apiMap[x.Name] = x
	}

	sorted := core.SortApiByDeps(opts.ApiList)

	ds := map[string]any{}
	for _, key := range sorted {
		x, ok := apiMap[key]
		if !ok {
			return nil, fmt.Errorf("API '%s' not found in apiMap", key)
		}

		ref := x.EndpointRef
		if ptr.Deref(x.KrateoGateway, false) {
			ref = &core.Reference{
				Name:      fmt.Sprintf("%s-clientconfig", opts.Subject),
				Namespace: opts.AuthnNS,
			}
		}

		ep, err := endpoints.Resolve(context.TODO(), opts.RESTConfig, ref)
		if err != nil {
			return nil, err
		}

		hc, err := api.HTTPClientForEndpoint(ep)
		if err != nil {
			return nil, err
		}

		rt, err := api.Call(ctx, hc, api.CallOptions{
			API:      x,
			Endpoint: ep,
			Tpl:      opts.Tpl,
			DS:       ds,
		})
		if err != nil {
			return nil, err
		}

		ds[x.Name] = rt
	}

	return ds, nil
}
