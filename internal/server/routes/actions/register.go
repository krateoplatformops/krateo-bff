package actions

import (
	"github.com/go-chi/chi/v5"
	"k8s.io/client-go/rest"
)

func Register(r *chi.Mux, rc *rest.Config, authnNS string) {
	r.Get(newHandler(rc, authnNS))
	r.Post(newHandler(rc, authnNS))
	r.Put(newHandler(rc, authnNS))
	r.Delete(newHandler(rc, authnNS))
}
