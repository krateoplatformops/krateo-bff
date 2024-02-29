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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "demo-system"
)

func TestFormTemplateGet(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formtemplates.NewClient(cfg)
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

func TestFormTemplateList(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formtemplates.NewClient(cfg)
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

// kubectl get cards plain -n demo-system -o yaml
func TestFormTemplatePlain(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formtemplates.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).Get(context.TODO(), "fireworksapp")
	if err != nil {
		t.Fatal(err)
	}

	//fin, _ := os.Create("ppp.json")
	//defer fin.Close()
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(res)
}

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
