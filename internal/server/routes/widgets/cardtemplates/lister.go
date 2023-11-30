package cardtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	listerPath                    = "/apis/widgets.ui.krateo.io/v1alpha1/cardtemplates"
	forbiddenAtClusterScopeMsgFmt = "forbidden: User %q cannot list resource \"cardtemplates\" in API group \"widgets.ui.krateo.io\" at cluster scope"
	forbiddenInNamespaceMsgFmt    = "forbidden: User %q cannot list resource \"cardtemplates\" in API group \"widgets.ui.krateo.io\" in namespace %s"
)

func newLister(rc *rest.Config) (string, http.HandlerFunc) {
	handler := &lister{rc: rc}
	return listerPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*lister)(nil)

type lister struct {
	rc *rest.Config
}

func (r *lister) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context()).With().Logger()

	//namespace := chi.URLParam(req, "namespace")

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
		Namespace: namespace,
	})
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("namespace", namespace).
			Bool("eval", eval).
			Msg("checking if 'get' verb is allowed")
		encode.Invalid(wri, err)
		return
	}

	if !ok {
		if len(namespace) > 0 {
			encode.Forbidden(wri, fmt.Errorf(forbiddenInNamespaceMsgFmt, sub, namespace))
		} else {
			encode.Forbidden(wri, fmt.Errorf(forbiddenAtClusterScopeMsgFmt, sub))
		}
		return
	}

	log.Debug().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("namespace", namespace).
		Bool("eval", eval).
		Msg("resolving card template list")

	res, err := resolvers.CardTemplateGetAll(context.Background(), r.rc, namespace, eval)
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("namespace", namespace).
			Msg("unable to resolve card templates")
		if apierrors.IsNotFound(err) {
			encode.NotFound(wri, err)
		} else {
			encode.Invalid(wri, err)
		}

		return
	}

	if res != nil {
		gr := cardtemplatev1alpha1.CardTemplateGroupVersionKind.GroupVersion().
			WithResource("cardtemplates").
			GroupResource()
		for _, el := range res.Items {
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
					Str("name", el.GetName()).
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
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		log.Err(err).Msg("unable to serve json encoded cardtemplates")
	}
}
