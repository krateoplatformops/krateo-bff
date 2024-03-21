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

func NewLister(rc *rest.Config) (Lister, error) {
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

	return &dynamicLister{
		dynamicClient: dynamicClient,
		mapper:        mapper,
	}, nil
}

var _ Lister = (*dynamicLister)(nil)

type dynamicLister struct {
	dynamicClient *dynamic.DynamicClient
	mapper        *restmapper.DeferredDiscoveryRESTMapper
}

func (l *dynamicLister) List(ctx context.Context, namespace string, gvk schema.GroupVersionKind) (*unstructured.UnstructuredList, error) {
	restMapping, err := l.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	var ri dynamic.ResourceInterface
	if restMapping.Scope.Name() == meta.RESTScopeNameRoot {
		ri = l.dynamicClient.Resource(restMapping.Resource)
	} else {
		ri = l.dynamicClient.Resource(restMapping.Resource).
			Namespace(namespace)
	}

	return ri.List(ctx, corev1.ListOptions{})
}
