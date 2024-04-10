//go:build integration
// +build integration

package formtemplates_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/formtemplates"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGet(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formtemplates.NewClient(rc, true)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Get(context.TODO(), formtemplates.GetOptions{
		Name:      "fireworksapp",
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

func TestList(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formtemplates.NewClient(rc, true)
	if err != nil {
		t.Fatal(err)
	}

	all, err := cli.List(context.TODO(), formtemplates.ListOptions{
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
