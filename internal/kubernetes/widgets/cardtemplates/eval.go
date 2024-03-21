package cardtemplates

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"

	"github.com/krateoplatformops/krateo-bff/internal/api/batch"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

const (
	actionPathFmt = "/apis/actions?group=%s&version=%s&plural=%s&sub=%s&orgs=%s&name=%s&namespace=%s"
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

	dict, err := batch.Call(ctx, batch.CallOptions{
		RESTConfig: opts.RESTConfig,
		AuthnNS:    opts.AuthnNS,
		Subject:    opts.Subject,
		ApiList:    in.Spec.APIList,
		Tpl:        tpl,
	})
	if err != nil {
		return err
	}

	tot := 1
	it := ptr.Deref(in.Spec.Iterator, "")
	if len(it) > 0 {
		len, err := tpl.Execute(fmt.Sprintf("${ %s | length }", it), dict)
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
		nfo, err := render(&in.Spec, tpl, dict, i)
		if err != nil {
			return err
		}

		in.Status.Cards[i] = nfo
	}

	if in.Status.Actions == nil {
		in.Status.Actions = []*v1alpha1.Action{}
	}

	gv, _ := schema.ParseGroupVersion(in.APIVersion)
	gr := dynamic.InferGroupResource(gv.Group, in.Kind)

	ok, err := rbacutil.CanDeleteResource(ctx, opts.RESTConfig,
		rbacutil.ResourceInfo{
			Subject:       opts.Subject,
			Groups:        opts.Groups,
			GroupResource: gr,
			ResourceName:  in.GetName(),
			Namespace:     in.GetNamespace(),
		})
	if err != nil {
		return err
	}

	if ok {
		in.Status.Actions = append(in.Status.Actions, &v1alpha1.Action{
			Verb: "delete",
			Path: fmt.Sprintf(actionPathFmt,
				gv.Group, gv.Version, gr.Resource,
				opts.Subject, strings.Join(opts.Groups, ","),
				in.GetName(), in.GetNamespace()),
		})
	}

	gr = dynamic.InferGroupResource(formtemplatesv1alpha1.Group, formtemplatesv1alpha1.FormTemplateKind)
	ok, err = rbacutil.CanListResource(ctx, opts.RESTConfig,
		rbacutil.ResourceInfo{
			Subject:       opts.Subject,
			Groups:        opts.Groups,
			GroupResource: gr,
			ResourceName:  in.Spec.FormTemplateRef.Name,
			Namespace:     in.Spec.FormTemplateRef.Namespace,
		})
	if err != nil {
		return err
	}

	if ok {
		in.Status.Actions = append(in.Status.Actions, &v1alpha1.Action{
			Verb: "get",
			Path: fmt.Sprintf(actionPathFmt,
				formtemplatesv1alpha1.Group, formtemplatesv1alpha1.Version, gr.Resource,
				opts.Subject, strings.Join(opts.Groups, ","),
				in.Spec.FormTemplateRef.Name, in.Spec.FormTemplateRef.Namespace),
		})
	}

	return nil
}

func render(spec *v1alpha1.CardTemplateSpec, tpl tmpl.JQTemplate, ds map[string]any, idx int) (res *v1alpha1.Card, err error) {
	it := ptr.Deref(spec.Iterator, "")

	hackQueryFn := func(q string) string {
		if len(it) == 0 {
			return q
		}

		el := fmt.Sprintf("%s[%d]", it, idx)
		q = strings.Replace(q, "${", fmt.Sprintf("${ %s | ", el), 1)
		return q
	}

	info := spec.App

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

	return res, err
}
