package util

import (
	"context"
	"slices"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ResourceInfo struct {
	Subject       string
	Groups        []string
	GroupResource schema.GroupResource
	ResourceName  string
	Namespace     string
}

func Can(what string, verbs []string) bool {
	if slices.Contains(verbs, "*") {
		return true
	}

	return slices.Contains(verbs, what)
}

func CanListResource(ctx context.Context, restConfig *rest.Config, opts ResourceInfo) (bool, error) {
	verbs, err := GetAllowedVerbs(ctx, restConfig, opts)
	if err != nil {
		return false, err
	}

	if slices.Contains(verbs, "*") {
		return true, nil
	}

	if slices.Contains(verbs, "list") {
		return true, nil
	}

	if slices.Contains(verbs, "get") {
		return true, nil
	}

	return false, nil
}

func CanDeleteResource(ctx context.Context, restConfig *rest.Config, opts ResourceInfo) (bool, error) {
	verbs, err := GetAllowedVerbs(ctx, restConfig, opts)
	if err != nil {
		return false, err
	}
	if slices.Contains(verbs, "*") {
		return true, nil
	}

	return slices.Contains(verbs, "delete"), nil
}

func GetAllowedVerbs(ctx context.Context, restConfig *rest.Config, opts ResourceInfo) ([]string, error) {
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

		ok := slices.Contains(x.Resources(), "*")
		if !ok {
			ok = slices.Contains(x.Resources(), gr.Resource)
			//ok = ok && !checkResourceName
		}
		if !ok {
			continue
		}

		if checkResourceName {
			if slices.Contains(x.ResourceNames(), resourceName) {
				result = append(result, x.Verbs()...)
				continue
			}
		}

		if len(x.ResourceNames()) > 0 {
			if slices.Contains(x.Verbs(), "get") {
				result = append(result, x.Verbs()...)
			}
			continue
		}

		result = append(result, x.Verbs()...)
	}

	return result
}
