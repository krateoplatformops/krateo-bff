package apps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/server/decode"
	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	applierPath = "/apis/apps/{name}"
)

func newApplier(rc *rest.Config, authnNS string) (string, http.HandlerFunc) {
	handler := &applier{
		rc:      rc,
		authnNS: authnNS,
	}
	return applierPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*applier)(nil)

type applier struct {
	rc      *rest.Config
	authnNS string
	client  *dynamic.Applier
}

func (r *applier) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")

	qs := req.URL.Query()

	namespace := qs.Get("namespace")
	sub := qs.Get("sub")
	orgs := strings.Split(qs.Get("orgs"), ",")
	version := qs.Get("version")
	if len(version) == 0 {
		version = "v1alpha1"
	}
	group := qs.Get("group")
	kind := qs.Get("kind")

	gr := dynamic.InferGroupResource(group, kind)

	log := zerolog.Ctx(req.Context()).With().
		Str("sub", sub).
		Strs("orgs", orgs).
		Str("name", name).
		Str("namespace", namespace).
		Str("version", version).
		Str("gr", gr.String()).
		Logger()

	ok, err := rbacutil.CanCreateOrUpdateResource(context.Background(), r.rc, rbacutil.ResourceInfo{
		Subject:       sub,
		Groups:        orgs,
		GroupResource: gr,
		ResourceName:  name,
		Namespace:     namespace,
	})
	if err != nil {
		log.Err(err).Msg("checking if [create,update] verbs are allowed")
		encode.InternalError(wri, err)
		return
	}
	if !ok {
		encode.Forbidden(wri,
			fmt.Errorf("forbidden: User %q cannot create or update resources %s", sub, gr.String()))
		return
	}

	if r.client == nil {
		cli, err := dynamic.NewApplier(r.rc)
		if err != nil {
			log.Err(err).Msg("unable to create resource applier client")
			encode.InternalError(wri, err)
			return
		}

		r.client = cli
	}

	content := map[string]any{}
	err = decode.JSONBody(wri, req, &content)
	if err != nil {
		log.Err(err).Msg("decoding JSON data")
		mr := &decode.MalformedRequest{}
		if errors.As(err, &mr) {
			encode.BadRequest(wri, err)
		} else {
			encode.InternalError(wri, err)
		}
		return
	}

	err = r.client.Apply(context.Background(), map[string]any{"spec": content}, dynamic.ApplyOptions{
		GVK: schema.GroupVersionKind{
			Group:   group,
			Version: version,
			Kind:    kind,
		},
		Name:      name,
		Namespace: namespace,
	})
	if err != nil {
		log.Err(err).Any("content", content).
			Msg("creating or updating resource")
		encode.InternalError(wri, err)
		return
	}

	encode.OK(wri, 200)
}
