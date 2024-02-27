package formtemplates

import (
	"context"

	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/formdefinitions"
	corev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func NewReferenceResolver(rc *rest.Config) (*ReferenceResolver, error) {
	formdefinitionsClient, err := formdefinitions.NewClient(rc)
	if err != nil {
		return nil, err
	}

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

	return &ReferenceResolver{
		formdefinitionsClient: formdefinitionsClient,
		dynamicClient:         dynamicClient,
		mapper:                mapper,
	}, nil
}

type ReferenceResolver struct {
	formdefinitionsClient *formdefinitions.Client
	dynamicClient         *dynamic.DynamicClient
	mapper                *restmapper.DeferredDiscoveryRESTMapper
}

func (rr *ReferenceResolver) Get(ctx context.Context, obj *formtemplatesv1alpha1.FormTemplate) (*unstructured.Unstructured, error) {
	ref, err := rr.formdefinitionsClient.Namespace(obj.Spec.DefinitionRef.Namespace).
		Get(ctx, obj.Spec.DefinitionRef.Name)
	if err != nil {
		return nil, err
	}

	gvk := schema.GroupVersionKind{
		Group:   ref.Spec.Schema.Group,
		Version: ref.Spec.Schema.Version,
		Kind:    ref.Spec.Schema.Kind,
	}

	restMapping, err := rr.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	return rr.dynamicClient.Resource(restMapping.Resource).
		Namespace(obj.Spec.ResourceRef.Namespace).
		Get(ctx, obj.Spec.ResourceRef.Name, corev1.GetOptions{})
}
