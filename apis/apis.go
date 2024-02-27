package apis

import (
	formdefinitionsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/core/formdefinitions/v1alpha1"
	cardtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	columnsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	rowsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/rows/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}

func init() {
	AddToSchemes = append(AddToSchemes,
		cardtemplatesv1alpha1.SchemeBuilder.AddToScheme,
		columnsv1alpha1.SchemeBuilder.AddToScheme,
		rowsv1alpha1.SchemeBuilder.AddToScheme,
		formtemplatesv1alpha1.SchemeBuilder.AddToScheme,
		formdefinitionsv1alpha1.SchemeBuilder.AddToScheme,
	)
}
