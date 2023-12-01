# Krateo BFF

Krateo **B**ackend **F**or **F**rontend.

## CardTemplate API

### List all card templates

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Verb**          | `GET`                                                                               |
| **Path**          | `/apis/widgets.ui.krateo.io/v1alpha1/cardtemplates`                                 |
| **Query Params**  | `eval`: if _"false"_ do not evaluate expressions (default: true)                    |
|                   | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates (optional)                       |

**Example**

```sh
curl "http://localhost:8080/apis/widgets.ui.krateo.io/v1alpha1/cardtemplates?sub=cyberjoker&orgs=devs&namespace=dev-system"
```

**Response On Failure**

```json
{
  "kind": "Status",
  "apiVersion": "v1",
  "status": "Failure",
  "message": "forbidden: User \"cyberjoker\" cannot list resource \"cardtemplates\" in API group \"widgets.ui.krateo.io\" in namespace dev-system",
  "reason": "Forbidden",
  "code": 403
}
```

**Response On Success**

```json
{
  "kind": "CardTemplateList",
  "apiVersion": "widgets.ui.krateo.io/v1alpha1",
  "metadata": {
    "resourceVersion": "2618"
  },
  "items": [
    {
      "kind": "CardTemplate",
      "apiVersion": "widgets.ui.krateo.io/v1alpha1",
      "metadata": {
        "name": "card-dev",
        "namespace": "dev-system",
        "uid": "2f05f33c-f3c8-4ed1-b39e-f398524bb33d",
        "resourceVersion": "275",
        "generation": 1,
        "creationTimestamp": "2023-12-01T15:37:18Z",
        "annotations": {
          "krateo.io/allowed-verbs": "get,list",
          "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"widgets.ui.krateo.io/v1alpha1\",\"kind\":\"CardTemplate\",\"metadata\":{\"annotations\":{},\"name\":\"card-dev\",\"namespace\":\"dev-system\"},\"spec\":{\"api\":[{\"endpointRef\":{\"name\":\"httpbin-endpoint\",\"namespace\":\"dev-system\"},\"headers\":[\"Accept: application/json\"],\"name\":\"httpbin\",\"path\":\"/anything\",\"payload\":\"{\\n  \\\"firstName\\\": \\\"Charles\\\",\\n  \\\"lastName\\\": \\\"Doe\\\",\\n  \\\"age\\\": 41,\\n  \\\"location\\\": {\\n    \\\"city\\\": \\\"San Fracisco\\\",\\n    \\\"postalCode\\\": \\\"94103\\\"\\n  },\\n  \\\"hobbies\\\": [\\n    \\\"chess\\\",\\n    \\\"netflix\\\"\\n  ]\\n}\\n\",\"verb\":\"GET\"}],\"app\":{\"actions\":[{\"enabled\":true,\"endpointRef\":{\"name\":\"krateogateway-endpoint\",\"namespace\":\"krateo-system\"},\"name\":\"crd\",\"path\":\"/apis/apiextensions.k8s.io/v1/customresourcedefinitions/postgresqls.composition.krateo.io\",\"verb\":\"DELETE\"}],\"color\":\"GREEN\",\"content\":\"${ .httpbin.json.location.city }\",\"icon\":\"ApartmentOutlined\",\"tags\":\"${ .httpbin.json.hobbies | join(\\\",\\\") }\",\"title\":\"${ (.httpbin.json.firstName  + \\\" \\\" + .httpbin.json.lastName) }\"}}}\n"
        },
        "managedFields": [
          {
            "manager": "kubectl-client-side-apply",
            "operation": "Update",
            "apiVersion": "widgets.ui.krateo.io/v1alpha1",
            "time": "2023-12-01T15:37:18Z",
            "fieldsType": "FieldsV1",
            "fieldsV1": {
              "f:metadata": {
                "f:annotations": {
                  ".": {},
                  "f:kubectl.kubernetes.io/last-applied-configuration": {}
                }
              },
              "f:spec": {
                ".": {},
                "f:api": {},
                "f:app": {
                  ".": {},
                  "f:actions": {},
                  "f:color": {},
                  "f:content": {},
                  "f:icon": {},
                  "f:tags": {},
                  "f:title": {}
                }
              }
            }
          }
        ]
      },
      "spec": {
        "app": {
          "title": "Charles Doe",
          "content": "San Fracisco",
          "icon": "ApartmentOutlined",
          "color": "GREEN",
          "tags": "chess,netflix",
          "actions": [
            {
              "name": "crd",
              "path": "/apis/apiextensions.k8s.io/v1/customresourcedefinitions/postgresqls.composition.krateo.io",
              "verb": "DELETE",
              "endpointRef": {
                "name": "krateogateway-endpoint",
                "namespace": "krateo-system"
              },
              "enabled": false
            }
          ]
        },
        "api": [
          {
            "name": "httpbin",
            "path": "/anything",
            "verb": "GET",
            "headers": [
              "Accept: application/json"
            ],
            "payload": "{\n  \"firstName\": \"Charles\",\n  \"lastName\": \"Doe\",\n  \"age\": 41,\n  \"location\": {\n    \"city\": \"San Fracisco\",\n    \"postalCode\": \"94103\"\n  },\n  \"hobbies\": [\n    \"chess\",\n    \"netflix\"\n  ]\n}\n",
            "endpointRef": {
              "name": "httpbin-endpoint",
              "namespace": "dev-system"
            },
            "enabled": true
          }
        ]
      }
    }
  ]
}
```

### Get one card template

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Method**        | `GET`                                                                               |
| **Path**          | `/apis/widgets.ui.krateo.io/v1alpha1/cardtemplates/${NAME}` |
| **Path Params**   | `${NAME}` is your `CardTemplate` object name                                        |
| **Query Params**  | `eval`: if _"false"_ do not evaluate expressions (default: true)                    |
|                   | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates                                  |

**Example**:

```sh
curl "http://localhost:8080/apis/widgets.ui.krateo.io/v1alpha1/cardtemplates/card-dev?sub=cyberjoker&orgs=devs&namespace=dev-system"
```

```json
{
  "kind": "CardTemplate",
  "apiVersion": "widgets.ui.krateo.io/v1alpha1",
  "metadata": {
    "name": "card-dev",
    "namespace": "dev-system",
    "uid": "2f05f33c-f3c8-4ed1-b39e-f398524bb33d",
    "resourceVersion": "275",
    "generation": 1,
    "creationTimestamp": "2023-12-01T15:37:18Z",
    "annotations": {
      "krateo.io/allowed-verbs": "get,list",
      "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"widgets.ui.krateo.io/v1alpha1\",\"kind\":\"CardTemplate\",\"metadata\":{\"annotations\":{},\"name\":\"card-dev\",\"namespace\":\"dev-system\"},\"spec\":{\"api\":[{\"endpointRef\":{\"name\":\"httpbin-endpoint\",\"namespace\":\"dev-system\"},\"headers\":[\"Accept: application/json\"],\"name\":\"httpbin\",\"path\":\"/anything\",\"payload\":\"{\\n  \\\"firstName\\\": \\\"Charles\\\",\\n  \\\"lastName\\\": \\\"Doe\\\",\\n  \\\"age\\\": 41,\\n  \\\"location\\\": {\\n    \\\"city\\\": \\\"San Fracisco\\\",\\n    \\\"postalCode\\\": \\\"94103\\\"\\n  },\\n  \\\"hobbies\\\": [\\n    \\\"chess\\\",\\n    \\\"netflix\\\"\\n  ]\\n}\\n\",\"verb\":\"GET\"}],\"app\":{\"actions\":[{\"enabled\":true,\"endpointRef\":{\"name\":\"krateogateway-endpoint\",\"namespace\":\"krateo-system\"},\"name\":\"crd\",\"path\":\"/apis/apiextensions.k8s.io/v1/customresourcedefinitions/postgresqls.composition.krateo.io\",\"verb\":\"DELETE\"}],\"color\":\"GREEN\",\"content\":\"${ .httpbin.json.location.city }\",\"icon\":\"ApartmentOutlined\",\"tags\":\"${ .httpbin.json.hobbies | join(\\\",\\\") }\",\"title\":\"${ (.httpbin.json.firstName  + \\\" \\\" + .httpbin.json.lastName) }\"}}}\n"
    },
    "managedFields": [
      {
        "manager": "kubectl-client-side-apply",
        "operation": "Update",
        "apiVersion": "widgets.ui.krateo.io/v1alpha1",
        "time": "2023-12-01T15:37:18Z",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:metadata": {
            "f:annotations": {
              ".": {},
              "f:kubectl.kubernetes.io/last-applied-configuration": {}
            }
          },
          "f:spec": {
            ".": {},
            "f:api": {},
            "f:app": {
              ".": {},
              "f:actions": {},
              "f:color": {},
              "f:content": {},
              "f:icon": {},
              "f:tags": {},
              "f:title": {}
            }
          }
        }
      }
    ]
  },
  "spec": {
    "app": {
      "title": "Charles Doe",
      "content": "San Fracisco",
      "icon": "ApartmentOutlined",
      "color": "GREEN",
      "tags": "chess,netflix",
      "actions": [
        {
          "name": "crd",
          "path": "/apis/apiextensions.k8s.io/v1/customresourcedefinitions/postgresqls.composition.krateo.io",
          "verb": "DELETE",
          "endpointRef": {
            "name": "krateogateway-endpoint",
            "namespace": "krateo-system"
          },
          "enabled": false
        }
      ]
    },
    "api": [
      {
        "name": "httpbin",
        "path": "/anything",
        "verb": "GET",
        "headers": [
          "Accept: application/json"
        ],
        "payload": "{\n  \"firstName\": \"Charles\",\n  \"lastName\": \"Doe\",\n  \"age\": 41,\n  \"location\": {\n    \"city\": \"San Fracisco\",\n    \"postalCode\": \"94103\"\n  },\n  \"hobbies\": [\n    \"chess\",\n    \"netflix\"\n  ]\n}\n",
        "endpointRef": {
          "name": "httpbin-endpoint",
          "namespace": "dev-system"
        },
        "enabled": true
      }
    ]
  }
}
```