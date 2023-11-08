package builtins

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type SecretsGetter interface {
	Secrets(namespace string) SecretInterface
}

type SecretInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.Secret, error)
}

func newSecrets(c *BuiltinsClient, namespace string) *secrets {
	return &secrets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

type secrets struct {
	client rest.Interface
	ns     string
}

func (c *secrets) Get(ctx context.Context, name string, options metav1.GetOptions) (result *corev1.Secret, err error) {
	result = &corev1.Secret{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("secrets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}
