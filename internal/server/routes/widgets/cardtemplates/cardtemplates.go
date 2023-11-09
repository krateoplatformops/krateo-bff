package cardtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/krateoplatformops/krateo-bff/apis/core"
	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/resolvers"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	"k8s.io/client-go/rest"
)

func GetPath() string {
	return fmt.Sprintf("/apis/%s/%s/namespaces/{namespace}/cardtemplates/{name}",
		cardtemplatev1alpha1.Group, cardtemplatev1alpha1.Version)
}

var _ http.Handler = (*cardTemplateGet)(nil)

func Get(rc *rest.Config) http.Handler {
	return &cardTemplateGet{rc: rc}
}

type cardTemplateGet struct {
	rc *rest.Config
}

func (r *cardTemplateGet) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context()).With().Logger()

	namespace := chi.URLParam(req, "namespace")
	name := chi.URLParam(req, "name")

	res, err := resolvers.CardTemplate(context.Background(), r.rc, &core.Reference{
		Name: name, Namespace: namespace,
	})
	if err != nil {
		log.Err(err).
			Str("name", name).
			Str("namespace", namespace).
			Msg("unable to resolve card template")
		encode.Error(wri, http.StatusInternalServerError, err)
		return
	}

	wri.WriteHeader(http.StatusOK)
	wri.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		log.Err(err).Msg("unable to serve json encoded card")
	}
}
