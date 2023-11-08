package builtins

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type BuiltinsInterface interface {
	RESTClient() rest.Interface
}

type BuiltinsClient struct {
	restClient rest.Interface
}

func (c *BuiltinsClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}

func (c *BuiltinsClient) Secrets(ns string) SecretInterface {
	return newSecrets(c, ns)
}

func NewForConfig(c *rest.Config) (*BuiltinsClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	config.APIPath = "/api"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &BuiltinsClient{client}, nil
}
