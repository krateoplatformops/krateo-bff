package v1alpha1

import (
	"github.com/krateoplatformops/krateo-bff/apis/core"
	columnsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RowSpec struct {
	// ColumnListRef reference to column list.
	// +optional
	ColumnListRef []*core.Reference `json:"columnListRef,omitempty"`
}

type RowStatus struct {
	// +optional
	Columns []*columnsv1alpha1.ColumnStatus `json:"columns,omitempty"`
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
