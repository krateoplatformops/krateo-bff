package evaluator

import (
	"context"

	cardtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	cardtemplatesevaluator "github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates/evaluator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Username   string
}

func Eval(ctx context.Context, in *v1alpha1.Column, opts EvalOptions) error {
	return evalCardTemplateList(ctx, in, opts)
}

func evalCardTemplateList(ctx context.Context, in *v1alpha1.Column, opts EvalOptions) error {
	ref := in.Spec.CardTemplateListRef
	if ref == nil {
		return nil
	}

	in.Status.CardTemplateList = []cardtemplatesv1alpha1.CardTemplate{}

	cli, err := cardtemplates.NewClient(opts.RESTConfig)
	if err != nil {
		return err
	}
	all, err := cli.Namespace(ref.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, el := range all.Items {
		obj := &el
		err := cardtemplatesevaluator.Eval(ctx, obj, cardtemplatesevaluator.EvalOptions{
			RESTConfig: opts.RESTConfig,
			AuthnNS:    opts.AuthnNS,
			Username:   opts.Username,
		})
		if err != nil {
			return err
		}

		in.Status.CardTemplateList = append(in.Status.CardTemplateList, *obj)
	}

	return nil
}
