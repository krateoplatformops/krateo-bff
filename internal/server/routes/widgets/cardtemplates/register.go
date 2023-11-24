package cardtemplates

import (
	"github.com/go-chi/chi/v5"
	"k8s.io/client-go/rest"
)

const (
	allowedVerbsAnnotationKey = "krateo.io/allowed-verbs"
)

func Register(r *chi.Mux, rc *rest.Config) {
	r.Get(newLister(rc))
	r.Get(newGetter(rc))
}
