package util

import (
	"context"
	"slices"
	"strings"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func filterClusterRoleBindings(ctx context.Context, rbacClient *rbac.RbacClient, sub string, groups []string) ([]string, error) {
	acceptFn := func(kind string, name string) bool {
		if strings.EqualFold(kind, "group") {
			return slices.Contains(groups, name)
		}

		if strings.EqualFold(kind, "user") {
			return (sub == name)
		}

		return false
	}

	all, err := rbacClient.ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for _, el := range all.Items {
		if acceptFn == nil {
			continue
		}

		for _, x := range el.Subjects {
			if acceptFn(x.Kind, x.Name) {
				res = append(res, el.Name)
			}
		}
	}

	return res, nil
}

func filterRoleBindings(ctx context.Context, rbacClient *rbac.RbacClient, namespace string, sub string, groups []string) ([]string, error) {
	acceptFn := func(kind string, name string) bool {
		if strings.EqualFold(kind, "group") {
			return slices.Contains(groups, name)
		}

		if strings.EqualFold(kind, "user") {
			return (sub == name)
		}

		return false
	}

	all, err := rbacClient.RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for _, el := range all.Items {
		if acceptFn == nil {
			continue
		}

		for _, x := range el.Subjects {
			if acceptFn(x.Kind, x.Name) {
				res = append(res, el.Name)
			}
		}
	}

	return res, nil
}
