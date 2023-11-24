package rbac

import (
	"crypto/x509/pkix"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestRolesForSubject(t *testing.T) {
	kubeconfig, err := os.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating clientConfig")

	restConfig, err := clientConfig.ClientConfig()
	assert.Nil(t, err, "expecting nil error getting restConfig")

	all, err := RolesForSubject(restConfig, pkix.Name{
		CommonName: "luca", Organization: []string{"devs"},
	}, "dev-system")
	assert.Nil(t, err, "expecting nil error getting roles by subject")

	for _, x := range all {
		fmt.Printf("%+v\n", x)
	}

}
