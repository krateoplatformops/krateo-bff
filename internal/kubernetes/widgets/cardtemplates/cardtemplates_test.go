//go:build integration
// +build integration

package cardtemplates_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGet(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := cardtemplates.NewClient(rc, true)
	if err != nil {
		t.Fatal(err)
	}

	obj, err := cli.Get(context.TODO(), cardtemplates.GetOptions{
		Namespace: "demo-system",
		Name:      "one",
		Subject:   "cyberjoker",
		Orgs:      []string{"devs"},
		AuthnNS:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(obj)
}

func TestList(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := cardtemplates.NewClient(rc, true)
	if err != nil {
		t.Fatal(err)
	}

	obj, err := cli.List(context.TODO(), cardtemplates.ListOptions{
		Namespace: "demo-system",
		Subject:   "cyberjoker",
		Orgs:      []string{"devs"},
		AuthnNS:   "",
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
