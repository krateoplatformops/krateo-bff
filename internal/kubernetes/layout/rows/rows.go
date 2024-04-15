package rows

import (
	"context"
	"fmt"
	"os"

	columnsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/apis/ui/rows/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/columns"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func NewClient(rc *rest.Config, verbose bool) (*Client, error) {
	dyn, err := dynamic.NewClient(rc)
	if err != nil {
		return nil, err
	}

	cc, err := columns.NewClient(rc, verbose)
	if err != nil {
		return nil, err
	}

	return &Client{
		dyn: dyn,
		cc:  cc,
		gvk: v1alpha1.RowGroupVersionKind,
	}, nil
}

type Client struct {
	dyn dynamic.Client
	gvk schema.GroupVersionKind
	cc  *columns.Client
}

type GetOptions struct {
	Name      string
	Namespace string
	Subject   string
	Orgs      []string
	AuthnNS   string
}

func (c *Client) Get(ctx context.Context, opts GetOptions) (*v1alpha1.Row, error) {
	uns, err := c.dyn.Get(ctx, opts.Name, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       c.gvk,
	})
	if err != nil {
		return nil, err
	}

	obj := &v1alpha1.Row{}
	err = c.dyn.Convert(uns.UnstructuredContent(), obj)
	if err != nil {
		return nil, err
	}

	err = c.resolveColumns(ctx, obj, resolveOptions{
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

func (c *Client) List(ctx context.Context, opts ListOptions) (*v1alpha1.RowList, error) {
	uns, err := c.dyn.List(ctx, dynamic.Options{
		Namespace: opts.Namespace,
		GVK:       c.gvk,
	})
	if err != nil {
		return nil, err
	}

	all := &v1alpha1.RowList{}
	err = c.dyn.Convert(uns.UnstructuredContent(), all)
	if err != nil {
		return nil, err
	}

	for i := range all.Items {
		err := c.resolveColumns(ctx, &all.Items[i], resolveOptions{
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

func (c *Client) resolveColumns(ctx context.Context, in *v1alpha1.Row, opts resolveOptions) error {
	refs := in.Spec.ColumnListRef
	if refs == nil {
		return nil
	}

	all := &columnsv1alpha1.ColumnList{
		Items: []columnsv1alpha1.Column{},
	}
	all.SetGroupVersionKind(columnsv1alpha1.SchemeGroupVersion.WithKind("ColumnList"))

	for _, ref := range refs {
		el, err := c.cc.Get(ctx, columns.GetOptions{
			Name:      ref.Name,
			Namespace: ref.Namespace,
			AuthnNS:   opts.authnNS,
			Subject:   opts.subject,
			Orgs:      opts.orgs,
		})
		if err != nil {
			if apierrors.IsNotFound(err) {
				fmt.Fprintf(os.Stderr, "WARN: column %q @ %s not found\n", ref.Name, ref.Namespace)
				continue
			}
			return err
		}

		all.Items = append(all.Items, *el)
	}

	in.Status.Content = &runtime.RawExtension{
		Object: all,
	}

	return nil
}
