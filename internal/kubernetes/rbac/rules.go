package rbac

import (
	"context"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/objects"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

func RulesForRole(restConfig *rest.Config, role RoleInfo) ([]PolicyRule, error) {
	resolver, err := objects.NewObjectResolver(restConfig)
	if err != nil {
		return nil, err
	}

	ref := corev1.ObjectReference{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       role.Kind(),
		Name:       role.Name(),
		Namespace:  role.Namespace(),
	}

	obj, err := resolver.ResolveReference(context.TODO(), &ref)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, err
	}

	return extractRules(obj), nil
}

func extractRules(obj *unstructured.Unstructured) []PolicyRule {
	all := []PolicyRule{}

	rules, _, _ := unstructured.NestedSlice(obj.Object, "rules")
	for _, x := range rules {
		y, ok := x.(map[string]interface{})
		if !ok {
			continue
		}

		pol := &policyRule{}
		pol.resources, ok, _ = unstructured.NestedStringSlice(y, "resources")
		if !ok {
			continue
		}

		pol.resourceNames, _, _ = unstructured.NestedStringSlice(y, "resourceNames")

		pol.verbs, ok, _ = unstructured.NestedStringSlice(y, "verbs")
		if !ok {
			continue
		}

		pol.apiGroups, ok, _ = unstructured.NestedStringSlice(y, "apiGroups")
		if !ok {
			continue
		}

		all = append(all, pol)
	}

	return all
}
