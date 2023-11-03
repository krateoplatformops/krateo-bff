package core

// API contains external api call info.
// +k8s:deepcopy-gen=true
type API struct {
	// +optional
	Name *string `json:"name,omitempty"`

	Server string `json:"server"`

	// +optional
	Path *string `json:"path,omitempty"`

	// +optional
	// +kubebuilder:default=GET
	Verb *string `json:"verb,omitempty"`

	// +optional
	Headers []string `json:"headers,omitempty"`

	// +optional
	EndpointRef *Reference `json:"endpointRef,omitempty"`

	// +optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`
}
