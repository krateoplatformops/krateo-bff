//go:build integration
// +build integration

package api_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/resolvers"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
)

func TestCall(t *testing.T) {
	apiInfo := core.API{
		Name: "test",
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

	authn, err := resolvers.EndpointGetOne(context.TODO(), rc, apiInfo.EndpointRef)
	if err != nil {
		t.Fatal(err)
	}

	httpClient, err := api.HTTPClientForEndpoint(authn)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := api.Call(context.TODO(), httpClient, api.CallOptions{
		API:      &apiInfo,
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
