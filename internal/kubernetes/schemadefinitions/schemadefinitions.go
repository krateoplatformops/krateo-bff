package schemadefinitions

import (
	"context"
	"fmt"
	"strings"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/strvals"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	openAPIV3SchemaFilter = `.spec.versions[] | select(.name="%s") | .schema.openAPIV3Schema`
)

func NewClient(rc *rest.Config) (*Client, error) {
	dyn, err := dynamic.NewClient(rc)
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
	dyn dynamic.Client
	gvk schema.GroupVersionKind
}

func (c *Client) Get(ctx context.Context, name, namespace string) (result *unstructured.Unstructured, err error) {
	return c.dyn.Get(ctx, name, dynamic.Options{
		Namespace: namespace,
		GVK:       c.gvk,
	})
}

func (c *Client) GVK(ctx context.Context, name, namespace string) (schema.GroupVersionKind, error) {
	obj, err := c.dyn.Get(ctx, name, dynamic.Options{
		Namespace: namespace,
		GVK:       c.gvk,
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
			fmt.Errorf("nested map %q not found in '%s/%s'", "spec.schema", namespace, name)
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
	name := dynamic.InferGroupResource(gkv.Group, gkv.Kind).String()

	crd, err := c.dyn.Get(ctx, name, dynamic.Options{
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

	filter := fmt.Sprintf(openAPIV3SchemaFilter, gkv.Version)
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
