package rows

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/rows"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/rows/evaluator"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	getterPath         = "/apis/layout.ui.krateo.io/rows/{name}"
	forbiddenGetMsgFmt = "forbidden: User %q cannot get resource %q in API group \"layout.ui.krateo.io\" in namespace %s"
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
	client  *rows.Client
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
		encode.Forbidden(wri, fmt.Errorf(forbiddenGetMsgFmt, sub, name, namespace))
		return
	}

	if r.client == nil {
		cli, err := rows.NewClient(r.rc)
		if err != nil {
			log.Err(err).Msg("unable to create rows rest client")
			encode.InternalError(wri, err)
			return
		}

		r.client = cli
	}

	obj, err := r.client.Namespace(namespace).Get(context.TODO(), name)
	if err != nil {
		log.Err(err).Msg("unable to resolve row")
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
			Str("object", obj.GetName()).
			Msg("unable to evaluate row")

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
				Str("object", obj.GetName()).
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
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(obj); err != nil {
		log.Err(err).Msg("unable to serve json encoded row")
	}
}
