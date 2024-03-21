package cardtemplates

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"k8s.io/apimachinery/pkg/runtime"
)

func Get(ctx context.Context, dyn dynamic.Getter, name, namespace string) (*v1alpha1.CardTemplate, error) {
	uns, err := dyn.Get(ctx, name, namespace, v1alpha1.CardTemplateGroupVersionKind)
	if err != nil {
		return nil, err
	}

	obj := &v1alpha1.CardTemplate{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(uns.UnstructuredContent(), obj)
	if err == nil {
		if len(obj.Spec.FormTemplateRef.Namespace) == 0 {
			obj.Spec.FormTemplateRef.Namespace = namespace
		}
	}
	return obj, err
}

func List(ctx context.Context, dyn dynamic.Lister, namespace string) (*v1alpha1.CardTemplateList, error) {
	uns, err := dyn.List(ctx, namespace, v1alpha1.CardTemplateGroupVersionKind)
	if err != nil {
		return nil, err
	}

	spew.Dump(uns)
	all := &v1alpha1.CardTemplateList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(uns.UnstructuredContent(), all)
	if err == nil {
		for _, el := range all.Items {
			if len(el.Spec.FormTemplateRef.Namespace) == 0 {
				el.Spec.FormTemplateRef.Namespace = namespace
			}
		}
	}
	return all, err
}
