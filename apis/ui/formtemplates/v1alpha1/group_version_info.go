// +kubebuilder:object:generate=true
// +groupName=widgets.ui.krateo.io
// +versionName=v1alpha1
package v1alpha1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	Group   = "widgets.ui.krateo.io"
	Version = "v1alpha1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

var (
	FormTemplateKind             = reflect.TypeOf(FormTemplate{}).Name()
	FormTemplateGroupKind        = schema.GroupKind{Group: Group, Kind: FormTemplateKind}.String()
	FormTemplateKindAPIVersion   = FormTemplateKind + "." + SchemeGroupVersion.String()
	FormTemplateGroupVersionKind = SchemeGroupVersion.WithKind(FormTemplateKind)
)

func init() {
	SchemeBuilder.Register(&FormTemplate{}, &FormTemplateList{})
}
