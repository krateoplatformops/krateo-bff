package verbs

import (
	"github.com/go-chi/chi/v5"
	"k8s.io/client-go/rest"
)

// GET /apis/allowed-verbs?sub=cyberjoker&orgs=devs&gr=cardtemplates.widgets.ui.krateo.io&name=card-dev&namespace=dev-system

func Register(r *chi.Mux, rc *rest.Config) {
	r.Get(newGetter(rc))
}
