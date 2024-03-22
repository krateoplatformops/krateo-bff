package dynamic

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func NewUnstructured(name string, content map[string]any, opts Options) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetUnstructuredContent(content)
	obj.SetGroupVersionKind(opts.GVK)
	obj.SetNamespace(opts.Namespace)
	obj.SetName(name)
	return obj
}
