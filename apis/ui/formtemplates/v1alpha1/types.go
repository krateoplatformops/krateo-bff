package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type DataItem struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

func (di *DataItem) String() string {
	return fmt.Sprintf("%s=%s",
		strings.TrimSpace(di.Path), strings.TrimSpace(di.Value))
}

type Action struct {
	//+kubebuilder:validation:Required
	Path string `json:"path"`

	// +optional
	// +kubebuilder:default=GET
	Verb string `json:"verb,omitempty"`
}

type FormTemplateSpec struct {
	SchemaDefinitionRef      *core.Reference `json:"schemaDefinitionRef,omitempty"`
	CompositionDefinitionRef *core.Reference `json:"compositionDefinitionRef,omitempty"`
}

type FormTemplateStatusContent struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	Schema *runtime.RawExtension `json:"schema,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	Instance *runtime.RawExtension `json:"instance,omitempty"`
}

type FormTemplateStatus struct {
	Content *FormTemplateStatusContent `json:"content,omitempty"`
	Actions []*Action                  `json:"actions,omitempty"`
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

type FormTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FormTemplate `json:"items"`
}
