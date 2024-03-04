package columns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/columns"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/columns/evaluator"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	listerPath                        = "/apis/layout.ui.krateo.io/columns"
	forbiddenListAtClusterScopeMsgFmt = "forbidden: User %q cannot list resource \"columns\" in API group \"layout.ui.krateo.io\" at cluster scope"
	forbiddenListInNamespaceMsgFmt    = "forbidden: User %q cannot list resource \"columns\" in API group \"layout.ui.krateo.io\" in namespace %s"
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
	rc      *rest.Config
	client  *columns.Client
	gr      schema.GroupResource
	authnNS string
}

func (r *lister) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context()).With().Logger()

	qs := req.URL.Query()

	namespace := qs.Get("namespace")
	sub := qs.Get("sub")
	orgs := strings.Split(qs.Get("orgs"), ",")

	ok, err := rbacutil.CanListResource(context.TODO(), r.rc, rbacutil.ResourceInfo{
		Subject:       sub,
		Groups:        orgs,
		GroupResource: r.gr,
		Namespace:     namespace,
	})
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("namespace", namespace).
			Msg("checking if 'get' verb is allowed")
		encode.InternalError(wri, err)
		return
	}

	if !ok {
		if len(namespace) > 0 {
			encode.Forbidden(wri, fmt.Errorf(forbiddenListInNamespaceMsgFmt, sub, namespace))
		} else {
			encode.Forbidden(wri, fmt.Errorf(forbiddenListAtClusterScopeMsgFmt, sub))
		}
		return
	}

	log.Debug().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("namespace", namespace).
		Msg("resolving column list")

	if r.client == nil {
		cli, err := columns.NewClient(r.rc)
		if err != nil {
			log.Err(err).
				Str("sub", sub).
				Strs("orgs", orgs).
				Str("namespace", namespace).
				Msg("unable to create column rest client")

			encode.InternalError(wri, err)
			return
		}

		r.client = cli
	}

	all, err := r.client.Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("namespace", namespace).
			Msg("unable to list columns")

		if apierrors.IsNotFound(err) {
			encode.NotFound(wri, err)
		} else {
			encode.Invalid(wri, err)
		}
		return
	}

	for i, el := range all.Items {
		obj := &el
		err = evaluator.Eval(context.Background(), obj, evaluator.EvalOptions{
			RESTConfig: r.rc, AuthnNS: r.authnNS, Subject: sub, Groups: orgs,
		})
		if err != nil {
			log.Err(err).Str("object", obj.GetName()).
				Msg("unable to evaluate column")

			encode.Invalid(wri, err)
			return
		}

		verbs, err := rbacutil.GetAllowedVerbs(context.TODO(), r.rc, util.ResourceInfo{
			Subject: sub, Groups: orgs,
			GroupResource: r.gr, ResourceName: obj.GetName(),
			Namespace: obj.GetNamespace(),
		})
		if err != nil {
			log.Err(err).
				Str("name", obj.GetName()).
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

		all.Items[i] = *obj
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(all); err != nil {
		log.Err(err).Msg("unable to serve json encoded column list")
	}
}
