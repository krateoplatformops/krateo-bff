package evaluator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Username   string
}

func Eval(ctx context.Context, in *v1alpha1.Column, opts EvalOptions) error {
	ds := map[string]any{}
	for _, x := range in.Spec.APIList {
		ref := x.EndpointRef
		if ptr.Deref(x.KrateoGateway, false) {
			ref = &core.Reference{
				Name:      fmt.Sprintf("%s-kubeconfig", opts.Username),
				Namespace: opts.AuthnNS,
			}
		}

		ep, err := endpoints.Resolve(context.TODO(), opts.RESTConfig, ref)
		if err != nil {
			return err
		}

		hc, err := api.HTTPClientForEndpoint(ep)
		if err != nil {
			return err
		}

		rt, err := api.Call(ctx, hc, api.CallOptions{
			API:      x,
			Endpoint: ep,
		})
		if err != nil {
			return err
		}

		ds[x.Name] = rt
	}

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(ds); err != nil {
		return err
	}

	in.Status = json.RawMessage(buf.Bytes())

	return nil
}
