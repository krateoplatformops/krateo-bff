//go:build integration
// +build integration

package util

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac"
)

func TestNamespaceRules(t *testing.T) {
	restConfig, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	rbacClient, err := rbac.NewForConfig(restConfig)
	if err != nil {
		t.Fatal(err)
	}

	all, err := namespaceRules(context.TODO(), rbacClient, namespace, []string{"dev"})
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(all)
}

func TestClusterRules(t *testing.T) {
	restConfig, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	rbacClient, err := rbac.NewForConfig(restConfig)
	if err != nil {
		t.Fatal(err)
	}

	all, err := clusterRules(context.TODO(), rbacClient, []string{"widgets-viewer"})
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(all)
}
