package cardtemplates

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/core"
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

var (
	schemaDefinitionGVK      = schema.FromAPIVersionAndKind("core.krateo.io/v1alpha1", "SchemaDefinition")
	compositionDefinitionGVK = schema.FromAPIVersionAndKind("core.krateo.io/v1alpha1", "CompositionDefinition")
)

type DefinitionDeref struct {
	group     string
	version   string
	kind      string
	resource  string
	name      string
	namespace string
}

func (dr *DefinitionDeref) String() string {
	sb := strings.Builder{}
	sb.WriteRune('{')
	sb.WriteString("group: ")
	sb.WriteString(dr.group)
	sb.WriteString(", version: ")
	sb.WriteString(dr.version)
	sb.WriteString(", kind: ")
	sb.WriteString(dr.kind)
	sb.WriteString(", resource: ")
	sb.WriteString(dr.resource)
	sb.WriteString(", name: ")
	sb.WriteString(dr.name)
	sb.WriteString(", namespace: ")
	sb.WriteString(dr.namespace)
	sb.WriteRune('}')

	return sb.String()
}

func NewClient(rc *rest.Config, verbose bool) (*Client, error) {
	dyn, err := dynamic.NewClient(rc)
	if err != nil {
		return nil, err
	}

	tpl, err := tmpl.New("${", "}")
	if err != nil {
		return nil, err
	}

	if verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}

	c := &Client{
		rc:      rc,
		dyn:     dyn,
		tpl:     tpl,
		verbose: verbose,
	}

	return c, nil
}

type Client struct {
	rc      *rest.Config
	dyn     dynamic.Client
	tpl     tmpl.JQTemplate
	verbose bool
}

type GetOptions struct {
	Name      string
	Namespace string
	Subject   string
	Orgs      []string
	AuthnNS   string
}

func (c *Client) SetVerbose(v bool) { c.verbose = v }

func (c *Client) Verbose() bool { return c.verbose }

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

	err = c.evalCard(ctx, obj, opts.Subject, opts.Orgs, opts.AuthnNS)

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

		err = c.evalCard(ctx, &all.Items[i], opts.Subject, opts.Orgs, opts.AuthnNS)
		if err != nil {
			return all, err
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

	if c.verbose {
		log.Printf("[DBG] resolved formtemplate reference: %s\n", ref)
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
		qs.Set("version", gvk.Version)
		qs.Set("sub", sub)
		qs.Set("orgs", strings.Join(orgs, ","))
		qs.Set("name", ref.name)
		qs.Set("namespace", ref.namespace)

		actions = append(actions, &v1alpha1.Action{
			Verb: "get",
			Path: fmt.Sprintf("/apis/widgets.ui.krateo.io/formtemplates/%s?%s", ref.name, qs.Encode()),
		})
	}

	return actions, nil
}

func (c *Client) resolveFormTemplateRef(ctx context.Context, in v1alpha1.FormTemplateRef) (*DefinitionDeref, error) {
	uns, err := c.dyn.Get(ctx, in.Name, dynamic.Options{
		Namespace: in.Namespace,
		GVK:       formtemplatesv1alpha1.FormTemplateGroupVersionKind,
	})
	if err != nil {
		return nil, err
	}

	obj := &formtemplatesv1alpha1.FormTemplate{}
	err = c.dyn.Convert(uns.UnstructuredContent(), obj)
	if err != nil {
		return nil, err
	}

	gvk := schemaDefinitionGVK
	ref := obj.Spec.SchemaDefinitionRef
	if ref == nil {
		ref = obj.Spec.CompositionDefinitionRef
		gvk = compositionDefinitionGVK
	}
	if ref == nil {
		return nil,
			fmt.Errorf("both 'schemaDefinitionRef' and 'compositionDefinitionRef' are undefined (%s@%s)", in.Name, in.Namespace)
	}

	if len(ref.Namespace) == 0 {
		ref.Namespace = in.Namespace
	}

	return c.resolveDefinitionRef(ctx, ref, gvk)
}

func (c *Client) resolveDefinitionRef(ctx context.Context, ref *core.Reference, gvk schema.GroupVersionKind) (*DefinitionDeref, error) {
	uns, err := c.dyn.Get(ctx, ref.Name, dynamic.Options{
		Namespace: ref.Namespace,
		GVK:       gvk,
	})
	if err != nil {
		return nil, err
	}

	status, ok, err := unstructured.NestedMap(uns.UnstructuredContent(), "status")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("status not found in '%s/%s'", ref.Namespace, ref.Name)
	}

	apiVersion, ok := status["apiVersion"].(string)
	if !ok {
		return nil, fmt.Errorf("status.apiVersion not found in '%s/%s'", ref.Namespace, ref.Name)
	}

	kind, ok := status["kind"].(string)
	if !ok {
		return nil, fmt.Errorf("status.kind not found in '%s/%s'", ref.Namespace, ref.Name)
	}

	refGVK := schema.FromAPIVersionAndKind(apiVersion, kind)
	refGR := dynamic.InferGroupResource(refGVK.Group, refGVK.Kind)

	return &DefinitionDeref{
		group:     refGVK.Group,
		version:   refGVK.Version,
		resource:  refGR.Resource,
		kind:      refGVK.Kind,
		name:      ref.Name,
		namespace: ref.Namespace,
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
