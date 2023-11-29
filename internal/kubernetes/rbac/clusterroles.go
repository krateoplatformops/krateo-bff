package rbac

import (
	"context"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ClusterRoleInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*rbacv1.ClusterRole, error)
	List(ctx context.Context, opts metav1.ListOptions) (*rbacv1.ClusterRoleList, error)
}

var _ ClusterRoleInterface = (*clusterRoles)(nil)

// clusterRoles implements ClusterRoleInterface
type clusterRoles struct {
	client rest.Interface
}

// newClusterRoles returns a ClusterRoles
func newClusterRoles(c *RbacClient) *clusterRoles {
	return &clusterRoles{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterRole, and returns the corresponding clusterRole object, and an error if there is any.
func (c *clusterRoles) Get(ctx context.Context, name string, options metav1.GetOptions) (result *rbacv1.ClusterRole, err error) {
	result = &rbacv1.ClusterRole{}
	err = c.client.Get().
		Resource("clusterroles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterRoles that match those selectors.
func (c *clusterRoles) List(ctx context.Context, opts metav1.ListOptions) (result *rbacv1.ClusterRoleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &rbacv1.ClusterRoleList{}
	err = c.client.Get().
		Resource("clusterroles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}
