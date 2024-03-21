//go:build integration
// +build integration

package cardtemplates_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/apis"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGet(t *testing.T) {
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

	spew.Dump(obj)
}

func TestList(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	dyn, err := dynamic.NewLister(rc)
	if err != nil {
		t.Fatal(err)
	}

	obj, err := cardtemplates.List(context.TODO(), dyn, "demo-system")
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(obj)
}

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
