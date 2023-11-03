package core

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"k8s.io/utils/ptr"
)

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

func (r *API) Execute(ctx context.Context, client *http.Client, authn *AuthInfo) (any, error) {
	uri := strings.TrimSuffix(r.Server, "/")
	if pt := ptr.Deref(r.Path, ""); len(pt) > 0 {
		uri = fmt.Sprintf(uri, pt)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, ptr.Deref(r.Verb, ""), u.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(r.Headers) > 0 {
		for _, el := range r.Headers {
			idx := strings.Index(el, ":")
			if idx <= 0 {
				continue
			}
			req.Header.Set(el[:idx], el[idx+1:])
		}
	}

	if authn != nil {
		if authn.IsBasicAuth() {
			req.SetBasicAuth(authn.Username, authn.Password)
		}

		if authn.IsTokenAuth() {
			req.Header.Set("Authorization", "Bearer: "+authn.Token)
		}

		cert, err := tls.X509KeyPair(authn.CertificateAuthorityData, authn.ClientKeyData)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(authn.CertificateAuthorityData)
		tlsConfig := &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{cert},
		}

		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		v, err := decodeResponseBody(resp)
		if err != nil {
			return nil, fmt.Errorf("http response: %s", resp.Status)
		}
		return v, nil
	}

	return decodeResponseBody(resp)
}

func decodeResponseBody(resp *http.Response) (any, error) {
	if !hasContentType(resp, "application/json") {
		return nil, fmt.Errorf("only 'application/json' media type is supported")
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	v := map[string]any{}
	if err := json.Unmarshal(dat, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Determine whether the request `content-type` includes a
// server-acceptable mime-type
//
// Failure should yield an HTTP 415 (`http.StatusUnsupportedMediaType`)
func hasContentType(r *http.Response, mimetype string) bool {
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
