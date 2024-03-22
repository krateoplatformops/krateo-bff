package formtemplates

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/formtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	listerPath = "/apis/widgets.ui.krateo.io/formtemplates"
)

func newLister(rc *rest.Config, authnNS string) (string, http.HandlerFunc) {
	handler := &lister{
		rc:      rc,
		authnNS: authnNS,
		gr: schema.GroupResource{
			Group:    group,
			Resource: resource,
		},
	}
	return listerPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*lister)(nil)

type lister struct {
	rc              *rest.Config
	gr              schema.GroupResource
	templatesClient *formtemplates.Client
	authnNS         string
}

func (r *lister) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
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
		Str("namespace", namespace).
		Str("version", version).
		Logger()

	if err := r.complete(); err != nil {
		log.Err(err).Msg("unable to initialize rest clients")
		encode.InternalError(wri, err)
		return
	}

	all, err := r.templatesClient.List(context.Background(), formtemplates.ListOptions{
		Namespace: namespace,
		Subject:   sub,
		Orgs:      orgs,
	})
	if err != nil {
		log.Err(err).Msg("unable to resolve form templates")
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
	if err := enc.Encode(all); err != nil {
		log.Err(err).Msg("unable to serve json encoded form template list")
	}
}

func (r *lister) complete() error {
	if r.templatesClient == nil {
		cli, err := formtemplates.NewClient(r.rc,
			formtemplates.AuthnNS(r.authnNS),
			formtemplates.Eval(true))
		if err != nil {
			return err
		}

		r.templatesClient = cli
	}

	return nil
}
