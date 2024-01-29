package v1alpha1

import (
	"github.com/krateoplatformops/krateo-bff/apis/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

type App struct {
	// +optional
	Props map[string]string `json:"props,omitempty"`
}

type ColumnSpec struct {
	// App is the column content
	App App `json:"app"`

	// CardTemplateListRef reference to card template list.
	// +optional
	CardTemplateListRef []*core.Reference `json:"cardTemplateListRef,omitempty"`
}

type ColumnStatus struct {
	//+kubebuilder:validation:EmbeddedResource
	Content *runtime.RawExtension `json:"content,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={krateo,layout,columns}
type Column struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ColumnSpec   `json:"spec"`
	Status ColumnStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type ColumnList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Column `json:"items"`
}
