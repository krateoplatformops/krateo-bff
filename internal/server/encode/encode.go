package encode

import (
	"encoding/json"
	"net/http"
)

func Success(w http.ResponseWriter, dat []byte) error {
	out := response{
		Code: http.StatusOK,
		Data: dat,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&out)
}

func Error(w http.ResponseWriter, status int, err error) error {
	out := response{
		Code:  status,
		Error: err.Error(),
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&out)
}

type response struct {
	Code  int             `json:"code"`
	Error string          `json:"error,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
}
