package crdutil

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

func InferGroupResource(gk schema.GroupKind) schema.GroupResource {
	kind := types.Type{Name: types.Name{Name: gk.Kind}}
	namer := namer.NewPrivatePluralNamer(nil)
	return schema.GroupResource{
		Group:    gk.Group,
		Resource: strings.ToLower(namer.Name(&kind)),
	}
}
