package rbac

import (
	"crypto/x509/pkix"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/utils/strings/slices"
)

func CanSubjectGetResource(restConfig *rest.Config, sub pkix.Name, gr schema.GroupResource, name, namespace string) (bool, error) {
	roles, err := RolesForSubject(restConfig, sub, namespace)
	if err != nil {
		return false, err
	}

	resourceName := len(name) > 0
	for _, x := range roles {
		all, err := RulesForRole(restConfig, x)
		if err != nil {
			return false, err
		}

		for _, y := range all {
			if !slices.Contains(y.APIGroups(), gr.Group) {
				continue
			}

			if !slices.Contains(y.Verbs(), "get") {
				continue
			}

			if resourceName {
				if !slices.Contains(y.ResourceNames(), name) {
					continue
				}
			}

			ok := slices.Contains(y.Resources(), "*")
			if !ok {
				ok = slices.Contains(y.Resources(), gr.Resource)
			}
			if !ok {
				continue
			}

			if x.Namespace() == namespace {
				return true, nil
			}
		}
	}

	return false, nil
}

func AllowedVerbsOnResourceForSubject(restConfig *rest.Config, sub pkix.Name, gr schema.GroupResource, name, namespace string) ([]string, error) {
	roles, err := RolesForSubject(restConfig, sub, namespace)
	if err != nil {
		return []string{}, err
	}

	rules := []PolicyRule{}
	for _, x := range roles {
		itm, err := RulesForRole(restConfig, x)
		if err != nil {
			return []string{}, err
		}

		if x.Namespace() == namespace {
			rules = append(rules, itm...)
		}
	}

	resourceName := len(name) > 0

	allow := []string{}
	for _, x := range rules {
		if !slices.Contains(x.APIGroups(), gr.Group) {
			continue
		}

		if resourceName {
			if slices.Contains(x.ResourceNames(), name) {
				allow = append(allow, x.Verbs()...)
				continue
			}
		}

		ok := slices.Contains(x.Resources(), "*")
		if !ok {
			ok = slices.Contains(x.Resources(), gr.Resource)
		}
		if !ok {
			continue
		}

		allow = append(allow, x.Verbs()...)
	}

	return allow, nil
}
