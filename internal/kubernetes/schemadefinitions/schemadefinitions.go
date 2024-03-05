package schemadefinitions

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	openAPIV3SchemaFilter = `.spec.versions[] | select(.name="%s") | .schema.openAPIV3Schema`
)

func NewClient(rc *rest.Config) (*Client, error) {
	dyn, err := dynamic.NewGetter(rc)
	if err != nil {
		return nil, err
	}

	return &Client{
		dyn: dyn,
		gvk: schema.GroupVersionKind{
			Group:   "core.krateo.io",
			Version: "v1alpha1",
			Kind:    "SchemaDefinition",
		},
	}, nil
}

type Client struct {
	dyn *dynamic.Getter
	gvk schema.GroupVersionKind
	ns  string
}

func (c *Client) Namespace(ns string) *Client {
	c.ns = ns
	return c
}

func (c *Client) Get(ctx context.Context, name string) (result *unstructured.Unstructured, err error) {
	result, err = c.dyn.Get(ctx, dynamic.GetOptions{
		GVK:       c.gvk,
		Namespace: c.ns,
		Name:      name,
	})

	return
}

func (c *Client) GVK(ctx context.Context, name string) (schema.GroupVersionKind, error) {
	obj, err := c.dyn.Get(ctx, dynamic.GetOptions{
		GVK:       c.gvk,
		Namespace: c.ns,
		Name:      name,
	})
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	data, ok, err := unstructured.NestedStringMap(obj.Object, "spec", "schema")
	if err != nil {
		return schema.GroupVersionKind{}, err
	}
	if !ok {
		return schema.GroupVersionKind{},
			fmt.Errorf("nested map %q not found in '%s/%s'", "spec.schema", c.ns, name)
	}

	kind := data["kind"]
	version := data["version"]
	if len(version) == 0 {
		version = "v1alpha1"
	}

	return schema.GroupVersionKind{
		Group:   "apps.krateo.io",
		Version: version,
		Kind:    kind,
	}, nil
}

func (c *Client) OpenAPISchema(ctx context.Context, gkv schema.GroupVersionKind) (*unstructured.Unstructured, error) {
	crd, err := c.Namespace("").dyn.Get(ctx, dynamic.GetOptions{
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1",
			Kind:    "CustomResourceDefinition",
		},
		Name: dynamic.InferGroupResource(gkv.Group, gkv.Kind).String(),
	})
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf(openAPIV3SchemaFilter, gkv.Version)
	sch, err := dynamic.Extract(ctx, crd, filter)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{
		Object: sch.(map[string]any),
	}, nil
}
