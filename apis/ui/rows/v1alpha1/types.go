package v1alpha1

import (
	"github.com/krateoplatformops/krateo-bff/apis/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

type RowSpec struct {
	// ColumnListRef reference to column list.
	// +optional
	ColumnListRef []*core.Reference `json:"columnListRef,omitempty"`
}

type RowStatus struct {
	//+kubebuilder:validation:EmbeddedResource
	Content *runtime.RawExtension `json:"columns,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={krateo,layout,rows}
type Row struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RowSpec   `json:"spec"`
	Status RowStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type RowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Row `json:"items"`
}
