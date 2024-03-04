package formtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/schemadefinitions"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/formtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	getterPath = "/apis/widgets.ui.krateo.io/formtemplates/{name}"
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
	rc                *rest.Config
	gr                schema.GroupResource
	templatesClient   *formtemplates.Client
	definitionsClient *schemadefinitions.Client
	authnNS           string
}

func (r *getter) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context()).With().Logger()

	name := chi.URLParam(req, "name")

	qs := req.URL.Query()

	namespace := qs.Get("namespace")
	sub := qs.Get("sub")
	orgs := strings.Split(qs.Get("orgs"), ",")
	version := qs.Get("version")
	if len(version) == 0 {
		version = "v1alpha1"
	}

	if err := r.complete(); err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("name", name).
			Str("namespace", namespace).
			Msg("unable to initialize rest clients")

		encode.InternalError(wri, err)
		return
	}

	obj, err := r.templatesClient.Namespace(namespace).Get(context.Background(), name)
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

	formGVK, err := r.definitionsClient.Namespace(obj.Spec.SchemaDefinitionRef.Namespace).
		GVK(context.Background(), obj.Spec.SchemaDefinitionRef.Name)
	if err != nil {
		log.Err(err).
			Str("sub", sub).
			Strs("orgs", orgs).
			Str("name", name).
			Str("namespace", namespace).
			Msg("unable to resolve form definition gvk")
		if apierrors.IsNotFound(err) {
			encode.NotFound(wri, err)
		} else {
			encode.Invalid(wri, err)
		}
		return
	}

	ok, err := rbacutil.CanListResource(context.TODO(), r.rc, rbacutil.ResourceInfo{
		Subject:       sub,
		Groups:        orgs,
		GroupResource: dynamic.InferGroupResource(formGVK.Group, formGVK.Kind),
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
		gr := dynamic.InferGroupResource(formGVK.Group, formGVK.Kind)
		encode.Forbidden(wri,
			fmt.Errorf("forbidden: User %q cannot get resource %q", sub, gr))
		return
	}

	log.Debug().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("name", name).
		Str("namespace", namespace).
		Msg("resolving form template")

	sch, err := r.definitionsClient.OpenAPISchema(context.Background(), formGVK)
	if err != nil {
		log.Err(err).
			Str("namespace", namespace).
			Str("object", obj.GetName()).
			Msg("unable to resolve schema definition openAPI schema")

		encode.Invalid(wri, err)
		return
	}

	obj.Status.Content = &formtemplatesv1alpha1.FormTemplateStatusContent{
		//Instance: &runtime.RawExtension{Object: vals},
		Schema: &runtime.RawExtension{Object: sch},
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(obj); err != nil {
		log.Err(err).Msg("unable to serve json encoded form template")
	}
}

func (r *getter) complete() error {
	if r.templatesClient == nil {
		cli, err := formtemplates.NewClient(r.rc)
		if err != nil {
			return err
		}

		r.templatesClient = cli
	}

	if r.definitionsClient == nil {
		cli, err := schemadefinitions.NewClient(r.rc)
		if err != nil {
			return err
		}

		r.definitionsClient = cli
	}

	return nil
}
