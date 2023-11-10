package resolvers

import (
	"context"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func CardTemplateGetAll(ctx context.Context, rc *rest.Config, ns string, eval bool) (*cardtemplatev1alpha1.CardTemplateList, error) {
	wcl, err := widgets.NewForConfig(rc)
	if err != nil {
		return nil, err
	}

	all, err := wcl.CardTemplates(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if !eval {
		return all, nil
	}

	for i := range all.Items {
		if err := cardTemplateEval(ctx, rc, &all.Items[i]); err != nil {
			return all, err
		}
	}

	return all, nil
}

func CardTemplateGetOne(ctx context.Context, rc *rest.Config, ref *core.Reference, eval bool) (*cardtemplatev1alpha1.CardTemplate, error) {
	wcl, err := widgets.NewForConfig(rc)
	if err != nil {
		return nil, err
	}
	res, err := wcl.CardTemplates(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if !eval {
		return res, nil
	}

	err = cardTemplateEval(ctx, rc, res)
	return res, err
}

func cardTemplateEval(ctx context.Context, rc *rest.Config, in *cardtemplatev1alpha1.CardTemplate) error {
	ds := map[string]any{}
	for _, x := range in.Spec.APIList {
		ep, err := EndpointGetOne(context.TODO(), rc, x.EndpointRef)
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

	in.Spec.App.Title, err = tpl.Execute(in.Spec.App.Title, ds)
	if err != nil {
		return err
	}

	in.Spec.App.Content, err = tpl.Execute(in.Spec.App.Content, ds)
	if err != nil {
		return err
	}

	in.Spec.App.Icon, err = tpl.Execute(in.Spec.App.Icon, ds)
	if err != nil {
		return err
	}

	in.Spec.App.Color, err = tpl.Execute(in.Spec.App.Color, ds)
	if err != nil {
		return err
	}

	in.Spec.App.Date, err = tpl.Execute(in.Spec.App.Date, ds)
	if err != nil {
		return err
	}

	in.Spec.App.Tags, err = tpl.Execute(in.Spec.App.Tags, ds)
	if err != nil {
		return err
	}

	return nil
}
