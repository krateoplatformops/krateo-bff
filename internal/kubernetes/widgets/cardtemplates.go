package widgets

import (
	"context"
	"time"

	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type CardTemplatesGetter interface {
	CardTemplates(ns string) CardTemplateInterface
}

type CardTemplateInterface interface {
	Get(ctx context.Context, name string, options metav1.GetOptions) (result *cardtemplatev1alpha1.CardTemplate, err error)
	List(ctx context.Context, opts metav1.ListOptions) (result *cardtemplatev1alpha1.CardTemplateList, err error)
}

func newCardTemplates(c *WidgetsClient, ns string) *cardTemplates {
	return &cardTemplates{
		client: c.RESTClient(),
		ns:     ns,
	}
}

type cardTemplates struct {
	client rest.Interface
	ns     string
}

func (rc *cardTemplates) Get(ctx context.Context, name string, options metav1.GetOptions) (result *cardtemplatev1alpha1.CardTemplate, err error) {
	result = &cardtemplatev1alpha1.CardTemplate{}
	err = rc.client.Get().
		Resource("cardtemplates").
		Name(name).Namespace(rc.ns).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

func (rc *cardTemplates) List(ctx context.Context, opts metav1.ListOptions) (result *cardtemplatev1alpha1.CardTemplateList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &cardtemplatev1alpha1.CardTemplateList{}
	err = rc.client.Get().
		Resource("cardtemplates").
		Namespace(rc.ns).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}
