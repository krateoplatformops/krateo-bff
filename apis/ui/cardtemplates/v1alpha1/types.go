package v1alpha1

import (
	"github.com/krateoplatformops/krateo-bff/apis/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Card struct {
	Title string `json:"title"`

	Content string `json:"content"`

	// +optional
	Icon string `json:"icon,omitempty"`

	// +optional
	Color string `json:"color,omitempty"`

	// +optional
	Date string `json:"date,omitempty"`

	// +optional
	Tags string `json:"tags,omitempty"`

	// +optional
	Actions []*core.API `json:"actions,omitempty"`

	// +optional
	AllowedActions []string `json:"allowedActions,omitempty"`
}

// CardTemplate is a template for a Krateo UI Card widget.
type CardTemplateSpec struct {
	// +optional
	Iterator *string `json:"iterator,omitempty"`

	// CardTemplateInfo is the card template info
	CardTemplateInfo Card `json:"app"`

	// APIList list of api calls.
	// +optional
	APIList []*core.API `json:"api,omitempty"`
}

type CardTemplateStatus struct {
	Cards []*Card `json:"content,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={krateo,cards,widgets}

// CardTemplate is ui widgets card configuration.
type CardTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CardTemplateSpec   `json:"spec,omitempty"`
	Status CardTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CardTemplateList contains a list of CardTemplate
type CardTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CardTemplate `json:"items"`
}
