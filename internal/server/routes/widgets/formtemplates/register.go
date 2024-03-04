package formtemplates

import (
	"github.com/go-chi/chi/v5"
	"k8s.io/client-go/rest"
)

const (
	allowedVerbsAnnotationKey = "krateo.io/allowed-verbs"
	group                     = "widgets.ui.krateo.io"
	resource                  = "formtemplates"
)

func Register(r *chi.Mux, rc *rest.Config, authnNS string) {
	r.Get(newGetter(rc, authnNS))
}
