package rbac

import (
	"context"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type RoleInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*rbacv1.Role, error)
	List(ctx context.Context, opts metav1.ListOptions) (*rbacv1.RoleList, error)
}

var _ RoleInterface = (*roles)(nil)

type roles struct {
	client rest.Interface
	ns     string
}

func newRoles(c *RbacClient, namespace string) *roles {
	return &roles{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the role, and returns the corresponding role object, and an error if there is any.
func (c *roles) Get(ctx context.Context, name string, options metav1.GetOptions) (result *rbacv1.Role, err error) {
	result = &rbacv1.Role{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("roles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Roles that match those selectors.
func (c *roles) List(ctx context.Context, opts metav1.ListOptions) (result *rbacv1.RoleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &rbacv1.RoleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("roles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}
