package crds

import (
	"context"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

func NewClient(rc *rest.Config) (*Client, error) {
	s := runtime.NewScheme()
	apiextensionsv1.AddToScheme(s)

	config := *rc
	config.APIPath = "/apis"
	config.GroupVersion = &schema.GroupVersion{
		Group: apiextensionsv1.SchemeGroupVersion.Group, Version: apiextensionsv1.SchemeGroupVersion.Version,
	}
	config.NegotiatedSerializer = serializer.NewCodecFactory(s).
		WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	cli, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	pc := runtime.NewParameterCodec(s)

	return &Client{rc: cli, pc: pc}, nil
}

const (
	resourceName = "customresourcedefinitions"
)

type Client struct {
	rc rest.Interface
	pc runtime.ParameterCodec
}

func (c *Client) Get(ctx context.Context, name string) (result *apiextensionsv1.CustomResourceDefinition, err error) {
	result = &apiextensionsv1.CustomResourceDefinition{}
	err = c.rc.Get().
		Resource(resourceName).
		Name(name).
		Do(ctx).
		Into(result)
	// issue: https://github.com/kubernetes/client-go/issues/541
	//result.SetGroupVersionKind(apiextensionsv1.SchemeGroupVersion.WithKind())
	return
}
