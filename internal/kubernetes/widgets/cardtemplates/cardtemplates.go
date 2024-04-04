package cardtemplates

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/api/batch"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

const (
	actionPathFmt = "/apis/actions?%s"
)

type FormTemplateDeref struct {
	group     string
	version   string
	resource  string
	kind      string
	name      string
	namespace string
}

func NewClient(rc *rest.Config, eval bool) (*Client, error) {
	dyn, err := dynamic.NewClient(rc)
	if err != nil {
		return nil, err
	}

	tpl, err := tmpl.New("${", "}")
	if err != nil {
		return nil, err
	}

	c := &Client{
		rc:   rc,
		dyn:  dyn,
		tpl:  tpl,
		eval: eval,
	}

	return c, nil
}

type Client struct {
	rc   *rest.Config
	dyn  dynamic.Client
	tpl  tmpl.JQTemplate
	eval bool
}

type GetOptions struct {
	Name      string
	Namespace string
	Subject   string
	Orgs      []string
	AuthnNS   string
}

func (c *Client) Get(ctx context.Context, opts GetOptions) (*v1alpha1.CardTemplate, error) {
	uns, err := c.dyn.Get(ctx, opts.Name, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       v1alpha1.CardTemplateGroupVersionKind,
	})
	if err != nil {
		return nil, err
	}

	obj := &v1alpha1.CardTemplate{}
	err = c.dyn.Convert(uns.UnstructuredContent(), obj)
	if err != nil {
		return nil, err
	}

	if len(obj.Spec.FormTemplateRef.Namespace) == 0 {
		obj.Spec.FormTemplateRef.Namespace = opts.Namespace
	}

	if c.eval {
		err = c.evalCard(ctx, obj, opts.Subject, opts.Orgs, opts.AuthnNS)
	}

	return obj, err
}

type ListOptions struct {
	Namespace string
	Subject   string
	AuthnNS   string
	Orgs      []string
}

func (c *Client) List(ctx context.Context, opts ListOptions) (*v1alpha1.CardTemplateList, error) {
	uns, err := c.dyn.List(ctx, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       v1alpha1.CardTemplateGroupVersionKind,
	})
	if err != nil {
		return nil, err
	}

	all := &v1alpha1.CardTemplateList{}
	err = c.dyn.Convert(uns.UnstructuredContent(), all)
	if err != nil {
		return nil, err
	}

	for i := range all.Items {
		if len(all.Items[i].Spec.FormTemplateRef.Namespace) == 0 {
			all.Items[i].Spec.FormTemplateRef.Namespace = opts.Namespace
		}

		if c.eval {
			err = c.evalCard(ctx, &all.Items[i], opts.Subject, opts.Orgs, opts.AuthnNS)
			if err != nil {
				return all, err
			}
		}
	}

	return all, nil
}

type DeleteOptions struct {
	Name      string
	Namespace string
}

func (c *Client) Delete(ctx context.Context, opts DeleteOptions) error {
	return c.dyn.Delete(ctx, opts.Name, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       v1alpha1.CardTemplateGroupVersionKind,
	})
}

func (c *Client) resolveActions(ctx context.Context, in *v1alpha1.CardTemplate, sub string, orgs []string) ([]*v1alpha1.Action, error) {
	actions := []*v1alpha1.Action{}

	ok, err := rbacutil.CanDeleteResource(ctx, c.rc,
		rbacutil.ResourceInfo{
			Subject: sub,
			Groups:  orgs,
			GroupResource: schema.GroupResource{
				Group: v1alpha1.Group, Resource: "cardtemplates",
			},
			ResourceName: in.GetName(),
			Namespace:    in.GetNamespace(),
		})
	if err != nil {
		return nil, err
	}

	if ok {
		qs := url.Values{}
		qs.Set("group", v1alpha1.Group)
		qs.Set("version", "v1alpha1")
		qs.Set("kind", in.Kind)
		qs.Set("plural", "cardtemplates")
		qs.Set("sub", sub)
		qs.Set("orgs", strings.Join(orgs, ","))
		qs.Set("name", in.Name)
		qs.Set("namespace", in.Namespace)

		actions = append(actions, &v1alpha1.Action{
			Verb: "delete",
			Path: fmt.Sprintf(actionPathFmt, qs.Encode()),
		})
	}

	ref, err := c.resolveFormTemplateRef(ctx, in.Spec.FormTemplateRef)
	if err != nil {
		return actions, err
	}

	ok, err = rbacutil.CanListResource(ctx, c.rc,
		rbacutil.ResourceInfo{
			Subject:       sub,
			Groups:        orgs,
			GroupResource: schema.ParseGroupResource(fmt.Sprintf("%s.%s", ref.resource, ref.group)),
			ResourceName:  ref.name,
			Namespace:     ref.namespace,
		})
	if err != nil {
		return actions, err
	}
	if ok {
		gvk := formtemplatesv1alpha1.FormTemplateGroupVersionKind
		qs := url.Values{}
		qs.Set("group", gvk.Group)        //ref.group)
		qs.Set("version", gvk.Version)    // ref.version)
		qs.Set("kind", gvk.Kind)          // ref.kind)
		qs.Set("plural", "formtemplates") // ref.resource)
		qs.Set("sub", sub)
		qs.Set("orgs", strings.Join(orgs, ","))
		qs.Set("name", ref.name)
		qs.Set("namespace", ref.namespace)

		actions = append(actions, &v1alpha1.Action{
			Verb: "get",
			Path: fmt.Sprintf(actionPathFmt, qs.Encode()),
		})
	}

	return actions, nil
}

func (c *Client) resolveFormTemplateRef(ctx context.Context, in v1alpha1.FormTemplateRef) (*FormTemplateDeref, error) {
	uns, err := c.dyn.Get(ctx, in.Name, dynamic.Options{
		Namespace: in.Namespace,
		GVK:       formtemplatesv1alpha1.FormTemplateGroupVersionKind,
	})
	if err != nil {
		return nil, err
	}

	schemaDefinitionRef, ok, _ := unstructured.NestedMap(uns.UnstructuredContent(), "spec", "schemaDefinitionRef")
	if !ok {
		return nil, fmt.Errorf("unable to resolve 'schemaDefinitionRef'")
	}

	name, ok, _ := unstructured.NestedString(schemaDefinitionRef, "name")
	if !ok {
		return nil, fmt.Errorf("unable to resolve 'schemaDefinitionRef.name'")
	}
	namespace, ok, _ := unstructured.NestedString(schemaDefinitionRef, "namespace")
	if !ok {
		return nil, fmt.Errorf("unable to resolve 'schemaDefinitionRef.namespace'")
	}

	uns, err = c.dyn.Get(ctx, name, dynamic.Options{
		Namespace: namespace,
		GVK:       schema.FromAPIVersionAndKind("core.krateo.io/v1alpha1", "SchemaDefinition"),
	})
	if err != nil {
		return nil, err
	}

	schemaSpec, ok, _ := unstructured.NestedMap(uns.UnstructuredContent(), "spec", "schema")
	if !ok {
		return nil,
			fmt.Errorf("unable to resolve 'schema.spec' for: %s.%s/%s @ %s",
				uns.GetName(), uns.GetAPIVersion(), uns.GetKind(), uns.GetNamespace())
	}

	version, ok, _ := unstructured.NestedString(schemaSpec, "version")
	if !ok {
		return nil,
			fmt.Errorf("unable to resolve 'schema.spec.version' for: %s.%s/%s @ %s",
				uns.GetName(), uns.GetAPIVersion(), uns.GetKind(), uns.GetNamespace())
	}
	kind, ok, _ := unstructured.NestedString(schemaSpec, "kind")
	if !ok {
		return nil,
			fmt.Errorf("unable to resolve 'schema.spec.kind' for: %s.%s/%s @ %s",
				uns.GetName(), uns.GetAPIVersion(), uns.GetKind(), uns.GetNamespace())
	}

	gv := schema.GroupVersion{Group: "apps.krateo.io", Version: version}
	gr := dynamic.InferGroupResource(gv.Group, kind)

	return &FormTemplateDeref{
		group: gv.Group, version: version, resource: gr.Resource, kind: kind,
		name: name, namespace: namespace,
	}, nil
}

func (c *Client) evalCard(ctx context.Context, in *v1alpha1.CardTemplate, sub string, orgs []string, authnNS string) error {
	cards, err := c.evalTemplate(ctx, in, sub, authnNS)
	if err != nil {
		return err
	}

	actions, err := c.resolveActions(ctx, in, sub, orgs)
	if err != nil {
		return err
	}

	if len(actions) > 0 {
		for i := range cards {
			cards[i].Actions = make([]*v1alpha1.Action, len(actions))
			copy(cards[i].Actions, actions)
		}
	}

	in.Status.Cards = cards
	return nil
}

func (c *Client) evalTemplate(ctx context.Context, in *v1alpha1.CardTemplate, sub string, authnNS string) ([]*v1alpha1.EvalCard, error) {
	dict, err := batch.Call(ctx, batch.CallOptions{
		RESTConfig: c.rc,
		Tpl:        c.tpl,
		ApiList:    in.Spec.APIList,
		AuthnNS:    authnNS,
		Subject:    sub,
	})
	if err != nil {
		return nil, err
	}

	tot := 1
	it := ptr.Deref(in.Spec.Iterator, "")
	if len(it) > 0 {
		len, err := c.tpl.Execute(fmt.Sprintf("${ %s | length }", it), dict)
		if err != nil {
			return nil, err
		}
		tot, err = strconv.Atoi(len)
		if err != nil {
			return nil, err
		}
	}

	cards := make([]*v1alpha1.EvalCard, tot)

	for i := 0; i < tot; i++ {
		cards[i] = c.renderCard(&in.Spec, dict, i)
	}

	return cards, nil
}

func (c *Client) renderCard(spec *v1alpha1.CardTemplateSpec, ds map[string]any, idx int) *v1alpha1.EvalCard {
	it := ptr.Deref(spec.Iterator, "")

	hackQueryFn := func(q string) string {
		if len(it) == 0 {
			return q
		}

		el := fmt.Sprintf("%s[%d]", it, idx)
		q = strings.Replace(q, "${", fmt.Sprintf("${ %s | ", el), 1)
		return q
	}

	render := func(s string, ds map[string]any) string {
		out, err := c.tpl.Execute(hackQueryFn(s), ds)
		if err != nil {
			out = err.Error()
		}
		return out
	}

	return &v1alpha1.EvalCard{
		Card: v1alpha1.Card{
			Title:   render(spec.App.Title, ds),
			Content: render(spec.App.Content, ds),
			Icon:    render(spec.App.Icon, ds),
			Color:   render(spec.App.Color, ds),
			Date:    render(spec.App.Date, ds),
			Tags:    render(spec.App.Tags, ds),
		},
	}
}
