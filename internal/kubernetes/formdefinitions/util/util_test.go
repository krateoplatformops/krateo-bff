package util

import (
	"testing"

	"github.com/krateoplatformops/krateo-bff/apis/core/formdefinitions/v1alpha1"
)

func TestInferGroupResource(t *testing.T) {
	table := []struct {
		in   v1alpha1.FormDefinition
		want string
	}{
		{
			in: v1alpha1.FormDefinition{
				Spec: v1alpha1.FormDefinitionSpec{
					Schema: v1alpha1.SchemaInfo{
						Group:   "apps.krateo.io",
						Version: "v1beta1",
						Kind:    "FireworksappForm",
					},
				},
			},
			want: "fireworksappforms.apps.krateo.io",
		},
		{
			in: v1alpha1.FormDefinition{
				Spec: v1alpha1.FormDefinitionSpec{
					Schema: v1alpha1.SchemaInfo{
						Group:   "examples.org",
						Version: "v1",
						Kind:    "MagnificentResource",
					},
				},
			},
			want: "magnificentresources.examples.org",
		},
	}

	for i, tc := range table {
		got := InferGroupResource(&tc.in)
		if got.String() != tc.want {
			t.Fatalf("[tc: %d] got: %v, expected: %v", i, got, tc.want)
		}
	}
}
