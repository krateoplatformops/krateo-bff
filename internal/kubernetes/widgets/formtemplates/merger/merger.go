package merger

import (
	"fmt"
	"strings"

	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/strvals"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Merge(src *formtemplatesv1alpha1.FormTemplate, dst *unstructured.Unstructured) error {
	lines := make([]string, len(src.Spec.Data))
	for i, di := range src.Spec.Data {
		lines[i] = di.String()
	}

	values := strings.Join(lines, ",")
	fmt.Println(values)

	err := strvals.ParseInto(values, dst.UnstructuredContent())
	if err != nil {
		return err
	}

	return nil //unstructured.SetNestedMap(dst.Object, spec, "spec")
}
