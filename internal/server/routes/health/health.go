package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
)

const (
	getterPathFmt = "/apis/health"
)

type Options struct {
	Version     string
	Build       string
	ServiceName string
	Healty      *int32
}

func Register(r *chi.Mux, opts Options) {
	r.Get(healthGetter(opts))
}

func healthGetter(opts Options) (string, http.HandlerFunc) {
	handler := &getter{
		healthy:     opts.Healty,
		version:     fmt.Sprintf("ver: %s (bld: %s)", opts.Version, opts.Build),
		serviceName: opts.ServiceName,
	}
	return getterPathFmt, func(wri http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(wri, req)
	}
}

var _ http.Handler = (*getter)(nil)

type getter struct {
	healthy     *int32
	version     string
	serviceName string
}

func (r *getter) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	if atomic.LoadInt32(r.healthy) == 1 {
		data := map[string]string{
			"name":    r.serviceName,
			"version": r.version,
		}

		wri.Header().Set("Content-Type", "application/json")
		wri.WriteHeader(http.StatusOK)
		json.NewEncoder(wri).Encode(data)
		return
	}
	wri.WriteHeader(http.StatusServiceUnavailable)
}
