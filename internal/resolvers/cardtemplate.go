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

func CardTemplate(ctx context.Context, rc *rest.Config, ref *core.Reference) (*cardtemplatev1alpha1.CardInfo, error) {
	wcl, err := widgets.NewForConfig(rc)
	if err != nil {
		return nil, err
	}
	res, err := wcl.CardTemplates(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	ds := map[string]any{}
	for _, x := range res.Spec.APIList {
		ep, err := Endpoint(context.TODO(), rc, x.EndpointRef)
		if err != nil {
			return nil, err
		}

		hc, err := api.HTTPClientForEndpoint(ep)
		if err != nil {
			return nil, err
		}

		rt, err := api.Call(context.TODO(), hc, api.CallOptions{
			API:      x,
			Endpoint: ep,
		})
		if err != nil {
			return nil, err
		}

		ds[x.Name] = rt
	}

	tpl, err := tmpl.New("${", "}")
	if err != nil {
		return nil, err
	}

	res.Spec.App.Title, err = tpl.Execute(res.Spec.App.Title, ds)
	if err != nil {
		return nil, err
	}

	res.Spec.App.Content, err = tpl.Execute(res.Spec.App.Content, ds)
	if err != nil {
		return nil, err
	}

	res.Spec.App.Icon, err = tpl.Execute(res.Spec.App.Icon, ds)
	if err != nil {
		return nil, err
	}

	res.Spec.App.Color, err = tpl.Execute(res.Spec.App.Color, ds)
	if err != nil {
		return nil, err
	}

	res.Spec.App.Date, err = tpl.Execute(res.Spec.App.Date, ds)
	if err != nil {
		return nil, err
	}

	res.Spec.App.Tags, err = tpl.Execute(res.Spec.App.Tags, ds)
	if err != nil {
		return nil, err
	}

	return res.Spec.App.DeepCopy(), nil
}
