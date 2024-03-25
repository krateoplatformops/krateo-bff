package formtemplates

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/formtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	getterPath = "/apis/widgets.ui.krateo.io/formtemplates/{name}"
)

func newGetter(rc *rest.Config, authnNS string) (string, http.HandlerFunc) {
	handler := &getter{
		rc: rc,
		gr: schema.GroupResource{
			Group:    group,
			Resource: resource,
		},
		authnNS: authnNS,
	}
	return getterPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*getter)(nil)

type getter struct {
	rc              *rest.Config
	gr              schema.GroupResource
	templatesClient *formtemplates.Client
	authnNS         string
}

func (r *getter) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")

	qs := req.URL.Query()

	namespace := qs.Get("namespace")
	sub := qs.Get("sub")
	orgs := strings.Split(qs.Get("orgs"), ",")
	version := qs.Get("version")
	if len(version) == 0 {
		version = "v1alpha1"
	}

	log := zerolog.Ctx(req.Context()).With().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("name", name).
		Str("namespace", namespace).
		Str("version", version).
		Logger()

	if err := r.complete(); err != nil {
		log.Err(err).Msg("unable to initialize rest clients")
		encode.InternalError(wri, err)
		return
	}

	obj, err := r.templatesClient.Get(context.Background(), formtemplates.GetOptions{
		Name:      name,
		Namespace: namespace,
		Subject:   sub,
		Orgs:      orgs,
	})
	if err != nil {
		log.Err(err).Msg("unable to resolve form template")
		if apierrors.IsNotFound(err) {
			encode.NotFound(wri, err)
		} else {
			encode.Invalid(wri, err)
		}
		return
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(obj); err != nil {
		log.Err(err).Msg("unable to serve json encoded form template")
	}
}

func (r *getter) complete() error {
	if r.templatesClient == nil {
		cli, err := formtemplates.NewClient(r.rc, true)
		if err != nil {
			return err
		}

		r.templatesClient = cli
	}

	return nil
}
