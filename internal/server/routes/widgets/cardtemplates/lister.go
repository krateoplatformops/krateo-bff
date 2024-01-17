package cardtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/cardtemplates/evaluator"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func newLister(rc *rest.Config, authnNS string) (string, http.HandlerFunc) {
	gr := cardtemplatev1alpha1.CardTemplateGroupVersionKind.GroupVersion().
		WithResource("cardtemplates").
		GroupResource()
	handler := &lister{rc: rc, authnNS: authnNS, gr: gr}
	return listerPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*lister)(nil)

type lister struct {
	rc      *rest.Config
	client  *cardtemplates.Client
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
			Msg("checking if 'get' verb is allowed")
		encode.InternalError(wri, err)
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
		Msg("resolving card template list")

	if r.client == nil {
		cli, err := cardtemplates.NewClient(r.rc)
		if err != nil {
			log.Err(err).
				Str("sub", sub).
				Strs("orgs", orgs).
				Str("namespace", namespace).
				Msg("unable to create card template rest client")

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
			Msg("unable to list card template")

		if apierrors.IsNotFound(err) {
			encode.NotFound(wri, err)
		} else {
			encode.Invalid(wri, err)
		}
		return
	}

	for _, el := range all.Items {
		obj := &el
		err = evaluator.Eval(context.Background(), obj, evaluator.EvalOptions{
			RESTConfig: r.rc, AuthnNS: r.authnNS, Username: sub,
		})
		if err != nil {
			log.Err(err).Str("object", obj.GetName()).
				Msg("unable to evaluate card template")

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

		obj.Status.AllowedActions = []string{}
		for _, x := range obj.Spec.App.Actions {
			verb := strings.ToLower(ptr.Deref(x.Verb, ""))

			if slices.Contains(verbs, verb) {
				obj.Status.AllowedActions = append(obj.Status.AllowedActions, x.Name)
			}
		}
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(all); err != nil {
		log.Err(err).Msg("unable to serve json encoded cardtemplates")
	}
}
