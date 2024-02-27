//go:build integration
// +build integration

package formdefinitions_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/formdefinitions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "demo-system"
)

func TestFormDefinitionGet(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formdefinitions.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).Get(context.TODO(), "fireworksapp")
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		t.Fatal(err)
	}
}

func TestFormDefinitionList(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formdefinitions.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	all, err := cli.Namespace(namespace).List(context.TODO(), v1.ListOptions{})
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
