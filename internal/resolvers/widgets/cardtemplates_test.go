package widgets

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/resolvers"
)

func TestCardTemplateGet(t *testing.T) {
	all := scheme.Scheme.KnownTypes(cardtemplatev1alpha1.SchemeGroupVersion)
	if len(all) == 0 {
		cardtemplatev1alpha1.SchemeBuilder.AddToScheme(scheme.Scheme)
	}

	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	uic, err := NewForConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := uic.CardTemplates("test-system").
		Get(context.TODO(), "card-test", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	ds := map[string]map[string]any{}
	for _, x := range res.Spec.APIList {
		ep, err := resolvers.GetEndpoint(cfg, x.EndpointRef)
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

	spew.Dump(ds)
}

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
