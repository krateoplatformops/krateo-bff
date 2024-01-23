//go:build integration
// +build integration

package evaluator_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/rows"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/rows/evaluator"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "demo-system"
	sub       = "cyberjoker"
)

func TestRowEval(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := rows.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).Get(context.TODO(), "two")
	if err != nil {
		t.Fatal(err)
	}

	err = evaluator.Eval(context.TODO(), res, evaluator.EvalOptions{
		RESTConfig: cfg,
		AuthnNS:    namespace,
		Subject:    sub,
		Groups:     []string{"devs"},
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

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
