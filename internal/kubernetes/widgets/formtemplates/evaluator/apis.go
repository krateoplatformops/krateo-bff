package evaluator

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

type callAPIsOptions struct {
	restConfig *rest.Config
	authnNS    string
	subject    string
	tpl        tmpl.JQTemplate
	apiList    []*core.API
}

func callAPIs(ctx context.Context, opts callAPIsOptions) (map[string]any, error) {
	apiMap := map[string]*core.API{}
	for _, x := range opts.apiList {
		apiMap[x.Name] = x
	}

	sorted := core.SortApiByDeps(opts.apiList)

	ds := map[string]any{}
	for _, key := range sorted {
		x, ok := apiMap[key]
		if !ok {
			return nil, fmt.Errorf("API '%s' not found in apiMap", key)
		}

		ref := x.EndpointRef
		if ptr.Deref(x.KrateoGateway, false) {
			ref = &core.Reference{
				Name:      fmt.Sprintf("%s-clientconfig", opts.subject),
				Namespace: opts.authnNS,
			}
		}

		ep, err := endpoints.Resolve(context.TODO(), opts.restConfig, ref)
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
			Tpl:      opts.tpl,
			DS:       ds,
		})
		if err != nil {
			return nil, err
		}

		ds[x.Name] = rt
	}

	return ds, nil
}
