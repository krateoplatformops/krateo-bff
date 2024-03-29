package columns

import (
	"github.com/go-chi/chi/v5"
	"k8s.io/client-go/rest"
)

const (
	group    = "layout.ui.krateo.io"
	resource = "columns"
)

func Register(r *chi.Mux, rc *rest.Config, authnNS string) {
	r.Get(newGetter(rc, authnNS))
	r.Get(newLister(rc, authnNS))
}
