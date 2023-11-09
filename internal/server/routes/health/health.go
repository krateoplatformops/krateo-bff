package health

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

func Check(healthy *int32, version string, serviceName string) http.Handler {
	return &healthRoute{
		healthy:     healthy,
		version:     version,
		serviceName: serviceName,
	}
}

var _ http.Handler = (*healthRoute)(nil)

type healthRoute struct {
	healthy     *int32
	version     string
	serviceName string
}

func (r *healthRoute) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
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
