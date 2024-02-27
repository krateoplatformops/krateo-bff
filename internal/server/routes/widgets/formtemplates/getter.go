package formtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/formtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/formtemplates/evaluator"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	getterPath = "/apis/widgets.ui.krateo.io/v1alpha1/formtemplates/{name}"
)

func newGetter(rc *rest.Config, authnNS string) (string, http.HandlerFunc) {
	gr := formtemplatesv1alpha1.FormTemplateGroupVersionKind.GroupVersion().
		WithResource("formtemplates").
		GroupResource()

	handler := &getter{
		rc:      rc,
		gr:      gr,
		authnNS: authnNS,
	}
	return getterPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*getter)(nil)

type getter struct {
	rc      *rest.Config
	client  *formtemplates.Client
	gr      schema.GroupResource
	authnNS string
}

func (r *getter) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context()).With().Logger()

	name := chi.URLParam(req, "name")

	qs := req.URL.Query()

	namespace := qs.Get("namespace")
	sub := qs.Get("sub")
	orgs := strings.Split(qs.Get("orgs"), ",")

	ok, err := rbacutil.CanListResource(context.TODO(), r.rc, rbacutil.ResourceInfo{
		Subject:       sub,
		Groups:        orgs,
		GroupResource: r.gr,
		ResourceName:  name,
		Namespace:     namespace,
	})
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("name", name).
			Str("namespace", namespace).
			Msg("checking if 'get' verb is allowed")
		encode.InternalError(wri, err)
		return
	}

	if !ok {
		encode.Forbidden(wri,
			fmt.Errorf("forbidden: User %q cannot get resource \"formtemplates/%s\" in API group \"widgets.ui.krateo.io\"", sub, name))
		return
	}

	log.Debug().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("name", name).
		Str("namespace", namespace).
		Msg("resolving form template")

	if r.client == nil {
		cli, err := formtemplates.NewClient(r.rc)
		if err != nil {
			log.Err(err).
				Str("sub", sub).
				Strs("orgs", orgs).
				Str("name", name).
				Str("namespace", namespace).
				Msg("unable to create form template rest client")

			encode.InternalError(wri, err)
			return
		}

		r.client = cli
	}

	obj, err := r.client.Namespace(namespace).Get(context.TODO(), name)
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("name", name).
			Str("namespace", namespace).
			Msg("unable to resolve form template")

		if apierrors.IsNotFound(err) {
			encode.NotFound(wri, err)
		} else {
			encode.Invalid(wri, err)
		}
		return
	}

	err = evaluator.Eval(context.Background(), obj, evaluator.EvalOptions{
		RESTConfig: r.rc, AuthnNS: r.authnNS, Subject: sub, Groups: orgs,
	})
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("name", name).
			Str("namespace", namespace).
			Str("object", obj.GetName()).
			Msg("unable to evaluate form template")

		encode.Invalid(wri, err)
		return
	}

	// log.Debug().
	// 	Str("name", obj.GetName()).
	// 	Str("namespace", namespace).
	// 	Strs("verbs", verbs).
	// 	Msg("successfully resolved allowed verbs for sub in orgs")

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(obj); err != nil {
		log.Err(err).Msg("unable to serve json encoded form template")
	}
}