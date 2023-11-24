package rbac

import (
	"crypto/x509/pkix"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCanSubjectGetResource(t *testing.T) {
	kubeconfig, err := os.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating clientConfig")

	restConfig, err := clientConfig.ClientConfig()
	assert.Nil(t, err, "expecting nil error getting restConfig")

	sub := pkix.Name{
		CommonName: "luca", Organization: []string{"devs"},
	}

	gr := schema.GroupResource{
		Group: "widgets.ui.krateo.io", Resource: "cardtemplates",
	}

	ok, err := CanSubjectGetResource(restConfig, sub, gr, "", "dev-system")
	assert.Nil(t, err, "expecting nil error checking if resource name is listable")
	assert.True(t, ok, "expecting resource name listable")
}
