package evaluator

import (
	"context"
	"strings"

	columnsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/apis/ui/rows/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/columns"
	columnsevaluator "github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/columns/evaluator"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

const (
	listKind = "ColumnList"
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Subject    string
	Groups     []string
}

func Eval(ctx context.Context, in *v1alpha1.Row, opts EvalOptions) error {
	refs := in.Spec.ColumnListRef
	if refs == nil {
		return nil
	}

	cli, err := columns.NewClient(opts.RESTConfig)
	if err != nil {
		return err
	}

	all := &columnsv1alpha1.ColumnList{
		Items: []columnsv1alpha1.Column{},
	}
	all.SetGroupVersionKind(columnsv1alpha1.SchemeGroupVersion.WithKind(listKind))

	for _, ref := range refs {
		obj, err := cli.Namespace(ref.Namespace).Get(ctx, ref.Name)
		if err != nil {
			return err
		}

		err = columnsevaluator.Eval(ctx, obj, columnsevaluator.EvalOptions{
			RESTConfig: opts.RESTConfig,
			AuthnNS:    opts.AuthnNS,
			Subject:    opts.Subject,
			Groups:     opts.Groups,
		})
		if err != nil {
			return err
		}

		all.Items = append(all.Items, *obj)
	}

	in.Status.Content = &runtime.RawExtension{
		Object: all,
	}

	return nil
}

const (
	allowedVerbsAnnotationKey = "krateo.io/allowed-verbs"
	resource                  = "rows"
)

type allowedVerbsInjectorOptions struct {
	restConfig *rest.Config
	subject    string
	groups     []string
}

func injectAllowedVerbs(in *v1alpha1.Row, opts allowedVerbsInjectorOptions) error {
	verbs, err := rbacutil.GetAllowedVerbs(context.TODO(), opts.restConfig, util.ResourceInfo{
		Subject: opts.subject,
		Groups:  opts.groups,
		GroupResource: v1alpha1.RowGroupVersionKind.GroupVersion().
			WithResource(resource).
			GroupResource(),
		ResourceName: in.GetName(),
		Namespace:    in.GetNamespace(),
	})
	if err != nil {
		return err
	}

	m := in.GetAnnotations()
	if len(m) == 0 {
		m = map[string]string{}
	}
	m[allowedVerbsAnnotationKey] = strings.Join(verbs, ",")
	in.SetAnnotations(m)

	return nil
}
