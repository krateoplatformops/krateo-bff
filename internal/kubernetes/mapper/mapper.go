package mapper

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func FindGVR(restConfig *rest.Config, gk schema.GroupKind) (schema.GroupVersionResource, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	groupResources, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

	mapping, err := mapper.RESTMapping(gk)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return mapping.Resource, nil
}
