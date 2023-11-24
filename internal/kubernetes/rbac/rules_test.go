package rbac

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestRulesForRole(t *testing.T) {
	kubeconfig, err := os.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating clientConfig")

	restConfig, err := clientConfig.ClientConfig()
	assert.Nil(t, err, "expecting nil error getting restConfig")

	all, err := RulesForRole(restConfig, &roleInfo{
		kind: "Role", name: "dev", namespace: "dev-system",
	})
	assert.Nil(t, err, "expecting nil error getting role")

	for _, x := range all {
		fmt.Printf("%+v\n", x)
	}
}
