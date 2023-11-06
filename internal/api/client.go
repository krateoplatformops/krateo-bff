package api

import (
	"fmt"
	"net/http"

	"github.com/krateoplatformops/krateo-bff/apis/core"
)

func HTTPClientForEndpoint(authn *core.Endpoint) (*http.Client, error) {
	rt, err := tlsConfigFor(authn)
	if err != nil {
		return http.DefaultClient, err
	}

	if authn.Debug {
		rt = &debuggingRoundTripper{
			delegatedRoundTripper: rt,
		}
	}

	// Set authentication wrappers
	switch {
	case authn.HasBasicAuth() && authn.HasTokenAuth():
		return nil, fmt.Errorf("username/password or bearer token may be set, but not both")

	case authn.HasTokenAuth():
		rt = &bearerAuthRoundTripper{
			bearer: authn.Token,
			rt:     rt,
		}

	case authn.HasBasicAuth():
		rt = &basicAuthRoundTripper{
			username: authn.Username,
			password: authn.Password,
			rt:       rt,
		}
	}

	return &http.Client{Transport: rt}, nil
}
