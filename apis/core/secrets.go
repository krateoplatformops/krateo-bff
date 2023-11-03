package core

// A Reference to a named object.
// +k8s:deepcopy-gen=true
type Reference struct {
	// Name of the referenced object.
	Name string `json:"name"`

	// Namespace of the referenced object.
	Namespace string `json:"namespace"`
}

// A SecretKeySelector is a reference to a secret key in an arbitrary namespace.
// +k8s:deepcopy-gen=true
type SecretKeySelector struct {
	Reference `json:",inline"`

	// The key to select.
	Key string `json:"key"`
}
