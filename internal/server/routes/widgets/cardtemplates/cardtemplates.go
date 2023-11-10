package cardtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/krateoplatformops/krateo-bff/apis/core"
	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/resolvers"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	"k8s.io/client-go/rest"
)

const (
	listerPathFmt = "/apis/%s/%s/namespaces/{namespace}/cardtemplates"
	getterPathFmt = "/apis/%s/%s/namespaces/{namespace}/cardtemplates/{name}"
)

func Register(r *chi.Mux, rc *rest.Config) {
	r.Get(cardTemplateLister(rc))
	r.Get(cardTemplateGetter(rc))
}

func cardTemplateGetter(rc *rest.Config) (string, http.HandlerFunc) {
	pattern := fmt.Sprintf(getterPathFmt, cardtemplatev1alpha1.Group, cardtemplatev1alpha1.Version)
	handler := &getter{rc: rc}
	return pattern, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

func cardTemplateLister(rc *rest.Config) (string, http.HandlerFunc) {
	pattern := fmt.Sprintf(listerPathFmt, cardtemplatev1alpha1.Group, cardtemplatev1alpha1.Version)
	handler := &lister{rc: rc}
	return pattern, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*getter)(nil)

type getter struct {
	rc *rest.Config
}

func (r *getter) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context()).With().Logger()

	namespace := chi.URLParam(req, "namespace")
	name := chi.URLParam(req, "name")
	eval := true
	if qs := req.URL.Query(); qs.Has("eval") {
		ok, err := strconv.ParseBool(qs.Get("eval"))
		if err == nil {
			eval = ok
		}
	}

	res, err := resolvers.CardTemplateGetOne(context.Background(), r.rc, &core.Reference{
		Name: name, Namespace: namespace,
	}, eval)
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
		log.Err(err).Msg("unable to serve json encoded cardtemplate")
	}
}

var _ http.Handler = (*lister)(nil)

type lister struct {
	rc *rest.Config
}

func (r *lister) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context()).With().Logger()

	namespace := chi.URLParam(req, "namespace")
	eval := true
	if qs := req.URL.Query(); qs.Has("eval") {
		ok, err := strconv.ParseBool(qs.Get("eval"))
		if err == nil {
			eval = ok
		}
	}

	res, err := resolvers.CardTemplateGetAll(context.Background(), r.rc, namespace, eval)
	if err != nil {
		log.Err(err).
			Str("namespace", namespace).
			Msg("unable to resolve card templates")
		encode.Error(wri, http.StatusInternalServerError, err)
		return
	}

	wri.WriteHeader(http.StatusOK)
	wri.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		log.Err(err).Msg("unable to serve json encoded cardtemplates")
	}
}
