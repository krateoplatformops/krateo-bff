package evaluator

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"

	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Subject    string
	Groups     []string
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
				Name:      fmt.Sprintf("%s-clientconfig", opts.Subject),
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

	in.Status.Cards = make([]*v1alpha1.Card, tot)

	for i := 0; i < tot; i++ {
		nfo, err := renderCard(&in.Spec, tpl, ds, i)
		if err != nil {
			return err
		}

		in.Status.Cards[i] = nfo
	}

	return injectAllowedVerbs(in, allowedVerbsInjectorOptions{
		restConfig: opts.RESTConfig,
		subject:    opts.Subject,
		groups:     opts.Groups,
	})
}

func renderCard(spec *v1alpha1.CardTemplateSpec, tpl tmpl.JQTemplate, ds map[string]any, idx int) (res *v1alpha1.Card, err error) {
	it := ptr.Deref(spec.Iterator, "")

	hackQueryFn := func(q string) string {
		if len(it) == 0 {
			return q
		}

		el := fmt.Sprintf("%s[%d]", it, idx)
		q = strings.Replace(q, "${", fmt.Sprintf("${ %s | ", el), 1)
		return q
	}

	info := spec.CardTemplateInfo

	res = &v1alpha1.Card{}
	res.Title, err = tpl.Execute(hackQueryFn(info.Title), ds)
	if err != nil {
		return
	}

	res.Content, err = tpl.Execute(hackQueryFn(info.Content), ds)
	if err != nil {
		return
	}

	res.Icon, err = tpl.Execute(hackQueryFn(info.Icon), ds)
	if err != nil {
		return
	}

	res.Color, err = tpl.Execute(hackQueryFn(info.Color), ds)
	if err != nil {
		return
	}

	res.Date, err = tpl.Execute(hackQueryFn(info.Date), ds)
	if err != nil {
		return
	}

	res.Tags, err = tpl.Execute(hackQueryFn(info.Tags), ds)

	res.Actions = make([]*core.API, len(info.Actions))
	for i, x := range info.Actions {
		pt := ptr.Deref(x.Path, "")
		if len(pt) > 0 {
			rt, err := tpl.Execute(hackQueryFn(pt), ds)
			if err != nil {
				return nil, err
			}
			pt = rt
		}

		res.Actions[i] = x.DeepCopy()
		res.Actions[i].Path = ptr.To(pt)
	}

	return res, err
}

const (
	allowedVerbsAnnotationKey = "krateo.io/allowed-verbs"
	resource                  = "cardtemplates"
)

type allowedVerbsInjectorOptions struct {
	restConfig *rest.Config
	subject    string
	groups     []string
}

func injectAllowedVerbs(in *v1alpha1.CardTemplate, opts allowedVerbsInjectorOptions) error {
	verbs, err := rbacutil.GetAllowedVerbs(context.TODO(), opts.restConfig, util.ResourceInfo{
		Subject: opts.subject,
		Groups:  opts.groups,
		GroupResource: v1alpha1.CardTemplateGroupVersionKind.GroupVersion().
			WithResource(resource).
			GroupResource(),
		ResourceName: in.GetName(),
		Namespace:    in.GetNamespace(),
	})
	if err != nil {
		return err
	}

	m := in.GetAnnotations()
	if len(m) == 0 {
		m = map[string]string{}
	}
	m[allowedVerbsAnnotationKey] = strings.Join(verbs, ",")
	in.SetAnnotations(m)

	if in.Status.Cards == nil {
		return nil
	}

	for _, el := range in.Status.Cards {
		el.AllowedActions = []string{}
		for _, x := range in.Spec.CardTemplateInfo.Actions {
			verb := strings.ToLower(ptr.Deref(x.Verb, ""))
			if slices.Contains(verbs, verb) {
				el.AllowedActions = append(el.AllowedActions, x.Name)
			}
		}
	}

	return nil
}
