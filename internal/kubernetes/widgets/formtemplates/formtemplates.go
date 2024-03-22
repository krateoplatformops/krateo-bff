package formtemplates

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/strvals"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	actionPathFmt         = "/apis/actions?%s"
	openAPIV3SchemaFilter = `.spec.versions[] | select(.name="%s") | .schema.openAPIV3Schema`
)

type ClientOption func(*Client)

func AuthnNS(s string) ClientOption {
	return func(c *Client) {
		c.authnNS = s
	}
}

func Eval(b bool) ClientOption {
	return func(c *Client) {
		c.eval = b
	}
}

type SchemaDefinitionDeref struct {
	group     string
	version   string
	resource  string
	name      string
	namespace string
}

func NewClient(rc *rest.Config, opts ...ClientOption) (*Client, error) {
	dyn, err := dynamic.NewClient(rc)
	if err != nil {
		return nil, err
	}

	c := &Client{
		rc:  rc,
		dyn: dyn,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

type Client struct {
	rc      *rest.Config
	dyn     dynamic.Client
	authnNS string
	eval    bool
}

type GetOptions struct {
	Name      string
	Namespace string
	Subject   string
	Orgs      []string
}

func (c *Client) Get(ctx context.Context, opts GetOptions) (*v1alpha1.FormTemplate, error) {
	uns, err := c.dyn.Get(ctx, opts.Name, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       v1alpha1.FormTemplateGroupVersionKind,
	})
	if err != nil {
		return nil, err
	}

	obj := &v1alpha1.FormTemplate{}
	err = c.dyn.Convert(uns.UnstructuredContent(), obj)
	if err != nil {
		return nil, err
	}

	if len(obj.Spec.SchemaDefinitionRef.Namespace) == 0 {
		obj.Spec.SchemaDefinitionRef.Namespace = opts.Namespace
	}

	if c.eval {
		err = c.doEval(ctx, obj, opts.Subject, opts.Orgs)
	}

	return obj, err
}

type ListOptions struct {
	Namespace string
	Subject   string
	Orgs      []string
}

func (c *Client) List(ctx context.Context, opts ListOptions) (*v1alpha1.FormTemplateList, error) {
	uns, err := c.dyn.List(ctx, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       v1alpha1.FormTemplateGroupVersionKind,
	})
	if err != nil {
		return nil, err
	}

	all := &v1alpha1.FormTemplateList{}
	err = c.dyn.Convert(uns.UnstructuredContent(), all)
	if err != nil {
		return nil, err
	}

	for i := range all.Items {
		if len(all.Items[i].Spec.SchemaDefinitionRef.Namespace) == 0 {
			all.Items[i].Spec.SchemaDefinitionRef.Namespace = opts.Namespace
		}

		if c.eval {
			err = c.doEval(ctx, &all.Items[i], opts.Subject, opts.Orgs)
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
		GVK:       v1alpha1.FormTemplateGroupVersionKind,
	})
}

func (c *Client) doEval(ctx context.Context, in *v1alpha1.FormTemplate, sub string, orgs []string) error {
	ref, err := c.resolveSchemaDefinitionRef(ctx, in)
	if err != nil {
		return err
	}

	err = c.createActions(ctx, in, sub, orgs, ref)
	if err != nil {
		return err
	}

	sch, err := c.openAPISchema(ctx, ref)
	if err != nil {
		return err
	}

	in.Status.Content = &v1alpha1.FormTemplateStatusContent{
		//Instance: &runtime.RawExtension{Object: vals},
		Schema: &runtime.RawExtension{Object: sch},
	}

	return nil
}

func (c *Client) createActions(ctx context.Context, in *v1alpha1.FormTemplate, sub string, orgs []string, ref *SchemaDefinitionDeref) error {
	if in.Status.Actions == nil {
		in.Status.Actions = []*v1alpha1.Action{}
	}

	ok, err := rbacutil.CanCreateOrUpdateResource(ctx, c.rc,
		rbacutil.ResourceInfo{
			Subject:       sub,
			Groups:        orgs,
			GroupResource: schema.ParseGroupResource(fmt.Sprintf("%s.%s", ref.resource, ref.group)),
			ResourceName:  ref.name,
			Namespace:     ref.namespace,
		})
	if err != nil {
		return err
	}
	if ok {
		qs := url.Values{}
		qs.Set("group", ref.group)
		qs.Set("version", ref.version)
		qs.Set("plural", ref.resource)
		qs.Set("sub", sub)
		qs.Set("orgs", strings.Join(orgs, ","))
		qs.Set("name", ref.name)
		qs.Set("namespace", ref.namespace)

		in.Status.Actions = append(in.Status.Actions, &v1alpha1.Action{
			Verb: "create",
			Path: fmt.Sprintf(actionPathFmt, qs.Encode()),
		})
	}

	return nil
}

func (c *Client) resolveSchemaDefinitionRef(ctx context.Context, in *v1alpha1.FormTemplate) (*SchemaDefinitionDeref, error) {
	name := in.Spec.SchemaDefinitionRef.Name
	namespace := in.Spec.SchemaDefinitionRef.Namespace

	uns, err := c.dyn.Get(ctx, name, dynamic.Options{
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

	return &SchemaDefinitionDeref{
		group: gv.Group, version: version, resource: gr.Resource,
		name: name, namespace: namespace,
	}, nil
}

func (c *Client) openAPISchema(ctx context.Context, ref *SchemaDefinitionDeref) (*unstructured.Unstructured, error) {
	gr := schema.GroupResource{Group: ref.group, Resource: ref.resource}

	crd, err := c.dyn.Get(ctx, gr.String(), dynamic.Options{
		Namespace: "",
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1",
			Kind:    "CustomResourceDefinition",
		},
	})
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf(openAPIV3SchemaFilter, ref.version)
	sch, err := dynamic.Extract(ctx, crd, filter)
	if err != nil {
		return nil, err
	}

	dict, ok := sch.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expecting 'map[string]any', got: %T", sch)
	}

	err = injectMetadata(dict)
	return &unstructured.Unstructured{
		Object: sch.(map[string]any),
	}, err
}

func injectMetadata(in map[string]any) error {
	lines := []string{
		"properties.metadata.type=object",
		"properties.metadata.properties.name.type=string",
		"properties.metadata.properties.namespace.type=string",
		"properties.metadata.properties.namespace.type=string",
		"properties.metadata.required={name,namespace}",
	}

	metadata := strings.Join(lines, ",")

	return strvals.ParseInto(metadata, in)
}
