package columns

import (
	"context"

	cardtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	listKind = "CardTemplateList"
)

func NewClient(rc *rest.Config, eval bool) (*Client, error) {
	dyn, err := dynamic.NewClient(rc)
	if err != nil {
		return nil, err
	}

	ctc, err := cardtemplates.NewClient(rc, eval)
	if err != nil {
		return nil, err
	}

	return &Client{
		dyn: dyn,
		ctc: ctc,
		gvk: v1alpha1.ColumnGroupVersionKind,
	}, nil
}

type Client struct {
	dyn  dynamic.Client
	gvk  schema.GroupVersionKind
	ctc  *cardtemplates.Client
	eval bool
}

type GetOptions struct {
	Name      string
	Namespace string
	Subject   string
	Orgs      []string
	AuthnNS   string
}

func (c *Client) Get(ctx context.Context, opts GetOptions) (*v1alpha1.Column, error) {
	uns, err := c.dyn.Get(ctx, opts.Name, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       c.gvk,
	})
	if err != nil {
		return nil, err
	}

	obj := &v1alpha1.Column{}
	err = c.dyn.Convert(uns.UnstructuredContent(), obj)
	if err != nil {
		return nil, err
	}

	err = c.resolveCards(ctx, obj, resolveOptions{
		authnNS: opts.AuthnNS,
		subject: opts.Subject,
		orgs:    opts.Orgs,
	})

	return obj, err
}

type ListOptions struct {
	Namespace string
	Subject   string
	Orgs      []string
	AuthnNS   string
}

func (c *Client) List(ctx context.Context, opts ListOptions) (*v1alpha1.ColumnList, error) {
	uns, err := c.dyn.List(ctx, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       c.gvk,
	})
	if err != nil {
		return nil, err
	}

	all := &v1alpha1.ColumnList{}
	err = c.dyn.Convert(uns.UnstructuredContent(), all)
	if err != nil {
		return nil, err
	}

	for i, _ := range all.Items {
		err := c.resolveCards(ctx, &all.Items[i], resolveOptions{
			authnNS: opts.AuthnNS,
			subject: opts.Subject,
			orgs:    opts.Orgs,
		})
		if err != nil {
			return all, err
		}
	}

	return all, nil
}

type resolveOptions struct {
	authnNS string
	subject string
	orgs    []string
}

func (c *Client) resolveCards(ctx context.Context, in *v1alpha1.Column, opts resolveOptions) error {
	refs := in.Spec.CardTemplateListRef
	if refs == nil {
		return nil
	}

	all := &cardtemplatesv1alpha1.CardTemplateList{
		Items: []cardtemplatesv1alpha1.CardTemplate{},
	}
	all.SetGroupVersionKind(cardtemplatesv1alpha1.SchemeGroupVersion.WithKind(listKind))

	for _, ref := range refs {
		el, err := c.ctc.Get(ctx, cardtemplates.GetOptions{
			Name:      ref.Name,
			Namespace: ref.Namespace,
			AuthnNS:   opts.authnNS,
			Subject:   opts.subject,
			Orgs:      opts.orgs,
		})
		if err != nil {
			return err
		}

		all.Items = append(all.Items, *el)
	}

	in.Status.Content = &runtime.RawExtension{
		Object: all,
	}

	return nil
}
