package dynamic_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	cardtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGet(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := dynamic.NewClient(rc)
	if err != nil {
		t.Fatal(err)
	}

	opts := dynamic.Options{
		Namespace: "demo-system",
		GVK:       cardtemplatesv1alpha1.CardTemplateGroupVersionKind,
	}

	uns, err := cli.Get(context.TODO(), "one", opts)
	if err != nil {
		t.Fatal(err)
	}

	res := cardtemplatesv1alpha1.CardTemplate{}
	err = cli.Convert(uns.UnstructuredContent(), &res)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Spec.FormTemplateRef.Namespace) == 0 {
		res.Spec.FormTemplateRef.Namespace = opts.Namespace
	}

	spew.Dump(res)
}

func TestList(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := dynamic.NewClient(rc)
	if err != nil {
		t.Fatal(err)
	}

	opts := dynamic.Options{
		Namespace: "demo-system",
		GVK:       cardtemplatesv1alpha1.CardTemplateGroupVersionKind,
	}

	uns, err := cli.List(context.TODO(), opts)
	if err != nil {
		t.Fatal(err)
	}

	res := cardtemplatesv1alpha1.CardTemplateList{}
	err = cli.Convert(uns.UnstructuredContent(), &res)
	if err != nil {
		t.Fatal(err)
	}

	for i := range res.Items {
		if len(res.Items[i].Spec.FormTemplateRef.Namespace) == 0 {
			res.Items[i].Spec.FormTemplateRef.Namespace = opts.Namespace
		}
	}

	spew.Dump(res)
}

// kubectl apply -f testdata/ns.yaml
// kubectl apply -f testdata/fireworksapp.crd.yaml
func TestApply(t *testing.T) {
	fin, err := os.Open("../../../testdata/fireworksapp.fake.json")
	if err != nil {
		t.Fatal(err)
	}
	defer fin.Close()

	var content map[string]any
	err = json.NewDecoder(fin).Decode(&content)
	if err != nil {
		t.Fatal(err)
	}
	//spew.Dump(content)

	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := dynamic.NewClient(rc)
	if err != nil {
		t.Fatal(err)
	}

	err = cli.Apply(context.TODO(), "demo-apply", map[string]any{"spec": content}, dynamic.Options{
		GVK: schema.GroupVersionKind{
			Group:   "apps.krateo.io",
			Version: "v1alpha1",
			Kind:    "Fireworksapp",
		},
		Namespace: "demo-system",
	})
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
