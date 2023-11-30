package cardtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/krateoplatformops/krateo-bff/apis/core"
	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/resolvers"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	"k8s.io/utils/strings/slices"
)

const (
	getterPath = "/apis/widgets.ui.krateo.io/v1alpha1/cardtemplates/{name}"
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

	name := chi.URLParam(req, "name")

	qs := req.URL.Query()

	namespace := qs.Get("namespace")
	sub := qs.Get("sub")
	orgs := strings.Split(qs.Get("orgs"), ",")
	eval := true
	if qs.Has("eval") {
		ok, err := strconv.ParseBool(qs.Get("eval"))
		if err == nil {
			eval = ok
		}
	}

	ok, err := rbacutil.CanListResource(context.TODO(), r.rc, rbacutil.ResourceInfo{
		Subject: sub,
		Groups:  orgs,
		GroupResource: schema.GroupResource{
			Group: cardtemplatev1alpha1.Group, Resource: "cardtemplates",
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
			Bool("eval", eval).
			Msg("checking if 'get' verb is allowed")
		encode.Invalid(wri, err)
		return
	}

	if !ok {
		encode.Forbidden(wri,
			fmt.Errorf("forbidden: User %q cannot get resource \"cardtemplates/%s\" in API group \"widgets.ui.krateo.io\"", sub, name))
		return
	}

	log.Debug().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("name", name).
		Str("namespace", namespace).
		Bool("eval", eval).
		Msg("resolving card template")

	el, err := resolvers.CardTemplateGetOne(context.Background(), r.rc, &core.Reference{
		Name: name, Namespace: namespace,
	}, eval)
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("name", name).
			Str("namespace", namespace).
			Msg("unable to resolve card template")

		if apierrors.IsNotFound(err) {
			encode.NotFound(wri, err)
		} else {
			encode.Invalid(wri, err)
		}
		return
	}

	if el != nil {
		gr := cardtemplatev1alpha1.CardTemplateGroupVersionKind.GroupVersion().
			WithResource("cardtemplates").
			GroupResource()
		all, err := rbacutil.GetAllowedVerbs(context.TODO(), r.rc, util.ResourceInfo{
			Subject: sub, Groups: orgs,
			GroupResource: gr, ResourceName: el.GetName(),
			Namespace: el.GetNamespace(),
		})
		if err != nil {
			log.Err(err).
				Str("sub", sub).
				Strs("orgs", orgs).
				Str("gr", gr.String()).
				Str("name", name).
				Str("namespace", namespace).
				Msg("unable to resolve allowed verbs for sub in orgs")
			encode.Invalid(wri, err)
			return
		}

		m := el.GetAnnotations()
		if len(m) == 0 {
			m = map[string]string{}
		}
		m[allowedVerbsAnnotationKey] = strings.Join(all, ",")
		el.SetAnnotations(m)

		for _, x := range el.Spec.App.Actions {
			verb := strings.ToLower(ptr.Deref(x.Verb, ""))
			x.Enabled = ptr.To(slices.Contains(all, verb))
		}

		log.Debug().
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("name", el.GetName()).
			Str("namespace", namespace).
			Strs("verbs", all).
			Msg("successfully resolved allowed verbs for sub in orgs")
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(el); err != nil {
		log.Err(err).Msg("unable to serve json encoded cardtemplate")
	}
}
