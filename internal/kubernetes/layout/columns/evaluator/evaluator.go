package evaluator

import (
	"context"

	cardtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	cardtemplatesevaluator "github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates/evaluator"
	"k8s.io/client-go/rest"
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Username   string
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

	if in.Status.Cards == nil {
		in.Status.Cards = []*cardtemplatesv1alpha1.CardInfo{}
	}

	for _, ref := range refs {
		obj, err := cli.Namespace(ref.Namespace).Get(ctx, ref.Name)
		if err != nil {
			return err
		}

		err = cardtemplatesevaluator.Eval(ctx, obj, cardtemplatesevaluator.EvalOptions{
			RESTConfig: opts.RESTConfig,
			AuthnNS:    opts.AuthnNS,
			Username:   opts.Username,
		})
		if err != nil {
			return err
		}

		in.Status.Cards = append(in.Status.Cards, obj.Status.Cards...)
	}

	return nil
}
