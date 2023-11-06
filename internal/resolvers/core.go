package resolvers

import (
	"fmt"
	"strconv"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"k8s.io/client-go/rest"
)

func GetEndpoint(rc *rest.Config, ref *core.Reference) (*core.Endpoint, error) {
	values, err := GetSecretData(rc, ref)
	if err != nil {
		return nil, err
	}

	res := &core.Endpoint{}
	if v, ok := values["server"]; ok {
		res.Server = v
	} else {
		return res, fmt.Errorf("missed required attribute for endpoint: server")
	}

	if v, ok := values["token"]; ok {
		res.Token = v
	}

	if v, ok := values["username"]; ok {
		res.Username = v
	}

	if v, ok := values["password"]; ok {
		res.Password = v
	}

	if v, ok := values["certificate-authority-data"]; ok {
		res.CertificateAuthorityData = []byte(v)
	}

	if v, ok := values["client-key-data"]; ok {
		res.ClientCertificateData = []byte(v)
	}

	if v, ok := values["client-certificate-data"]; ok {
		res.ClientCertificateData = []byte(v)
	}

	if v, ok := values["debug"]; ok {
		res.Debug, _ = strconv.ParseBool(v)
	}

	return res, nil
}
