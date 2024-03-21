package dynamic

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	corev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func NewGetter(rc *rest.Config) (Getter, error) {
	dynamicClient, err := dynamic.NewForConfig(rc)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(rc)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(
		cacheddiscovery.NewMemCacheClient(discoveryClient),
	)

	return &dynamicGetter{
		dynamicClient: dynamicClient,
		mapper:        mapper,
	}, nil
}

var _ Getter = (*dynamicGetter)(nil)

type dynamicGetter struct {
	dynamicClient *dynamic.DynamicClient
	mapper        *restmapper.DeferredDiscoveryRESTMapper
}

func (g *dynamicGetter) Get(ctx context.Context, name, namespace string, gvk schema.GroupVersionKind) (*unstructured.Unstructured, error) {
	restMapping, err := g.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	var ri dynamic.ResourceInterface
	if restMapping.Scope.Name() == meta.RESTScopeNameRoot {
		ri = g.dynamicClient.Resource(restMapping.Resource)
	} else {
		ri = g.dynamicClient.Resource(restMapping.Resource).
			Namespace(namespace)
	}

	return ri.Get(ctx, name, corev1.GetOptions{})
}
