package v1alpha1

import (
	rtv1 "github.com/krateoplatformops/provider-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SchemaInfo struct {
	// Url of the values.schema.json file
	// +kubebuilder:validation:Required
	Url string `json:"url"`

	// Group: collection of kinds.
	// +kubebuilder:validation:Required
	Group string `json:"group"`

	// Version: allow Kubernetes to release groups as tagged versions.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Kind: the name of the object you are trying to generate
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
}

// FormDefinitionSpec is the specification of a Definition.
type FormDefinitionSpec struct {
	rtv1.ManagedSpec `json:",inline"`

	// Schema: the schema info
	// +immutable
	Schema SchemaInfo `json:"schema"`
}

// FormDefinitionStatus is the status of a Definition.
type FormDefinitionStatus struct {
	rtv1.ManagedStatus `json:",inline"`

	// Resource: the generated custom resource
	// +optional
	Resource string `json:"resource,omitempty"`

	// Digest: schema digest
	// +optional
	Digest *string `json:"digest,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced,categories={krateo,definition,frontend,forms}
//+kubebuilder:printcolumn:name="RESOURCE",type="string",JSONPath=".status.resource"
//+kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",priority=10

// FormDefinition is a definition type with a spec and a status.
type FormDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FormDefinitionSpec   `json:"spec,omitempty"`
	Status FormDefinitionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FormDefinitionList is a list of Definition objects.
type FormDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []FormDefinition `json:"items"`
}
