//go:build integration
// +build integration

package cardtemplates_test

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/apis"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestEval(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	dyn, err := dynamic.NewGetter(rc)
	if err != nil {
		t.Fatal(err)
	}

	obj, err := cardtemplates.Get(context.TODO(), dyn, "one", "demo-system")
	if err != nil {
		t.Fatal(err)
	}

	err = cardtemplates.Eval(context.TODO(), obj, cardtemplates.EvalOptions{
		RESTConfig: rc,
		Groups:     []string{"admins"},
	})
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(obj)
}
