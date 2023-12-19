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

// CardTemplate type metadata.
var (
	CardTemplateKind             = reflect.TypeOf(CardTemplate{}).Name()
	CardTemplateGroupKind        = schema.GroupKind{Group: Group, Kind: CardTemplateKind}.String()
	CardTemplateKindAPIVersion   = CardTemplateKind + "." + SchemeGroupVersion.String()
	CardTemplateGroupVersionKind = SchemeGroupVersion.WithKind(CardTemplateKind)
)

func init() {
	SchemeBuilder.Register(&CardTemplate{}, &CardTemplateList{})
}
