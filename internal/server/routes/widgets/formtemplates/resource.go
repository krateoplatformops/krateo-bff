package formtemplates

import (
	"context"
	"fmt"

	formdefinitionsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/core/formdefinitions/v1alpha1"
	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	formdefinitionsutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/formdefinitions/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	openAPIV3SchemaFilter = `.spec.versions[] | select(.name="%s") | .schema.openAPIV3Schema`
)

func getFormValues(ctx context.Context, rc *rest.Config, ref *formdefinitionsv1alpha1.FormDefinition, in *formtemplatesv1alpha1.FormTemplate) (*unstructured.Unstructured, error) {
	dyn, err := dynamic.NewGetter(rc)
	if err != nil {
		return nil, err
	}

	return dyn.Get(ctx, dynamic.GetOptions{
		GVK: schema.GroupVersionKind{
			Group:   ref.Spec.Schema.Group,
			Version: ref.Spec.Schema.Version,
			Kind:    ref.Spec.Schema.Kind,
		},
		Namespace: in.Spec.ResourceRef.Namespace,
		Name:      in.Spec.ResourceRef.Name,
	})
}

func getFormSchema(ctx context.Context, rc *rest.Config, ref *formdefinitionsv1alpha1.FormDefinition) (*unstructured.Unstructured, error) {
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
		Name: formdefinitionsutil.InferGroupResource(ref).String(),
	})
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf(openAPIV3SchemaFilter, ref.Spec.Schema.Version)
	sch, err := dynamic.Extract(ctx, crd, filter)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{
		Object: sch.(map[string]any),
	}, nil
}
