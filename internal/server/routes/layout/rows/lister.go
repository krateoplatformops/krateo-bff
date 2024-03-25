package rows

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/layout/rows"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	listerPath                        = "/apis/layout.ui.krateo.io/rows"
	forbiddenListAtClusterScopeMsgFmt = "forbidden: User %q cannot list resource \"rows\" in API group \"layout.ui.krateo.io\" at cluster scope"
	forbiddenListInNamespaceMsgFmt    = "forbidden: User %q cannot list resource \"rows\" in API group \"layout.ui.krateo.io\" in namespace %s"
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
	client  *rows.Client
	gr      schema.GroupResource
	authnNS string
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

	ok, err := rbacutil.CanListResource(context.TODO(), r.rc, rbacutil.ResourceInfo{
		Subject:       sub,
		Groups:        orgs,
		GroupResource: r.gr,
		Namespace:     namespace,
	})
	if err != nil {
		log.Err(err).Msg("checking if 'get' verb is allowed")
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

	if r.client == nil {
		cli, err := rows.NewClient(r.rc, true)
		if err != nil {
			log.Err(err).Msg("unable to create rows rest client")
			encode.InternalError(wri, err)
			return
		}

		r.client = cli
	}

	all, err := r.client.List(context.Background(), rows.ListOptions{
		Namespace: namespace,
		Subject:   sub,
		Orgs:      orgs,
		AuthnNS:   r.authnNS,
	})
	if err != nil {
		log.Err(err).Msg("unable to list rows")
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
		log.Err(err).Msg("unable to serve json encoded column list")
	}
}
