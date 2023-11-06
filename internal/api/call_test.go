//go:build integration
// +build integration

package api

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/resolvers"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
)

func TestCall(t *testing.T) {
	api := core.API{
		Name: ptr.To("test"),
		Path: ptr.To("/anything"),
		Verb: ptr.To("POST"),
		Headers: []string{
			"User-Agent: Krateo",
			"X-Data-1: XXXXXX",
			"X-Data-2: YYYYYY",
		},
		EndpointRef: &core.Reference{
			Name:      "httpbin-endpoint",
			Namespace: "test-system",
		},
		Payload: ptr.To(`{"name": "John", "surname": "Doe"}`),
	}

	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	authn, err := resolvers.GetEndpoint(rc, api.EndpointRef)
	if err != nil {
		t.Fatal(err)
	}

	httpClient, err := HTTPClientForEndpoint(authn)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := Call(context.TODO(), httpClient, CallOptions{
		API:      &api,
		Endpoint: authn,
	})
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(resp)
}

func newRestConfig() (*rest.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
