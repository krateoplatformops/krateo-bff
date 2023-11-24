package objects

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
)

func TestListPods(t *testing.T) {
	kubeconfig, err := os.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating clientConfig")

	restConfig, err := clientConfig.ClientConfig()
	assert.Nil(t, err, "expecting nil error getting restConfig")

	resolver, err := NewObjectResolver(restConfig)
	assert.Nil(t, err, "expecting nil error creating object resolver")

	all, err := resolver.List(context.TODO(), schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}, "dev-system")
	assert.Nil(t, err, "expecting nil error listing")

	if all == nil {
		return
	}

	for _, el := range all.Items {
		spew.Dump(el)
	}
}
func TestListObjects(t *testing.T) {
	kubeconfig, err := os.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating clientConfig")

	restConfig, err := clientConfig.ClientConfig()
	assert.Nil(t, err, "expecting nil error getting restConfig")

	resolver, err := NewObjectResolver(restConfig)
	assert.Nil(t, err, "expecting nil error creating object resolver")

	all, err := resolver.List(context.TODO(), schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "RoleBinding",
	}, "")
	assert.Nil(t, err, "expecting nil error listing")

	if all == nil {
		return
	}

	for _, el := range all.Items {
		spew.Dump(el)
		roleKind, _, _ := unstructured.NestedString(el.Object, "roleRef", "kind")
		roleName, _, _ := unstructured.NestedString(el.Object, "roleRef", "name")
		fmt.Println(roleKind, roleName)

		subs, _, _ := unstructured.NestedSlice(el.Object, "subjects")
		for _, x := range subs {
			y, ok := x.(map[string]interface{})
			if !ok {
				continue
			}
			subjectKind, _, _ := unstructured.NestedString(y, "kind")
			subjectName, _, _ := unstructured.NestedString(y, "name")
			fmt.Println(subjectKind, subjectName)
		}
		fmt.Println("------")
		break
	}
}
