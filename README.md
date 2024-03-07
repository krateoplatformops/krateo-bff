# Krateo BFF

Krateo **B**ackend **F**or **F**rontend.

## CardTemplate API

### List all card templates

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Verb**          | `GET`                                                                               |
| **Path**          | `/apis/widgets.ui.krateo.io/cardtemplates`                                          |
| **Query Params**  | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates (optional)                       |
|                   | `version`: cardtemplates schema version (optional)                                  |

**Example**

```sh
curl "http://localhost:8090/apis/widgets.ui.krateo.io/cardtemplates?sub=cyberjoker&orgs=devs&namespace=demo-system"
```

### Get one card template

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Method**        | `GET`                                                                               |
| **Path**          | `/apis/widgets.ui.krateo.io/cardtemplates/${NAME}`                                  |
| **Path Params**   | `${NAME}` is your `CardTemplate` object name                                        |
| **Query Params**  | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates                                  |
|                   | `version`: cardtemplates schema version (optional)                                  |

**Example**:

```sh
curl "http://localhost:8090/apis/widgets.ui.krateo.io/cardtemplates/one?sub=cyberjoker&orgs=devs&namespace=demo-system"
```

## Column API

### List all columns

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Verb**          | `GET`                                                                               |
| **Path**          | `/apis/layout.ui.krateo.io/columns`                                          |
| **Query Params**  | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates (optional)                       |
|                   | `version`: cardtemplates schema version (optional)                                  |

**Example**:

```sh
curl "http://localhost:8090/apis/layout.ui.krateo.io/columns?sub=cyberjoker&orgs=devs&namespace=demo-system"
```

### Get one column

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Method**        | `GET`                                                                               |
| **Path**          | `/apis/layout.ui.krateo.io/columns/${NAME}`                                  |
| **Path Params**   | `${NAME}` is your `Column` object name                                        |
| **Query Params**  | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates                                  |
|                   | `version`: cardtemplates schema version (optional)                                  |

**Example**:

```sh
curl "http://localhost:8090/apis/layout.ui.krateo.io/columns/three?sub=cyberjoker&orgs=devs&namespace=demo-system"
```

## Rows API

### List all rows

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Verb**          | `GET`                                                                               |
| **Path**          | `/apis/layout.ui.krateo.io/rows`                                          |
| **Query Params**  | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates (optional)                       |
|                   | `version`: cardtemplates schema version (optional)                                  |

**Example**:

```sh
curl "http://localhost:8090/apis/layout.ui.krateo.io/rows?sub=cyberjoker&orgs=devs&namespace=demo-system"
```

### Get one row

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Method**        | `GET`                                                                               |
| **Path**          | `/apis/layout.ui.krateo.io/rows/${NAME}`                                  |
| **Path Params**   | `${NAME}` is your `Row` object name                                        |
| **Query Params**  | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates                                  |
|                   | `version`: cardtemplates schema version (optional)                                  |

**Example**:

```sh
curl "http://localhost:8090/apis/layout.ui.krateo.io/rows/two?sub=cyberjoker&orgs=devs&namespace=demo-system"
```

## FormTemplate API


### List all rows

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Verb**          | `GET`                                                                               |
| **Path**          | `/apis/widgets.ui.krateo.io/formtemplates`                                          |
| **Query Params**  | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates (optional)                       |
|                   | `version`: cardtemplates schema version (optional)                                  |

**Example**:

```sh
curl "http://localhost:8090/apis/widgets.ui.krateo.io/formtemplates?sub=cyberjoker&orgs=devs&namespace=demo-system"
```

### Get one formtemplate

|                   |                                                                                     |
|:------------------|:------------------------------------------------------------------------------------|
| **Method**        | `GET`                                                                               |
| **Path**          | `/apis/widgets.ui.krateo.io/formtemplates/${NAME}`                                  |
| **Path Params**   | `${NAME}` is your `Row` object name                                        |
| **Query Params**  | `sub`: username (subject)                                                           |
|                   | `orgs`: comma separated organizations                                               |
|                   | `namespace`: namespace where to list cardtemplates                                  |
|                   | `version`: cardtemplates schema version (optional)                                  |

**Example**:

```sh
curl "http://localhost:8090/apis/widgets.ui.krateo.io/formtemplates/fireworksapp?sub=cyberjoker&orgs=devs&namespace=demo-system"
```

# Apps API

## Create or Update an App

|                   |                                                                   |
|:------------------|:------------------------------------------------------------------|
| **Method**        | `POST`                                                            |
| **Path**          | `/apis/apps/${NAME}`                                              |
| **Path Params**   | `${NAME}` is your _"app"_ object name                             |
| **Query Params**  | `sub`: username (subject)                                         |
|                   | `orgs`: comma separated organizations                             |
|                   | `namespace`: namespace where your _"app"_ resource belongs        |
|                   | `group`: api group of your custom _"app"_ resource                |
|                   | `kind`: kind of your _"app"_ resource                             |
|                   | `version`: api schema version of your _"app"_ resource (optional) |

**Example**:

```sh
curl -X POST -H "Content-Type: application/json" \
    --data @fireworksapp.json \
    "http://localhost:8090/apis/apps/fireworksapp?sub=cyberjoker&orgs=devs&kind=Fireworksapp&group=apps.krateo.io&namespace=demo-system"
```
