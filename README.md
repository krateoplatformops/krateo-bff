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
