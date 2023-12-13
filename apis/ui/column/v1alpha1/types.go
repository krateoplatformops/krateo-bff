package v1alpha1

import (
	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Content struct {
	CardTemplateList []v1alpha1.CardTemplate `json:"cardTemplateList"`
}

type App struct {
	Content Content `json:"content"`

	// +optional
	Props map[string]string `json:"props,omitempty"`
}

type ColumnSpec struct {
	// App is the column content
	App App `json:"app"`

	// APIList list of api calls.
	// +optional
	APIList []*core.API `json:"api,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,categories={krateo,layout}

type Column struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ColumnSpec `json:"spec"`
}

// +kubebuilder:object:root=true

type ColumnList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Column `json:"items"`
}
