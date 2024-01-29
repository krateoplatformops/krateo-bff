package evaluator

import (
	"context"
	"strings"

	cardtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	cardtemplatesevaluator "github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates/evaluator"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

const (
	listKind = "CardTemplateList"
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Subject    string
	Groups     []string
}

func Eval(ctx context.Context, in *v1alpha1.Column, opts EvalOptions) error {
	return evalCardTemplateRefs(ctx, in, opts)
}

func evalCardTemplateRefs(ctx context.Context, in *v1alpha1.Column, opts EvalOptions) error {
	refs := in.Spec.CardTemplateListRef
	if refs == nil {
		return nil
	}

	cli, err := cardtemplates.NewClient(opts.RESTConfig)
	if err != nil {
		return err
	}

	all := &cardtemplatesv1alpha1.CardTemplateList{
		Items: []cardtemplatesv1alpha1.CardTemplate{},
	}
	all.SetGroupVersionKind(cardtemplatesv1alpha1.SchemeGroupVersion.WithKind(listKind))

	for _, ref := range refs {
		obj, err := cli.Namespace(ref.Namespace).Get(ctx, ref.Name)
		if err != nil {
			return err
		}

		err = cardtemplatesevaluator.Eval(ctx, obj, cardtemplatesevaluator.EvalOptions{
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
	resource                  = "columns"
)

type allowedVerbsInjectorOptions struct {
	restConfig *rest.Config
	subject    string
	groups     []string
}

func injectAllowedVerbs(in *v1alpha1.Column, opts allowedVerbsInjectorOptions) error {
	verbs, err := rbacutil.GetAllowedVerbs(context.TODO(), opts.restConfig, util.ResourceInfo{
		Subject: opts.subject,
		Groups:  opts.groups,
		GroupResource: v1alpha1.ColumnGroupVersionKind.GroupVersion().
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
