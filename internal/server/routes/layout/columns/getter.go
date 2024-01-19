package columns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	columnsv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/columns/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/columns"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/columns/evaluator"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	getterPath         = "/apis/layout.ui.krateo.io/v1alpha1/columns/{name}"
	forbiddenGetMsgFmt = "forbidden: User %q cannot get resource %q in API group \"layout.ui.krateo.io\" in namespace %s"
)

func newGetter(rc *rest.Config, authnNS string) (string, http.HandlerFunc) {
	gr := columnsv1alpha1.ColumnGroupVersionKind.GroupVersion().
		WithResource(resources).
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
	client  *columns.Client
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
		Subject: sub,
		Groups:  orgs,
		GroupResource: schema.GroupResource{
			Group: columnsv1alpha1.Group, Resource: resources,
		},
		ResourceName: name,
		Namespace:    namespace,
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
		encode.Forbidden(wri, fmt.Errorf(forbiddenGetMsgFmt, sub, name, namespace))
		return
	}

	log.Debug().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("name", name).
		Str("namespace", namespace).
		Msg("resolving columns")

	if r.client == nil {
		cli, err := columns.NewClient(r.rc)
		if err != nil {
			log.Err(err).
				Str("sub", sub).
				Strs("orgs", orgs).
				Str("name", name).
				Str("namespace", namespace).
				Msg("unable to create columns rest client")

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
			Msg("unable to resolve column")

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
			Msg("unable to evaluate column")

		encode.Invalid(wri, err)
		return
	}

	if obj != nil {
		verbs, err := rbacutil.GetAllowedVerbs(context.TODO(), r.rc, util.ResourceInfo{
			Subject: sub, Groups: orgs,
			GroupResource: r.gr, ResourceName: obj.GetName(),
			Namespace: obj.GetNamespace(),
		})
		if err != nil {
			log.Err(err).
				Str("name", name).
				Str("namespace", namespace).
				Msg("unable to resolve allowed verbs")
			encode.Invalid(wri, err)
			return
		}

		m := obj.GetAnnotations()
		if len(m) == 0 {
			m = map[string]string{}
		}
		m[allowedVerbsAnnotationKey] = strings.Join(verbs, ",")
		obj.SetAnnotations(m)

		// _, err = r.client.Namespace(namespace).UpdateStatus(context.TODO(), obj)
		// if err != nil {
		// 	log.Err(err).Str("object", obj.GetName()).Msg("unable to update object status")
		// }
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(obj); err != nil {
		log.Err(err).Msg("unable to serve json encoded column")
	}
}
