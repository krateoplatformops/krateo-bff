package cardtemplates

import (
	"context"
	"time"

	"github.com/krateoplatformops/krateo-bff/apis"
	"github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

const (
	resourceName = "cardtemplates"
)

func NewClient(rc *rest.Config) (*Client, error) {
	s := runtime.NewScheme()
	apis.AddToScheme(s)

	config := *rc
	config.APIPath = "/apis"
	config.GroupVersion = &schema.GroupVersion{
		Group: v1alpha1.Group, Version: v1alpha1.Version,
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

type Client struct {
	rc rest.Interface
	pc runtime.ParameterCodec
	ns string
}

func (c *Client) Namespace(ns string) *Client {
	c.ns = ns
	return c
}

func (c *Client) Get(ctx context.Context, name string) (result *v1alpha1.CardTemplate, err error) {
	result = &v1alpha1.CardTemplate{}
	err = c.rc.Get().
		Namespace(c.ns).
		Resource(resourceName).
		Name(name).
		Do(ctx).
		Into(result)
	// issue: https://github.com/kubernetes/client-go/issues/541
	result.SetGroupVersionKind(v1alpha1.CardTemplateGroupVersionKind)
	return
}

func (c *Client) List(ctx context.Context, opts metav1.ListOptions) (result *v1alpha1.CardTemplateList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.CardTemplateList{}
	err = c.rc.Get().
		Namespace(c.ns).
		Resource(resourceName).
		VersionedParams(&opts, c.pc).
		Timeout(timeout).
		Do(ctx).
		Into(result)

	result.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("List"))
	return
}

func (c *Client) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.rc.Get().
		Namespace(c.ns).
		Resource(resourceName).
		VersionedParams(&opts, c.pc).
		Timeout(timeout).
		Watch(ctx)
}

func (c *Client) Create(ctx context.Context, obj *v1alpha1.CardTemplate, opts metav1.CreateOptions) (result *v1alpha1.CardTemplate, err error) {
	result = &v1alpha1.CardTemplate{}
	err = c.rc.Post().
		Namespace(c.ns).
		Resource(resourceName).
		VersionedParams(&opts, c.pc).
		Body(obj).
		Do(ctx).
		Into(result)
	// issue: https://github.com/kubernetes/client-go/issues/541
	result.SetGroupVersionKind(v1alpha1.CardTemplateGroupVersionKind)
	return
}

func (c *Client) Update(ctx context.Context, obj *v1alpha1.CardTemplate, opts metav1.UpdateOptions) (result *v1alpha1.CardTemplate, err error) {
	result = &v1alpha1.CardTemplate{}
	err = c.rc.Put().
		Namespace(c.ns).
		Resource(resourceName).
		Name(obj.Name).
		VersionedParams(&opts, c.pc).
		Body(obj).
		Do(ctx).
		Into(result)
	// issue: https://github.com/kubernetes/client-go/issues/541
	result.SetGroupVersionKind(v1alpha1.CardTemplateGroupVersionKind)
	return
}

func (c *Client) UpdateStatus(ctx context.Context, obj *v1alpha1.CardTemplate) (result *v1alpha1.CardTemplate, err error) {
	result = &v1alpha1.CardTemplate{}
	err = c.rc.Put().
		Namespace(c.ns).
		Resource(resourceName).
		Name(obj.Name).
		SubResource("status").
		Body(obj).
		Do(ctx).
		Into(result)
	// issue: https://github.com/kubernetes/client-go/issues/541
	result.SetGroupVersionKind(v1alpha1.CardTemplateGroupVersionKind)
	return
}

func (c *Client) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.rc.Delete().
		Namespace(c.ns).
		Resource(resourceName).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *Client) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.rc.Delete().
		Namespace(c.ns).
		Resource(resourceName).
		VersionedParams(&listOpts, c.pc).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *Client) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1alpha1.CardTemplate, err error) {
	result = &v1alpha1.CardTemplate{}
	err = c.rc.Patch(pt).
		Namespace(c.ns).
		Resource(resourceName).
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, c.pc).
		Body(data).
		Do(ctx).
		Into(result)
	// issue: https://github.com/kubernetes/client-go/issues/541
	result.SetGroupVersionKind(v1alpha1.CardTemplateGroupVersionKind)
	return
}
