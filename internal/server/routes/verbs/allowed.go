package verbs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	getterPath = "/apis/allowed-verbs"
)

func newGetter(rc *rest.Config) (string, http.HandlerFunc) {
	handler := &getter{rc: rc}
	return getterPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*getter)(nil)

type getter struct {
	rc *rest.Config
}

func (r *getter) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context()).With().Logger()

	qs := req.URL.Query()
	namespace := qs.Get("namespace")
	name := qs.Get("name")

	if len(qs.Get("gr")) == 0 {
		encode.BadRequest(wri, fmt.Errorf("missing parameter 'gr'"))
		return
	}
	gr := schema.ParseGroupResource(qs.Get("gr"))

	sub := qs.Get("sub")
	orgs := strings.Split(qs.Get("orgs"), ",")
	if len(sub) == 0 && len(orgs) == 0 {
		encode.BadRequest(wri, fmt.Errorf("both parameters 'sub' and 'orgs' cannot be empty"))
		return
	}

	log.Debug().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("name", name).
		Str("namespace", namespace).
		Str("gr", gr.String()).
		Msg("resolving allowed verbs")

	all, err := rbacutil.GetAllowedVerbs(context.TODO(), r.rc,
		util.GetAllowedVerbsOption{
			Subject: sub, Groups: orgs,
			GroupResource: gr,
			ResourceName:  name,
			Namespace:     namespace,
		})
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("gr", gr.String()).
			Str("name", name).
			Str("namespace", namespace).
			Msg("unable to resolve allowed verbs")
		encode.Invalid(wri, err)
		return
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(all); err != nil {
		log.Err(err).Msg("unable to serve json encoded allowed verbs")
	}
}
