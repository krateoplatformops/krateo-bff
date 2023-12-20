package evaluator

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Username   string
}

func Eval(ctx context.Context, in *v1alpha1.CardTemplate, opts EvalOptions) error {
	ds := map[string]any{}
	for _, x := range in.Spec.APIList {
		ref := x.EndpointRef
		if ptr.Deref(x.KrateoGateway, false) {
			ref = &core.Reference{
				Name:      fmt.Sprintf("%s-kubeconfig", opts.Username),
				Namespace: opts.AuthnNS,
			}
		}

		ep, err := endpoints.Resolve(context.TODO(), opts.RESTConfig, ref)
		if err != nil {
			return err
		}

		hc, err := api.HTTPClientForEndpoint(ep)
		if err != nil {
			return err
		}

		rt, err := api.Call(ctx, hc, api.CallOptions{
			API:      x,
			Endpoint: ep,
		})
		if err != nil {
			return err
		}

		ds[x.Name] = rt
	}

	tpl, err := tmpl.New("${", "}")
	if err != nil {
		return err
	}

	in.Status.Title, err = tpl.Execute(in.Spec.App.Title, ds)
	if err != nil {
		return err
	}

	in.Status.Content, err = tpl.Execute(in.Spec.App.Content, ds)
	if err != nil {
		return err
	}

	in.Status.Icon, err = tpl.Execute(in.Spec.App.Icon, ds)
	if err != nil {
		return err
	}

	in.Status.Color, err = tpl.Execute(in.Spec.App.Color, ds)
	if err != nil {
		return err
	}

	in.Status.Date, err = tpl.Execute(in.Spec.App.Date, ds)
	if err != nil {
		return err
	}

	in.Status.Tags, err = tpl.Execute(in.Spec.App.Tags, ds)
	if err != nil {
		return err
	}

	return nil
}
