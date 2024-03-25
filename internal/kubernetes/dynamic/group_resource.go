package dynamic

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

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
