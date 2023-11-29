//go:build integration
// +build integration

package util

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGetAllowedVerbs(t *testing.T) {
	restConfig, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	all, err := GetAllowedVerbs(context.TODO(), restConfig, GetAllowedVerbsOption{
		Subject: "cyberjoker",
		Groups:  []string{"devs"},
		GroupResource: schema.GroupResource{
			Group: "widgets.ui.krateo.io", Resource: "cardtemplates",
		},
		ResourceName: "card-dev",
		Namespace:    "dev-system",
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(all)
}

func TestParseGR(t *testing.T) {
	want := schema.GroupResource{
		Group: "widgets.ui.krateo.io", Resource: "cardtemplates",
	}

	s := "cardtemplates.widgets.ui.krateo.io"
	got := schema.ParseGroupResource(s)
	if diff := cmp.Diff(want, got); len(diff) > 0 {
		t.Fatalf("diff: %s", diff)
	}
}
