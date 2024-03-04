//go:build integration
// +build integration

package schemadefinitions_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/schemadefinitions"
	"github.com/krateoplatformops/krateo-bff/internal/strvals"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "demo-system"
)

func TestGet(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := schemadefinitions.NewClient(cfg)
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

func TestGVK(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := schemadefinitions.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Namespace(namespace).GVK(context.TODO(), "fireworksapp")
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		t.Fatal(err)
	}
}

func TestOpenAPISchema(t *testing.T) {
	cfg, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := schemadefinitions.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	gvk, err := cli.Namespace(namespace).GVK(context.TODO(), "fireworksapp")
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.OpenAPISchema(context.TODO(), gvk)
	if err != nil {
		t.Fatal(err)
	}

	lines := []string{
		"properties.metadata.type=object",
		"properties.metadata.properties.name.type=string",
		"properties.metadata.properties.namespace.type=string",
		"properties.metadata.properties.namespace.type=string",
		"properties.metadata.required={name,namespace}",
	}

	metadata := strings.Join(lines, ",")
	fmt.Println(metadata)

	err = strvals.ParseInto(metadata, res.UnstructuredContent())
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		t.Fatal(err)
	}
}

/*
metadata.type=object
metadata.properties.name.type=string
metadata.properties.namespace.type=string
metadata.properties.namespace.type=string
metadata.required={name,namespace}

	{
	   "metadata":{
	      "type":"object",
	      "properties":{
	         "name":{
	            "type":"string"
	         },
	         "namespace":{
	            "type":"string"
	         }
	      },
	      "required":[
	         "name",
	         "namespace"
	      ]
	   }
	}
*/
func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
