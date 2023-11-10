# Krateo BFF

Krateo **B**ackend **F**or **F**rontend.

## CardTemplate API

### Get all card templates

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Verb**          | `GET`                                                                               |
| **Path**          | `/apis/widgets.ui.krateo.io/v1alpha1/namespaces/${NAMESPACE}/cardtemplates`         |
| **Path Params**   | `${NAMESPACE}` is the namespace where your `CardTemplate` object is located         |
| **Query Params**  | `eval`: if _false_ do not evaluate `CardTemplate` expressions (default: true)       |

> Example:
>
> ```sh
> curl http://localhost:8080/apis/widgets.ui.krateo.io/v1alpha1/namespaces/dev-system/cardtemplates
> ```

### Get one card template

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Method**        | `GET`                                                                               |
| **Path**          | `/apis/widgets.ui.krateo.io/v1alpha1/namespaces/${NAMESPACE}/cardtemplates/${NAME}` |
| **Path Params**   | `${NAMESPACE}` is the namespace where your `CardTemplate` object is located         |
|                   | `${NAME}` is your `CardTemplate` object name                                        |
| **Query Params**  | `eval`: if _false_ do not evaluate `CardTemplate` expressions (default: true)       |

> Example:
>
> ```sh
> curl http://localhost:8080/apis/widgets.ui.krateo.io/v1alpha1/namespaces/dev-system/cardtemplates/card-dev
> ```
