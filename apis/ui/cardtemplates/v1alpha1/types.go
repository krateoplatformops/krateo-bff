package v1alpha1

import (
	"github.com/krateoplatformops/krateo-bff/apis/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FormTemplateRef struct {
	// Name of the referenced object.
	Name string `json:"name"`

	// Namespace of the referenced object.
	Namespace string `json:"namespace,omitempty"`
}

type Action struct {
	//+kubebuilder:validation:Required
	Path string `json:"path"`

	// +optional
	// +kubebuilder:default=GET
	Verb string `json:"verb,omitempty"`
}

type Card struct {
	//+kubebuilder:validation:Required
	Title string `json:"title"`

	//+kubebuilder:validation:Required
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
	// AllowedActions []string `json:"allowedActions,omitempty"`
}

// CardTemplate is a template for a Krateo UI Card widget.
type CardTemplateSpec struct {
	//+kubebuilder:validation:Required
	FormTemplateRef FormTemplateRef `json:"formTemplateRef"`

	//+kubebuilder:validation:Required
	// App is the card template info
	App Card `json:"app"`

	// +optional
	Iterator *string `json:"iterator,omitempty"`

	// APIList list of api calls.
	// +optional
	APIList []*core.API `json:"api,omitempty"`
}

type CardTemplateStatus struct {
	Cards   []*Card   `json:"content,omitempty"`
	Actions []*Action `json:"actions,omitempty"`
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
