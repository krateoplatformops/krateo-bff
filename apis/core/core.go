package core

// Endpoint contains information that describes identity information.
// +k8s:deepcopy-gen=true
type Endpoint struct {
	Server string `json:"server"`

	// CertificateAuthorityData contains PEM-encoded certificate authority certificates.
	CertificateAuthorityData []byte `json:"certificate-authority-data,omitempty"`

	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS.
	ClientCertificateData []byte `json:"client-certificate-data,omitempty"`

	// ClientKeyData contains PEM-encoded data from a client key file for TLS.
	ClientKeyData []byte `json:"client-key-data,omitempty"`

	// Token is the bearer token for authentication to the server.
	Token string `json:"token,omitempty"`

	// Username is the username for basic authentication to the server.
	Username string `json:"username,omitempty"`

	// Password is the password for basic authentication to the server.
	Password string `json:"password,omitempty"`

	Debug bool `json:"debug,omitempty"`
}

// HasCA returns whether the configuration has a certificate authority or not.
func (ep *Endpoint) HasCA() bool {
	return len(ep.CertificateAuthorityData) > 0
}

// HasBasicAuth returns whether the configuration has basic authentication or not.
func (ep *Endpoint) HasBasicAuth() bool {
	return len(ep.Password) != 0
}

// HasTokenAuth returns whether the configuration has token authentication or not.
func (ep *Endpoint) HasTokenAuth() bool {
	return len(ep.Token) != 0
}

// HasCertAuth returns whether the configuration has certificate authentication or not.
func (ep *Endpoint) HasCertAuth() bool {
	return len(ep.ClientCertificateData) != 0 && len(ep.ClientKeyData) != 0
}

// API contains external api call info.
// +k8s:deepcopy-gen=true
type API struct {
	Name string `json:"name"`

	// +optional
	Path *string `json:"path,omitempty"`

	// +optional
	// +kubebuilder:default=GET
	Verb *string `json:"verb,omitempty"`

	// +optional
	Headers []string `json:"headers,omitempty"`

	// +optional
	Payload *string `json:"payload,omitempty"`

	// +optional
	EndpointRef *Reference `json:"endpointRef,omitempty"`

	// +optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`
}
