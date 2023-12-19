package clients

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

type Client struct {
	//Codecs serializer.CodecFactory
	parameterCodec runtime.ParameterCodec
	restClient     *rest.RESTClient
}

func (c *Client) RESTClient() *rest.RESTClient {
	return c.restClient
}

func ForGroupVersion(rc *rest.Config, gv schema.GroupVersion, s *runtime.Scheme) (*Client, error) {
	config := *rc
	config.APIPath = "/apis"
	config.GroupVersion = &gv
	config.NegotiatedSerializer = serializer.NewCodecFactory(s).
		WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	cli, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{
		restClient:     cli,
		parameterCodec: runtime.NewParameterCodec(s),
	}, nil
}
