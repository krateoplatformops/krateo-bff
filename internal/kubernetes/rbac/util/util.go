package util

import (
	"context"
	"slices"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type GetAllowedVerbsOption struct {
	Subject       string
	Groups        []string
	GroupResource schema.GroupResource
	ResourceName  string
	Namespace     string
}

func GetAllowedVerbs(ctx context.Context, restConfig *rest.Config, opts GetAllowedVerbsOption) ([]string, error) {
	rbacClient, err := rbac.NewForConfig(restConfig)
	if err != nil {
		return []string{}, err
	}

	allow := []string{}

	if len(opts.Namespace) > 0 {
		roles, err := filterRoleBindings(ctx, rbacClient, opts.Namespace, opts.Subject, opts.Groups)
		if err != nil {
			return nil, err
		}

		rules, err := namespaceRules(ctx, rbacClient, opts.Namespace, roles)
		if err != nil {
			return nil, err
		}

		all := allowedVerbs(rules, opts.GroupResource, opts.ResourceName)
		if len(all) > 0 {
			allow = append(allow, all...)
		}
	}

	roles, err := filterClusterRoleBindings(ctx, rbacClient, opts.Subject, opts.Groups)
	if err != nil {
		return nil, err
	}

	rules, err := clusterRules(ctx, rbacClient, roles)
	if err != nil {
		return nil, err
	}

	all := allowedVerbs(rules, opts.GroupResource, opts.ResourceName)
	if len(all) > 0 {
		allow = append(allow, all...)
	}

	return slices.Compact(allow), nil
}

func allowedVerbs(rules []PolicyRule, gr schema.GroupResource, resourceName string) []string {
	checkResourceName := len(resourceName) > 0

	result := []string{}
	for _, x := range rules {
		if !slices.Contains(x.APIGroups(), gr.Group) {
			continue
		}

		if checkResourceName {
			if slices.Contains(x.ResourceNames(), resourceName) {
				result = append(result, x.Verbs()...)
				continue
			}
		}

		if len(x.ResourceNames()) > 0 {
			continue
		}

		ok := slices.Contains(x.Resources(), "*")
		if !ok {
			ok = slices.Contains(x.Resources(), gr.Resource)
			ok = ok && !checkResourceName
		}
		if !ok {
			continue
		}

		result = append(result, x.Verbs()...)
	}

	return result
}
