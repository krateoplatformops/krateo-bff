package rbac

import (
	"context"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type RoleBindingInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*rbacv1.RoleBinding, error)
	List(ctx context.Context, opts metav1.ListOptions) (*rbacv1.RoleBindingList, error)
}

var _ RoleBindingInterface = (*roleBindings)(nil)

type roleBindings struct {
	client rest.Interface
	ns     string
}

// newRoleBindings returns a RoleBindings
func newRoleBindings(c *RbacClient, namespace string) *roleBindings {
	return &roleBindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the roleBinding, and returns the corresponding roleBinding object, and an error if there is any.
func (c *roleBindings) Get(ctx context.Context, name string, options metav1.GetOptions) (result *rbacv1.RoleBinding, err error) {
	result = &rbacv1.RoleBinding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rolebindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RoleBindings that match those selectors.
func (c *roleBindings) List(ctx context.Context, opts metav1.ListOptions) (result *rbacv1.RoleBindingList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &rbacv1.RoleBindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rolebindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}
