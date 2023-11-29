package rbac

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type RbacInterface interface {
	RESTClient() rest.Interface
}

type RbacClient struct {
	restClient rest.Interface
}

func (c *RbacClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}

func (c *RbacClient) ClusterRoles() ClusterRoleInterface {
	return newClusterRoles(c)
}

func (c *RbacClient) ClusterRoleBindings() ClusterRoleBindingInterface {
	return newClusterRoleBindings(c)
}

func (c *RbacClient) Roles(namespace string) RoleInterface {
	return newRoles(c, namespace)
}

func (c *RbacClient) RoleBindings(namespace string) RoleBindingInterface {
	return newRoleBindings(c, namespace)
}

func NewForConfig(rc *rest.Config) (*RbacClient, error) {
	gv := rbacv1.SchemeGroupVersion
	config := *rc
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}

	client, err := rest.RESTClientForConfigAndClient(&config, httpClient)
	if err != nil {
		return nil, err
	}

	return &RbacClient{client}, nil
}
