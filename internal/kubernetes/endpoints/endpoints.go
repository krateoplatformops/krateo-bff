package endpoints

import (
	"context"
	"fmt"
	"strconv"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/secrets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func Resolve(ctx context.Context, rc *rest.Config, ref *core.Reference) (*core.Endpoint, error) {
	//bcl, err := builtins.NewForConfig(rc)
	cli, err := secrets.NewClient(rc)
	if err != nil {
		return nil, err
	}

	sec, err := cli.Namespace(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	res := &core.Endpoint{}
	if v, ok := sec.Data["server-url"]; ok {
		res.ServerURL = string(v)
	} else {
		return res, fmt.Errorf("missed required attribute for endpoint: server-url")
	}

	if v, ok := sec.Data["proxy-url"]; ok {
		res.ProxyURL = string(v)
	}

	if v, ok := sec.Data["token"]; ok {
		res.Token = string(v)
	}

	if v, ok := sec.Data["username"]; ok {
		res.Username = string(v)
	}

	if v, ok := sec.Data["password"]; ok {
		res.Password = string(v)
	}

	if v, ok := sec.Data["certificate-authority-data"]; ok {
		res.CertificateAuthorityData = string(v)
	}

	if v, ok := sec.Data["client-key-data"]; ok {
		res.ClientKeyData = string(v)
	}

	if v, ok := sec.Data["client-certificate-data"]; ok {
		res.ClientCertificateData = string(v)
	}

	if v, ok := sec.Data["debug"]; ok {
		res.Debug, _ = strconv.ParseBool(string(v))
	}

	return res, nil
}
