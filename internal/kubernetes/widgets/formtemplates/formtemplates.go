package formtemplates

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/core"
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

func NewClient(rc *rest.Config, eval bool) (*Client, error) {
	dyn, err := dynamic.NewClient(rc)
	if err != nil {
		return nil, err
	}

	c := &Client{
		rc:   rc,
		dyn:  dyn,
		eval: eval,
	}

	return c, nil
}

type Client struct {
	rc   *rest.Config
	dyn  dynamic.Client
	eval bool
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

	if obj.Spec.SchemaDefinitionRef != nil {
		if len(obj.Spec.SchemaDefinitionRef.Namespace) == 0 {
			obj.Spec.SchemaDefinitionRef.Namespace = opts.Namespace
		}
	}

	if obj.Spec.CompositionDefinitionRef != nil {
		if len(obj.Spec.CompositionDefinitionRef.Namespace) == 0 {
			obj.Spec.CompositionDefinitionRef.Namespace = opts.Namespace
		}
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
		if all.Items[i].Spec.SchemaDefinitionRef != nil {
			if len(all.Items[i].Spec.SchemaDefinitionRef.Namespace) == 0 {
				all.Items[i].Spec.SchemaDefinitionRef.Namespace = opts.Namespace
			}
		}

		if all.Items[i].Spec.CompositionDefinitionRef != nil {
			if len(all.Items[i].Spec.CompositionDefinitionRef.Namespace) == 0 {
				all.Items[i].Spec.CompositionDefinitionRef.Namespace = opts.Namespace
			}
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
	gvk := schemaDefinitionGVK
	ref := in.Spec.SchemaDefinitionRef
	if ref == nil {
		ref = in.Spec.CompositionDefinitionRef
		gvk = compositionDefinitionGVK
	}
	if ref == nil {
		return fmt.Errorf("both 'schemaDefinitionRef' and 'compositionDefinitionRef' are undefined (%s@%s)",
			in.GetName(), in.GetNamespace())
	}

	res, err := c.resolveDefinitionRef(ctx, ref, gvk)
	if err != nil {
		return err
	}

	err = c.createActions(ctx, in, sub, orgs, res)
	if err != nil {
		return err
	}

	sch, err := c.openAPISchema(ctx, res)
	if err != nil {
		return err
	}

	in.Status.Content = &v1alpha1.FormTemplateStatusContent{
		Schema: &runtime.RawExtension{Object: sch},
	}

	return nil
}

func (c *Client) createActions(ctx context.Context, in *v1alpha1.FormTemplate, sub string, orgs []string, ref *DefinitionDeref) error {
	if in.Status.Actions == nil {
		in.Status.Actions = []*v1alpha1.Action{}
	}

	allowed, err := rbacutil.GetAllowedVerbs(ctx, c.rc,
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
	if rbacutil.Can("create", allowed) {
		qs := url.Values{}
		qs.Set("group", ref.group)
		qs.Set("version", ref.version)
		qs.Set("kind", ref.kind)
		qs.Set("plural", ref.resource)
		qs.Set("sub", sub)
		qs.Set("orgs", strings.Join(orgs, ","))
		//qs.Set("name", ref.name)
		//qs.Set("namespace", ref.namespace)

		in.Status.Actions = append(in.Status.Actions, &v1alpha1.Action{
			Verb: "create",
			Path: fmt.Sprintf(actionPathFmt, qs.Encode()),
		})
	}

	return nil
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

func (c *Client) openAPISchema(ctx context.Context, ref *DefinitionDeref) (*unstructured.Unstructured, error) {
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
