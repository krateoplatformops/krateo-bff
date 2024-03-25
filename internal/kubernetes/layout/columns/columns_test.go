//go:build integration
// +build integration

package columns_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/columns"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestColumnGet(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := columns.NewClient(rc, true)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Get(context.TODO(), columns.GetOptions{
		Name:      "one",
		Namespace: "demo-system",
		Subject:   "cyberjoker",
		Orgs:      []string{"devs"},
	})
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		t.Fatal(err)
	}
}

func TestColumnList(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := columns.NewClient(rc, true)
	if err != nil {
		t.Fatal(err)
	}

	all, err := cli.List(context.TODO(), columns.ListOptions{
		Namespace: "demo-system",
		Subject:   "cyberjoker",
		Orgs:      []string{"devs"},
	})
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(all); err != nil {
		t.Fatal(err)
	}
}

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
