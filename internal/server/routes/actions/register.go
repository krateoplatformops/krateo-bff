package actions

import (
	"github.com/go-chi/chi/v5"
	"k8s.io/client-go/rest"
)

func Register(r *chi.Mux, rc *rest.Config, authnNS string) {
	r.Post(newCaller(rc, authnNS))
}
