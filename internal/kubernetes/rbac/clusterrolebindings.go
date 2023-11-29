package rbac

import (
	"context"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ClusterRoleBindingInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*rbacv1.ClusterRoleBinding, error)
	List(ctx context.Context, opts metav1.ListOptions) (*rbacv1.ClusterRoleBindingList, error)
}

var _ ClusterRoleBindingInterface = (*clusterRoleBindings)(nil)

type clusterRoleBindings struct {
	client rest.Interface
}

// newClusterRoleBindings returns a ClusterRoleBindings
func newClusterRoleBindings(c *RbacClient) *clusterRoleBindings {
	return &clusterRoleBindings{
		client: c.RESTClient(),
	}
}

func (c *clusterRoleBindings) Get(ctx context.Context, name string, options metav1.GetOptions) (result *rbacv1.ClusterRoleBinding, err error) {
	result = &rbacv1.ClusterRoleBinding{}
	err = c.client.Get().
		Resource("clusterrolebindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterRoleBindings that match those selectors.
func (c *clusterRoleBindings) List(ctx context.Context, opts metav1.ListOptions) (result *rbacv1.ClusterRoleBindingList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &rbacv1.ClusterRoleBindingList{}
	err = c.client.Get().
		Resource("clusterrolebindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}
