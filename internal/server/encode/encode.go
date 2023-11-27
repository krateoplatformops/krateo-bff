package encode

import (
	"encoding/json"
	"net/http"

	cardtemplatev1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplate/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Error(w http.ResponseWriter, reason metav1.StatusReason, code int, err error) error {
	out := metav1.Status{
		Status: "Failure",
		Reason: reason,
		Code:   int32(code),
		Details: &metav1.StatusDetails{
			Group: cardtemplatev1alpha1.Group,
			Kind:  cardtemplatev1alpha1.CardTemplateKind,
		},
		Message: err.Error(),
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&out)
}
