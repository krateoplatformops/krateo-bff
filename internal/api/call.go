package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"k8s.io/utils/ptr"
)

type CallOptions struct {
	API      *core.API
	Endpoint *core.Endpoint
}

func Call(ctx context.Context, client *http.Client, opts CallOptions) (map[string]any, error) {
	uri := strings.TrimSuffix(opts.Endpoint.ServerURL, "/")
	if pt := ptr.Deref(opts.API.Path, ""); len(pt) > 0 {
		uri = fmt.Sprintf("%s/%s", uri, strings.TrimPrefix(pt, "/"))
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	verb := ptr.Deref(opts.API.Verb, http.MethodGet)

	var body io.Reader
	if s := ptr.Deref(opts.API.Payload, ""); len(s) > 0 {
		body = strings.NewReader(s)
	}

	req, err := http.NewRequestWithContext(ctx, verb, u.String(), body)
	if err != nil {
		return nil, err
	}

	if len(opts.API.Headers) > 0 {
		for _, el := range opts.API.Headers {
			idx := strings.Index(el, ":")
			if idx <= 0 {
				continue
			}
			req.Header.Set(el[:idx], el[idx+1:])
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

func decodeResponseBody(resp *http.Response) (map[string]any, error) {
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
