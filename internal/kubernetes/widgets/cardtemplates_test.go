//go:build integration
// +build integration

package widgets_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets"
	"github.com/krateoplatformops/krateo-bff/internal/resolvers"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
)

func TestCardTemplateGet(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	uic, err := widgets.NewForConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := uic.CardTemplates("test-system").
		Get(context.TODO(), "card-test", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	ds := map[string]any{}
	for _, x := range res.Spec.APIList {
		ep, err := resolvers.Endpoint(context.TODO(), cfg, x.EndpointRef)
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

	res.Spec.App.Title, err = tpl.Execute(res.Spec.App.Title, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Spec.App.Content, err = tpl.Execute(res.Spec.App.Content, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Spec.App.Icon, err = tpl.Execute(res.Spec.App.Icon, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Spec.App.Color, err = tpl.Execute(res.Spec.App.Color, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Spec.App.Date, err = tpl.Execute(res.Spec.App.Date, ds)
	if err != nil {
		t.Fatal(err)
	}

	res.Spec.App.Tags, err = tpl.Execute(res.Spec.App.Tags, ds)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", res.Spec.App)
}

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
