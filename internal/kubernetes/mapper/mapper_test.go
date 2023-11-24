package mapper

import (
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/clientcmd"
)

func TestFindGVR(t *testing.T) {
	const sample = `
	{
		"apiVersion":"v1",
		"items":[
		   {
			  "apiVersion":"v1",
			  "data":{
				 "password":"MWYyZDFlMmU2N2Rm",
				 "username":"YWRtaW4="
			  },
			  "kind":"Secret",
			  "metadata":{
				 "annotations":{
					"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"data\":{\"password\":\"MWYyZDFlMmU2N2Rm\",\"username\":\"YWRtaW4=\"},\"kind\":\"Secret\",\"metadata\":{\"annotations\":{},\"name\":\"mysecret\",\"namespace\":\"dev-system\"},\"type\":\"Opaque\"}\n"
				 },
				 "creationTimestamp":"2023-11-20T09:18:14Z",
				 "name":"mysecret",
				 "namespace":"dev-system",
				 "resourceVersion":"680",
				 "uid":"51ce28e8-8775-468a-891a-3a9264f65f1a"
			  },
			  "type":"Opaque"
		   },
		   {
			  "apiVersion":"v1",
			  "data":{
				 "password":"MWYyZDFlMmU2N2Rm",
				 "username":"YWRtaW4="
			  },
			  "kind":"Secret",
			  "metadata":{
				 "annotations":{
					"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"data\":{\"password\":\"MWYyZDFlMmU2N2Rm\",\"username\":\"YWRtaW4=\"},\"kind\":\"Secret\",\"metadata\":{\"annotations\":{},\"name\":\"one\",\"namespace\":\"dev-system\"},\"type\":\"Opaque\"}\n"
				 },
				 "creationTimestamp":"2023-11-20T09:27:11Z",
				 "name":"one",
				 "namespace":"dev-system",
				 "resourceVersion":"1376",
				 "uid":"a3b6f87f-49ed-4173-87fb-4da01aefce02"
			  },
			  "type":"Opaque"
		   },
		   {
			  "apiVersion":"v1",
			  "data":{
				 "password":"enppbWllaSE=",
				 "username":"cGlwcG8="
			  },
			  "kind":"Secret",
			  "metadata":{
				 "annotations":{
					"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Secret\",\"metadata\":{\"annotations\":{},\"name\":\"two\",\"namespace\":\"dev-system\"},\"stringData\":{\"password\":\"zzimiei!\",\"username\":\"pippo\"},\"type\":\"Opaque\"}\n"
				 },
				 "creationTimestamp":"2023-11-20T09:27:11Z",
				 "name":"two",
				 "namespace":"dev-system",
				 "resourceVersion":"1377",
				 "uid":"990cefe2-782b-4f30-9f9d-bfc761ade0b4"
			  },
			  "type":"Opaque"
		   }
		],
		"kind":"List",
		"metadata":{
		   "resourceVersion":""
		}
	 }`

	/*
		qry := `.items[] | "\(.apiVersion),\(.kind)"`

		buf := bytes.Buffer{}
		err := chisel.Polish(strings.NewReader(sample), &buf, qry)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(buf.String())
	*/

	kubeconfig, err := os.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating clientConfig")

	restConfig, err := clientConfig.ClientConfig()
	assert.Nil(t, err, "expecting nil error getting restConfig")

	list := unstructured.UnstructuredList{}
	if err := list.UnmarshalJSON([]byte(sample)); err != nil {
		t.Fatal(err)
	}

	// copyList := unstructured.UnstructuredList{
	// 	Object: list.Object,
	// 	Items:  make([]unstructured.Unstructured, 0, len(list.Items)),
	// }

	for _, el := range list.Items {
		gvr, err := FindGVR(restConfig, el.GroupVersionKind().GroupKind())
		if err != nil {
			t.Fatal(err)
		}

		all, err := rbac.AllowedVerbsOnResourceForSubject(restConfig, pkix.Name{
			CommonName: "luca", Organization: []string{"devs"},
		}, gvr.GroupResource(), "", "dev-system")
		if err != nil {
			t.Fatal(err)
		}

		m := el.GetAnnotations()
		if len(m) == 0 {
			m = map[string]string{}
		}
		m["krateo.io/allowed-verbs"] = strings.Join(all, ",")
		el.SetAnnotations(m)

		// copyObj := el.DeepCopy()
		// m := copyObj.GetAnnotations()
		// if len(m) == 0 {
		// 	m = map[string]string{}
		// }
		// m["krateo.io/allowed-verbs"] = strings.Join(all, ",")
		// copyObj.SetAnnotations(m)

		// copyList.Items = append(copyList.Items, *copyObj)
		//dat, err := getPatchData(el, copyObj)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//_, err = f.dynamicClient.Resource(getDeploymentGVR()).Namespace(obj.GetNamespace()).Patch(f.ctx, obj.GetName(), types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})

		//spew.Dump(all)
	}
	//spew.Dump(list)

	dat, err := list.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(dat))

}

// getPatchData will return difference between original and modified document
func getPatchData(originalObj, modifiedObj interface{}) ([]byte, error) {
	originalData, err := json.Marshal(originalObj)
	if err != nil {
		return nil, errors.Wrapf(err, "failed marshal original data")
	}
	modifiedData, err := json.Marshal(modifiedObj)
	if err != nil {
		return nil, errors.Wrapf(err, "failed marshal original data")
	}

	// Using strategicpatch package can cause below error
	// Error: CreateTwoWayMergePatch failed: unable to find api field in struct Unstructured for the json field "spec"
	//patchBytes, err := strategicpatch.CreateTwoWayMergePatch(originalData, modifiedData, originalObj)
	// if err != nil {
	// 	return nil, errors.Errorf("CreateTwoWayMergePatch failed: %v", err)
	// }

	patchBytes, err := jsonpatch.CreateMergePatch(originalData, modifiedData)
	if err != nil {
		return nil, errors.Errorf("CreateTwoWayMergePatch failed: %v", err)
	}
	return patchBytes, nil
}
