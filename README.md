# Krateo BFF

Krateo **B**ackend **F**or **F**rontend.

## API

### CardTemplates

Get a single card template.

|         |       |
|:--------|:------------------------------------------------------------------------------------|
|*Method* | `GET`                                                                               |
|*Path*   | `/apis/widgets.ui.krateo.io/v1alpha1/namespaces/${NAMESPACE}/cardtemplates/${NAME}` |

where:

- `${NAMESPACE}` is the namespace where your `CardTemplate` object is located
- `${NAME}` is your `CardTemplate` object name
