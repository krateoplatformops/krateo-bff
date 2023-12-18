package resolvers

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type CardTemplateGetAllOpts struct {
	RESTConfig *rest.Config
	Username   string
	AuthnNS    string
	Namespace  string
}

func CardTemplateGetAll(ctx context.Context, opts CardTemplateGetAllOpts) (*cardtemplatev1alpha1.CardTemplateList, error) {
	wcl, err := widgets.NewForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}

	all, err := wcl.CardTemplates(opts.Namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	all.APIVersion = cardtemplatev1alpha1.CardTemplateGroupVersionKind.GroupVersion().String()
	all.Kind = "CardTemplateList"

	for i := range all.Items {
		err := cardTemplateEval(ctx, &all.Items[i], cardTemplateEvalOpts{
			rc: opts.RESTConfig, username: opts.Username, authnNS: opts.AuthnNS,
		})
		if err != nil {
			return all, err
		}
	}

	return all, nil
}

type CardTemplateGetOneOpts struct {
	RESTConfig *rest.Config
	Username   string
	AuthnNS    string
}

func CardTemplateGetOne(ctx context.Context, ref *core.Reference, opts CardTemplateGetOneOpts) (*cardtemplatev1alpha1.CardTemplate, error) {
	wcl, err := widgets.NewForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}
	res, err := wcl.CardTemplates(ref.Namespace).
		Get(ctx, ref.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	res.APIVersion = cardtemplatev1alpha1.CardTemplateGroupVersionKind.GroupVersion().String()
	res.Kind = cardtemplatev1alpha1.CardTemplateGroupVersionKind.Kind

	err = cardTemplateEval(ctx, res, cardTemplateEvalOpts{
		rc: opts.RESTConfig, authnNS: opts.AuthnNS, username: opts.Username,
	})
	return res, err
}

type cardTemplateEvalOpts struct {
	rc       *rest.Config
	authnNS  string
	username string
}

func cardTemplateEval(ctx context.Context, in *cardtemplatev1alpha1.CardTemplate, opts cardTemplateEvalOpts) error {
	ds := map[string]any{}
	for _, x := range in.Spec.APIList {
		ref := x.EndpointRef
		if ptr.Deref(x.KrateoGateway, false) {
			ref = &core.Reference{
				Name:      fmt.Sprintf("%s-kubeconfig", opts.username),
				Namespace: opts.authnNS,
			}
		}

		ep, err := EndpointGetOne(context.TODO(), opts.rc, ref)
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
