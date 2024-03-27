package actions

import (
	"net/http"
	"strconv"
)

type options struct {
	name      string
	namespace string
	subject   string
	orgs      string
	plural    string
	group     string
	version   string
	kind      string
	verbose   bool
}

func optionsFromRequest(req *http.Request) (opts options) {
	qs := req.URL.Query()

	opts.verbose, _ = strconv.ParseBool(qs.Get("verbose"))
	opts.version = qs.Get("version")
	opts.group = qs.Get("group")
	opts.kind = qs.Get("kind")
	opts.plural = qs.Get("plural")
	opts.name = qs.Get("name")
	opts.namespace = qs.Get("namespace")
	opts.subject = qs.Get("sub")
	opts.orgs = qs.Get("orgs")

	return
}
