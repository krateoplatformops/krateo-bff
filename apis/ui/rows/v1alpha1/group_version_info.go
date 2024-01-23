// +kubebuilder:object:generate=true
// +groupName=layout.ui.krateo.io
// +versionName=v1alpha1
package v1alpha1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	Group   = "layout.ui.krateo.io"
	Version = "v1alpha1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// Row type metadata.
var (
	RowKind             = reflect.TypeOf(Row{}).Name()
	RowGroupKind        = schema.GroupKind{Group: Group, Kind: RowKind}.String()
	RowKindAPIVersion   = RowKind + "." + SchemeGroupVersion.String()
	RowGroupVersionKind = SchemeGroupVersion.WithKind(RowKind)
)

func init() {
	SchemeBuilder.Register(&Row{}, &RowList{})
}
