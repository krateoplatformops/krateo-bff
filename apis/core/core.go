package core

// Endpoint contains information that describes identity information.
// +k8s:deepcopy-gen=true
type Endpoint struct {
	ServerURL string `json:"server-url"`

	// +optional
	ProxyURL string `json:"proxy-url,omitempty"`

	// CertificateAuthorityData contains Base64 PEM-encoded certificate authority certificates.
	CertificateAuthorityData string `json:"certificate-authority-data,omitempty"`

	// ClientCertificateData contains Base64 PEM-encoded data from a client cert file for TLS.
	ClientCertificateData string `json:"client-certificate-data,omitempty"`

	// ClientKeyData contains Base64 PEM-encoded data from a client key file for TLS.
	ClientKeyData string `json:"client-key-data,omitempty"`

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
	//Enabled *bool `json:"enabled,omitempty"`

	// +optional
	KrateoGateway *bool `json:"krateoGateway,omitempty"`
}
