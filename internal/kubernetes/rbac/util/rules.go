package util

import (
	"context"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func namespaceRules(ctx context.Context, rbacClient *rbac.RbacClient, namespace string, roles []string) ([]PolicyRule, error) {
	cli := rbacClient.Roles(namespace)

	res := []PolicyRule{}
	for _, name := range roles {
		obj, err := cli.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return res, err
		}

		for _, rule := range obj.Rules {
			el := &policyRule{
				apiGroups:     make([]string, len(rule.APIGroups)),
				verbs:         make([]string, len(rule.Verbs)),
				resources:     make([]string, len(rule.Resources)),
				resourceNames: make([]string, len(rule.ResourceNames)),
			}
			copy(el.apiGroups, rule.APIGroups)
			copy(el.verbs, rule.Verbs)
			copy(el.resources, rule.Resources)
			copy(el.resourceNames, rule.ResourceNames)

			res = append(res, el)
		}
	}

	return res, nil
}

func clusterRules(ctx context.Context, rbacClient *rbac.RbacClient, roles []string) ([]PolicyRule, error) {
	cli := rbacClient.ClusterRoles()

	res := []PolicyRule{}
	for _, name := range roles {
		obj, err := cli.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return res, err
		}

		for _, rule := range obj.Rules {
			el := &policyRule{
				apiGroups:     make([]string, len(rule.APIGroups)),
				verbs:         make([]string, len(rule.Verbs)),
				resources:     make([]string, len(rule.Resources)),
				resourceNames: make([]string, len(rule.ResourceNames)),
			}
			copy(el.apiGroups, rule.APIGroups)
			copy(el.verbs, rule.Verbs)
			copy(el.resources, rule.Resources)
			copy(el.resourceNames, rule.ResourceNames)

			res = append(res, el)
		}
	}

	return res, nil
}
