package apis

import (
	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	columnv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/column/v1alpha1"
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
		cardtemplatev1alpha1.SchemeBuilder.AddToScheme,
		columnv1alpha1.SchemeBuilder.AddToScheme,
	)
}
