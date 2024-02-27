package v1alpha1

import (
	"github.com/krateoplatformops/krateo-bff/apis/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type DataItem struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

type FormTemplateSpec struct {
	// DefinitionRef: reference to FormDefintion
	DefinitionRef *core.Reference `json:"definitionRef"`

	// ResourceRef: reference to resource instance
	ResourceRef *core.Reference `json:"resourceRef"`

	// +optional
	Data []*DataItem `json:"data,omitempty"`

	// APIList list of api calls.
	// +optional
	APIList []*core.API `json:"api,omitempty"`
}

type FormTemplateStatusContent struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	Schema *runtime.RawExtension `json:"schema,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	Instance *runtime.RawExtension `json:"instance,omitempty"`
}

type FormTemplateStatus struct {
	Content *FormTemplateStatusContent `json:"content,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={krateo,template,widgets}

type FormTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FormTemplateSpec `json:"spec,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	Status FormTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FormTemplateList contains a list of FormTemplate
type FormTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FormTemplate `json:"items"`
}
