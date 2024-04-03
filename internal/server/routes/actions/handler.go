package actions

import (
	"encoding/json"
	"net/http"

	"github.com/krateoplatformops/krateo-bff/internal/server/encode"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
)

const (
	apiPath = "/apis/actions"
)

func newHandler(rc *rest.Config, opts HandlerOptions) (string, http.HandlerFunc) {
	handler := &handler{
		x: &invoker{
			rc:            rc,
			authnNS:       opts.AuthnNS,
			kubeServerURL: opts.KubeServerURL,
			kubeProxyURL:  opts.KubeProxyURL,
		},
	}
	return apiPath, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*handler)(nil)

type handler struct {
	x *invoker
}

func (r *handler) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	rt, err := r.x.do(req)
	if err != nil {
		encode.Invalid(wri, err)
		return
	}

	wri.Header().Set("Content-Type", "application/json")
	wri.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(wri)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rt); err != nil {
		log.Err(err).Msg("unable to serve api call response")
	}
}
