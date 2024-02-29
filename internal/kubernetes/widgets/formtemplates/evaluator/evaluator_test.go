//go:build integration
// +build integration

package evaluator

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

func TestEval(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := formtemplates.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	namespace := "demo-system"
	name := "fireworksapp-with-api"

	res, err := cli.Namespace(namespace).Get(context.TODO(), name)
	if err != nil {
		t.Fatal(err)
	}

	err = Eval(context.TODO(), res, EvalOptions{
		RESTConfig: cfg,
		AuthnNS:    namespace,
		Subject:    "",
	})
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
