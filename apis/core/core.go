package core

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

// AuthInfo contains information that describes identity information.
// +k8s:deepcopy-gen=true
type AuthInfo struct {
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
}

func (r *AuthInfo) IsBasicAuth() bool {
	return len(r.Password) > 0
}

func (r *AuthInfo) IsTokenAuth() bool {
	return len(r.Token) > 0
}

func (r *AuthInfo) IsCertAuth() bool {
	return len(r.ClientCertificateData) > 0 && len(r.ClientKeyData) > 0
}

func (r *AuthInfo) SetTLSClientConfig(client *http.Client) error {
	if len(r.ClientCertificateData) == 0 || len(r.ClientKeyData) == 0 {
		return nil
	}

	cert, err := tls.X509KeyPair(r.ClientCertificateData, r.ClientKeyData)
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(r.CertificateAuthorityData)
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return nil
}
