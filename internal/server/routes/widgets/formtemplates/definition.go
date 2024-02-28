package formtemplates

import (
	"context"

	formdefinitionsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/core/formdefinitions/v1alpha1"
	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/formdefinitions"
	"k8s.io/client-go/rest"
)

func getFormDefinition(ctx context.Context, rc *rest.Config, in *formtemplatesv1alpha1.FormTemplate) (*formdefinitionsv1alpha1.FormDefinition, error) {
	formdefinitionsClient, err := formdefinitions.NewClient(rc)
	if err != nil {
		return nil, err
	}

	return formdefinitionsClient.
		Namespace(in.Spec.DefinitionRef.Namespace).
		Get(ctx, in.Spec.DefinitionRef.Name)
}
