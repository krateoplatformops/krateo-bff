//go:build integration
// +build integration

package cardtemplates_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "dev-system"
	name      = "card-dev"
)

func TestCardTemplateGet(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := cardtemplates.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).Get(context.TODO(), name)
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		t.Fatal(err)
	}
}

func TestCardTemplateList(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := cardtemplates.NewClient(cfg)
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

func TestCardTemplateUpdateStatus(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := cardtemplates.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).Get(context.TODO(), name)
	if err != nil {
		t.Fatal(err)
	}

	ds := map[string]any{}
	for _, x := range res.Spec.APIList {
		ep, err := endpoints.Resolve(context.TODO(), cfg, x.EndpointRef)
		if err != nil {
			t.Fatal(err)
		}

		hc, err := api.HTTPClientForEndpoint(ep)
		if err != nil {
			t.Fatal(err)
		}

		rt, err := api.Call(context.TODO(), hc, api.CallOptions{
			API:      x,
			Endpoint: ep,
		})
		if err != nil {
			t.Fatal(err)
		}

		ds[x.Name] = rt
	}

	tpl, err := tmpl.New("${", "}")
	if err != nil {
		t.Fatal(err)
	}

	res.Status.Title, err = tpl.Execute(res.Spec.App.Title, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Status.Content, err = tpl.Execute(res.Spec.App.Content, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Status.Icon, err = tpl.Execute(res.Spec.App.Icon, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Status.Color, err = tpl.Execute(res.Spec.App.Color, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Status.Date, err = tpl.Execute(res.Spec.App.Date, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Status.Tags, err = tpl.Execute(res.Spec.App.Tags, ds)
	if err != nil {
		t.Fatal(err)
	}

	res, err = cli.Namespace(namespace).UpdateStatus(context.TODO(), res)
	if err != nil {
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
