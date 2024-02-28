package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/core/formdefinitions/v1alpha1"
	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/formdefinitions"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

const (
	openAPIV3SchemaFilter = `.spec.versions[] | select(.name="%s") | .schema.openAPIV3Schema`
)

func GetFormSchema(ctx context.Context, rc *rest.Config, in *formtemplatesv1alpha1.FormTemplate) (*runtime.RawExtension, error) {
	formdefinitionsClient, err := formdefinitions.NewClient(rc)
	if err != nil {
		return nil, err
	}

	ref, err := formdefinitionsClient.
		Namespace(in.Spec.DefinitionRef.Namespace).
		Get(ctx, in.Spec.DefinitionRef.Name)
	if err != nil {
		return nil, err
	}

	dyn, err := dynamic.NewGetter(rc)
	if err != nil {
		return nil, err
	}

	crd, err := dyn.Get(ctx, dynamic.GetOptions{
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1",
			Kind:    "CustomResourceDefinition",
		},
		Name: InferGroupResource(ref).String(),
	})
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf(openAPIV3SchemaFilter, ref.Spec.Schema.Version)
	sch, err := dynamic.Extract(ctx, crd, filter)
	if err != nil {
		return nil, err
	}

	return &runtime.RawExtension{Object: &unstructured.Unstructured{
		Object: sch.(map[string]any),
	}}, nil
}

func InferGroupResource(obj *v1alpha1.FormDefinition) schema.GroupResource {
	gk := schema.GroupKind{
		Group: obj.Spec.Schema.Group,
		Kind:  obj.Spec.Schema.Kind,
	}

	kind := types.Type{Name: types.Name{Name: gk.Kind}}
	namer := namer.NewPrivatePluralNamer(nil)
	return schema.GroupResource{
		Group:    gk.Group,
		Resource: strings.ToLower(namer.Name(&kind)),
	}
}
