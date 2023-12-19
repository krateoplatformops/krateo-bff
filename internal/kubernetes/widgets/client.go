package widgets

import (
	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type WidgetsInterface interface {
	RESTClient() rest.Interface
}

type WidgetsClient struct {
	restClient rest.Interface
}

func (c *WidgetsClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}

func (c *WidgetsClient) CardTemplates(ns string) CardTemplateInterface {
	return newCardTemplates(c, ns)
}

func NewForConfig(c *rest.Config) (*WidgetsClient, error) {
	cardtemplatev1alpha1.SchemeBuilder.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{
		Group:   cardtemplatev1alpha1.Group,
		Version: cardtemplatev1alpha1.Version,
	}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &WidgetsClient{client}, nil
}
