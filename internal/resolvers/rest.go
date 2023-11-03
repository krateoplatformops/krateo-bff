package resolvers

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func UnversionedRESTClientFor(rc *rest.Config, gv schema.GroupVersion) (*rest.RESTClient, error) {
	config := *rc
	config.ContentConfig.GroupVersion = &gv
	if len(gv.Group) == 0 {
		config.APIPath = "/api"
	} else {
		config.APIPath = "/apis"
	}
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	return rest.UnversionedRESTClientFor(&config)
}
