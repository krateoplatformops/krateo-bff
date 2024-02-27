package crdutil

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestInferGroupResource(t *testing.T) {
	table := []struct {
		group string
		kind  string
		exp   string
	}{
		{
			group: "apps.krateo.io",
			kind:  "FireworksappForm",
			exp:   "fireworksappforms.apps.krateo.io",
		},
	}

	for i, tc := range table {
		got := InferGroupResource(schema.GroupKind{
			Group: tc.group,
			Kind:  tc.kind,
		})
		if got.String() != tc.exp {
			t.Fatalf("[tc: %d] - got: %v, expected: %v", i, got, tc.exp)
		}
	}
}
