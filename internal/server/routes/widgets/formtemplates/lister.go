package formtemplates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	formtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/schemadefinitions"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/widgets/formtemplates"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	listerPath = "/apis/widgets.ui.krateo.io/formtemplates"
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
	rc                *rest.Config
	gr                schema.GroupResource
	templatesClient   *formtemplates.Client
	definitionsClient *schemadefinitions.Client
	authnNS           string
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

	if err := r.complete(); err != nil {
		log.Err(err).Msg("unable to initialize rest clients")
		encode.InternalError(wri, err)
		return
	}

	all, err := r.templatesClient.Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msg("unable to resolve form templates")
		if apierrors.IsNotFound(err) {
			encode.NotFound(wri, err)
		} else {
			encode.Invalid(wri, err)
		}
		return
	}

	for i := 0; i < len(all.Items); i++ {
		obj := &all.Items[i]
		formGVK, err := r.definitionsClient.Namespace(obj.Spec.SchemaDefinitionRef.Namespace).
			GVK(context.Background(), obj.Spec.SchemaDefinitionRef.Name)
		if err != nil {
			log.Err(err).Msg("unable to resolve form definition gvk")
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
			Namespace:     namespace,
		})
		if err != nil {
			log.Err(err).Msg("checking if 'get' verb is allowed")
			encode.InternalError(wri, err)
			return
		}

		if !ok {
			gr := dynamic.InferGroupResource(formGVK.Group, formGVK.Kind)
			encode.Forbidden(wri,
				fmt.Errorf("forbidden: User %q cannot get resource %q", sub, gr))
			return
		}

		sch, err := r.definitionsClient.OpenAPISchema(context.Background(), formGVK)
		if err != nil {
			log.Err(err).Msg("unable to resolve schema definition openAPI schema")
			encode.Invalid(wri, err)
			return
		}

		obj.Status.Content = &formtemplatesv1alpha1.FormTemplateStatusContent{
			//Instance: &runtime.RawExtension{Object: vals},
			Schema: &runtime.RawExtension{Object: sch},
		}
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(all); err != nil {
		log.Err(err).Msg("unable to serve json encoded form template list")
	}
}

func (r *lister) complete() error {
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
