//go:build integration
// +build integration

package util

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	sub       = "cyberjoker"
	orgs      = "devs"
	namespace = "dev-system"
)

func TestFilterClusterRoleBindings(t *testing.T) {
	restConfig, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	rbacClient, err := rbac.NewForConfig(restConfig)
	if err != nil {
		t.Fatal(err)
	}

	all, err := filterClusterRoleBindings(context.TODO(), rbacClient, sub, strings.Split(orgs, ","))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n   cluster role bindings: %v", all)
}

func TestFilterRoleBindings(t *testing.T) {
	restConfig, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	rbacClient, err := rbac.NewForConfig(restConfig)
	if err != nil {
		t.Fatal(err)
	}

	all, err := filterRoleBindings(context.TODO(), rbacClient, namespace, sub, strings.Split(orgs, ","))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n   role bindings (namespace: %s): %v", namespace, all)
}

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
