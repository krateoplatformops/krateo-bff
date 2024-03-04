package evaluator

import (
	"context"

	"github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"

	"k8s.io/client-go/rest"
)

const (
	filter = `.spec.versions[] | select(.name="%s") | .schema.openAPIV3Schema`
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Subject    string
}

func Eval(ctx context.Context, in *v1alpha1.FormTemplate, opts EvalOptions) error {
	tpl, err := tmpl.New("${", "}")
	if err != nil {
		return err
	}

	ds, err := callAPIs(ctx, callAPIsOptions{
		restConfig: opts.RESTConfig,
		authnNS:    opts.AuthnNS,
		subject:    opts.Subject,
		tpl:        tpl,
		apiList:    in.Spec.APIList,
	})

	for _, el := range in.Spec.Data {
		el.Value, err = tpl.Execute(el.Value, ds)
		if err != nil {
			return err
		}
	}

	return nil
}
