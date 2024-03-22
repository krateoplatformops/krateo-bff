package cardtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	getterPath = "/apis/widgets.ui.krateo.io/cardtemplates/{name}"
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
	rc      *rest.Config
	client  *cardtemplates.Client
	gr      schema.GroupResource
	authnNS string
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

	ok, err := rbacutil.CanListResource(context.TODO(), r.rc, rbacutil.ResourceInfo{
		Subject:       sub,
		Groups:        orgs,
		GroupResource: r.gr,
		ResourceName:  name,
		Namespace:     namespace,
	})
	if err != nil {
		log.Err(err).Msg("checking if 'get' verb is allowed")
		encode.InternalError(wri, err)
		return
	}

	if !ok {
		encode.Forbidden(wri,
			fmt.Errorf("forbidden: User %q cannot get resource \"cardtemplates/%s\" in API group \"widgets.ui.krateo.io\"", sub, name))
		return
	}

	if r.client == nil {
		cli, err := cardtemplates.NewClient(r.rc,
			cardtemplates.AuthnNS(r.authnNS),
			cardtemplates.Eval(true))
		if err != nil {
			log.Err(err).Msg("unable to create cardtemplates rest client")
			encode.InternalError(wri, err)
			return
		}

		r.client = cli
	}

	obj, err := r.client.Get(context.TODO(), cardtemplates.GetOptions{
		Namespace: namespace,
		Name:      name,
		Subject:   sub,
		Orgs:      orgs,
	})
	if err != nil {
		log.Err(err).Msg("unable to resolve card template")
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
		log.Err(err).Msg("unable to serve json encoded cardtemplate")
	}
}
