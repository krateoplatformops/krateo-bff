package actions

import (
	"net/http"
	"path"
)

func buildApiPath(method string, opts options) string {
	switch method {
	case http.MethodPost:
		return path.Join("/apis",
			opts.group, opts.version,
			"namespaces", opts.namespace,
			opts.plural)
	default:
		return path.Join("/apis",
			opts.group, opts.version,
			"namespaces", opts.namespace,
			opts.plural, opts.name)
	}
}
