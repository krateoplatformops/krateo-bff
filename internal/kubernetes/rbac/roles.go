package rbac

import (
	"context"
	"crypto/x509/pkix"
	"strings"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/objects"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/utils/strings/slices"
)

func RolesForSubject(restConfig *rest.Config, sub pkix.Name, namespace string) ([]RoleInfo, error) {
	res := []RoleInfo{}

	if len(namespace) == 0 {
		all, err := findClusterRolesForSubject(restConfig, sub)
		if err != nil {
			return nil, err
		}
		res = append(res, all...)
	}

	all, err := findRolesForSubject(restConfig, sub, namespace)
	if err != nil {
		return nil, err
	}
	res = append(res, all...)

	return res, nil
}

func findClusterRolesForSubject(restConfig *rest.Config, sub pkix.Name) ([]RoleInfo, error) {
	resolver, err := objects.NewObjectResolver(restConfig)
	if err != nil {
		return nil, err
	}

	binds, err := resolver.List(context.TODO(), schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "ClusterRoleBinding",
	}, "")
	if err != nil {
		return nil, err
	}

	acceptFn := func(si SubjectInfo) bool {
		ok := strings.EqualFold(si.Kind(), "group")
		ok = ok && slices.Contains(sub.Organization, si.Name())
		return ok
	}

	roles := extractRoleInfo(binds, acceptFn)
	return roles, nil
}

func findRolesForSubject(restConfig *rest.Config, sub pkix.Name, namespace string) ([]RoleInfo, error) {
	resolver, err := objects.NewObjectResolver(restConfig)
	if err != nil {
		return nil, err
	}

	binds, err := resolver.List(context.TODO(), schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "RoleBinding",
	}, namespace)
	if err != nil {
		return nil, err
	}

	acceptFn := func(si SubjectInfo) bool {
		ok := strings.EqualFold(si.Kind(), "group")
		ok = ok && slices.Contains(sub.Organization, si.Name())
		return ok
	}

	roles := extractRoleInfo(binds, acceptFn)
	return roles, nil
}

func extractRoleInfo(binds *unstructured.UnstructuredList, acceptFn func(SubjectInfo) bool) []RoleInfo {
	res := []RoleInfo{}

	for _, el := range binds.Items {
		roleKind, _, _ := unstructured.NestedString(el.Object, "roleRef", "kind")
		roleName, _, _ := unstructured.NestedString(el.Object, "roleRef", "name")
		roleNamespace, ok, err := unstructured.NestedString(el.Object, "roleRef", "namespace")
		if !ok || err != nil {
			roleNamespace, _, _ = unstructured.NestedString(el.Object, "metadata", "namespace")
		}

		subs, _, _ := unstructured.NestedSlice(el.Object, "subjects")
		for _, x := range subs {
			y, ok := x.(map[string]interface{})
			if !ok {
				continue
			}
			subjectKind, _, _ := unstructured.NestedString(y, "kind")
			subjectName, _, _ := unstructured.NestedString(y, "name")

			if acceptFn != nil && acceptFn(&subjectInfo{
				kind: subjectKind,
				name: subjectName,
			}) {
				res = append(res, &roleInfo{
					kind:      roleKind,
					name:      roleName,
					namespace: roleNamespace,
				})
			}
		}
	}

	return res
}
