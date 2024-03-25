package dynamic

import (
	"context"
	"encoding/json"

	"github.com/krateoplatformops/krateo-bff/apis"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/utils/ptr"
)

const (
	PatchedByField = "app.kubernetes.io/patched-by"
	PatchedByValue = "krateo"
)

func NewClient(rc *rest.Config) (Client, error) {
	s := runtime.NewScheme()
	apis.AddToScheme(s)

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

	return &unstructuredClient{
		dynamicClient: dynamicClient,
		mapper:        mapper,
		serializer:    serializer.NewCodecFactory(s).WithoutConversion(),
		converter:     runtime.DefaultUnstructuredConverter,
	}, nil
}

type Options struct {
	Namespace string
	GVK       schema.GroupVersionKind
}

type Client interface {
	Get(ctx context.Context, name string, opts Options) (*unstructured.Unstructured, error)
	List(ctx context.Context, opts Options) (*unstructured.UnstructuredList, error)
	Delete(ctx context.Context, name string, opts Options) error
	Apply(ctx context.Context, name string, content map[string]any, opts Options) error
	Convert(in map[string]any, out any) error
}

var _ Client = (*unstructuredClient)(nil)

type unstructuredClient struct {
	dynamicClient *dynamic.DynamicClient
	mapper        *restmapper.DeferredDiscoveryRESTMapper
	serializer    runtime.NegotiatedSerializer
	converter     runtime.UnstructuredConverter
}

func (uc *unstructuredClient) Get(ctx context.Context, name string, opts Options) (*unstructured.Unstructured, error) {
	ri, err := uc.resourceInterfaceFor(opts)
	if err != nil {
		return nil, err
	}

	return ri.Get(ctx, name, metav1.GetOptions{})
}

func (uc *unstructuredClient) List(ctx context.Context, opts Options) (*unstructured.UnstructuredList, error) {
	ri, err := uc.resourceInterfaceFor(opts)
	if err != nil {
		return nil, err
	}

	return ri.List(ctx, metav1.ListOptions{})
}

func (uc *unstructuredClient) Delete(ctx context.Context, name string, opts Options) error {
	ri, err := uc.resourceInterfaceFor(opts)
	if err != nil {
		return err
	}

	return ri.Delete(ctx, name, metav1.DeleteOptions{})
}

func (uc *unstructuredClient) Apply(ctx context.Context, name string, content map[string]any, opts Options) error {
	if len(content) == 0 {
		return nil
	}

	ri, err := uc.resourceInterfaceFor(opts)
	if err != nil {
		return err
	}

	obj := NewUnstructured(name, content, opts)
	data, err := json.Marshal(&obj)
	if err != nil {
		return err
	}

	// create or Update the object with SSA (types.ApplyPatchType indicates SSA).
	_, err = ri.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: PatchedByValue,
		Force:        ptr.To(true),
	})

	return err
}

func (uc *unstructuredClient) Convert(in map[string]any, out any) error {
	return uc.converter.FromUnstructured(in, out)
}

func (uc *unstructuredClient) resourceInterfaceFor(opts Options) (dynamic.ResourceInterface, error) {
	restMapping, err := uc.mapper.RESTMapping(opts.GVK.GroupKind(), opts.GVK.Version)
	if err != nil {
		return nil, err
	}

	var ri dynamic.ResourceInterface
	if restMapping.Scope.Name() == meta.RESTScopeNameRoot {
		ri = uc.dynamicClient.Resource(restMapping.Resource)
	} else {
		ri = uc.dynamicClient.Resource(restMapping.Resource).
			Namespace(opts.Namespace)
	}
	return ri, nil
}
