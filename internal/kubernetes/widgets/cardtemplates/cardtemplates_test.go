//go:build integration
// +build integration

package cardtemplates_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates/evaluator"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "demo-system"
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

	res, err := cli.Namespace(namespace).Get(context.TODO(), "sample")
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

// kubectl get cards plain -n demo-system -o yaml
func TestCardTemplatePlain(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := cardtemplates.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).Get(context.TODO(), "plain")
	if err != nil {
		t.Fatal(err)
	}

	err = evaluator.Eval(context.TODO(), res, evaluator.EvalOptions{
		RESTConfig: cfg,
		AuthnNS:    namespace,
		Username:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err = cli.Namespace(namespace).UpdateStatus(context.TODO(), res)
	if err != nil {
		t.Fatal(err)
	}
}

// kubectl get cards one -n demo-system -o yaml
func TestCardTemplateWithoutIterator(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := cardtemplates.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).Get(context.TODO(), "one")
	if err != nil {
		t.Fatal(err)
	}

	err = evaluator.Eval(context.TODO(), res, evaluator.EvalOptions{
		RESTConfig: cfg,
		AuthnNS:    namespace,
		Username:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err = cli.Namespace(namespace).UpdateStatus(context.TODO(), res)
	if err != nil {
		t.Fatal(err)
	}
}

// kubectl get cards all -n demo-system -o yaml
func TestCardTemplateWithIterator(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := cardtemplates.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).Get(context.TODO(), "all")
	if err != nil {
		t.Fatal(err)
	}

	err = evaluator.Eval(context.TODO(), res, evaluator.EvalOptions{
		RESTConfig: cfg,
		AuthnNS:    namespace,
		Username:   "",
	})
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
