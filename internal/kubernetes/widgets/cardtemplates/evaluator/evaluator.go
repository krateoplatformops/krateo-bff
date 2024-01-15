package evaluator

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	tpl, err := tmpl.New("${", "}")
	if err != nil {
		return err
	}

	apiMap := map[string]*core.API{}
	for _, x := range in.Spec.APIList {
		apiMap[x.Name] = x
	}

	sorted := core.SortApiByDeps(in.Spec.APIList)

	ds := map[string]any{}
	for _, key := range sorted {
		x, ok := apiMap[key]
		if !ok {
			return fmt.Errorf("API '%s' not found in apiMap", key)
		}

		ref := x.EndpointRef
		if ptr.Deref(x.KrateoGateway, false) {
			ref = &core.Reference{
				Name:      fmt.Sprintf("%s-clientconfig", opts.Username),
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
			Tpl:      tpl,
			DS:       ds,
		})
		if err != nil {
			return err
		}

		ds[x.Name] = rt
	}

	tot := 1
	it := ptr.Deref(in.Spec.Iterator, "")
	if len(it) > 0 {
		len, err := tpl.Execute(fmt.Sprintf("${ %s | length }", it), ds)
		if err != nil {
			return err
		}
		tot, err = strconv.Atoi(len)
		if err != nil {
			return err
		}
	}

	in.Status.Cards = make([]*v1alpha1.CardInfo, tot)

	for i := 0; i < tot; i++ {
		nfo, err := eval(in.Spec, tpl, ds, i)
		if err != nil {
			return err
		}

		in.Status.Cards[i] = nfo
	}

	return nil
}

func eval(spec v1alpha1.CardTemplateSpec, tpl tmpl.JQTemplate, ds map[string]any, idx int) (res *v1alpha1.CardInfo, err error) {
	it := ptr.Deref(spec.Iterator, "")

	hackQueryFn := func(q string) string {

		if len(it) == 0 {
			return q
		}

		el := fmt.Sprintf("%s[%d]", it, idx)
		q = strings.Replace(q, "${", fmt.Sprintf("${ %s | ", el), 1)
		return q
	}

	res = &v1alpha1.CardInfo{}
	res.Title, err = tpl.Execute(hackQueryFn(spec.App.Title), ds)
	if err != nil {
		return
	}

	res.Content, err = tpl.Execute(hackQueryFn(spec.App.Content), ds)
	if err != nil {
		return
	}

	res.Icon, err = tpl.Execute(hackQueryFn(spec.App.Icon), ds)
	if err != nil {
		return
	}

	res.Color, err = tpl.Execute(hackQueryFn(spec.App.Color), ds)
	if err != nil {
		return
	}

	res.Date, err = tpl.Execute(hackQueryFn(spec.App.Date), ds)
	if err != nil {
		return
	}

	res.Tags, err = tpl.Execute(hackQueryFn(spec.App.Tags), ds)

	return res, err
}
