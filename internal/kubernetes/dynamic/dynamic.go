package dynamic

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

type Getter interface {
	Get(ctx context.Context, name, namespace string, gvk schema.GroupVersionKind) (*unstructured.Unstructured, error)
}

type Lister interface {
	List(ctx context.Context, namespace string, gvk schema.GroupVersionKind) (*unstructured.UnstructuredList, error)
}

func InferGroupResource(g, k string) schema.GroupResource {
	gk := schema.GroupKind{
		Group: g,
		Kind:  k,
	}

	kind := types.Type{Name: types.Name{Name: gk.Kind}}
	namer := namer.NewPrivatePluralNamer(nil)
	return schema.GroupResource{
		Group:    gk.Group,
		Resource: strings.ToLower(namer.Name(&kind)),
	}
}
