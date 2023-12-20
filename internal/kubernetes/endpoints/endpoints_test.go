//go:build integration
// +build integration

package endpoints_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestResolve(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	obj, err := endpoints.Resolve(context.TODO(), rc, &core.Reference{
		Name: "httpbin-endpoint", Namespace: "dev-system",
	})
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
