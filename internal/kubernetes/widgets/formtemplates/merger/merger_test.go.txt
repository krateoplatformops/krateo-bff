package merger

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/formtemplates"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestMerge(t *testing.T) {
	ctx := context.TODO()
	namespace := "demo-system"
	name := "fireworksapp"

	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	dyn, err := dynamic.NewGetter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	dst, err := dyn.Get(ctx, dynamic.GetOptions{
		GVK: schema.GroupVersionKind{
			Group:   "apps.krateo.io",
			Version: "v1alpha1",
			Kind:    "FireworksappForm",
		},
		Namespace: namespace,
		Name:      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formtemplates.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	src, err := cli.Namespace(namespace).Get(ctx, name)
	if err != nil {
		t.Fatal(err)
	}

	if err := Merge(src, dst); err != nil {
		t.Fatal(err)
	}

	dump(os.Stdout, dst)
}

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}

func dump(w io.Writer, v any) {
	//fin, _ := os.Create("ppp.json")
	//defer fin.Close()
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
