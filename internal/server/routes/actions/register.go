package actions

import (
	"github.com/go-chi/chi/v5"
	"k8s.io/client-go/rest"
)

type HandlerOptions struct {
	AuthnNS       string
	KubeServerURL string
	KubeProxyURL  string
}

func Register(r *chi.Mux, rc *rest.Config, opts HandlerOptions) {
	r.Get(newHandler(rc, opts))
	r.Post(newHandler(rc, opts))
	r.Put(newHandler(rc, opts))
	r.Delete(newHandler(rc, opts))
}
