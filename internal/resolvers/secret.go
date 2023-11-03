package resolvers

import (
	"context"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func GetSecretData(rc *rest.Config, sel *core.Reference) (map[string]string, error) {
	cli, err := UnversionedRESTClientFor(rc, schema.GroupVersion{Group: "", Version: "v1"})
	if err != nil {
		return nil, err
	}

	res := &corev1.Secret{}
	err = cli.Get().Resource("secrets").
		Namespace(sel.Namespace).Name(sel.Name).
		Do(context.Background()).
		Into(res)
	if err != nil {
		return nil, err
	}

	values := map[string]string{}
	for k, v := range res.Data {
		values[k] = string(v)
	}
	return values, nil
}
