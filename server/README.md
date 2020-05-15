# Layer2 Server

English|[中文](README_CN.md)

The Layer2 server consolidates Layer2 transaction information and provider query service. This service facilitiates fetching transaction address related details and other transaction records such as deposit and withdrawal history.

## Configuration

```json
{
  "log_level": 2,
  "rest_port": 30334,
  "version": "1.0.0",
  "http_max_connections":10000,
  "explorerdb_url":"127.0.0.1:3306",
  "explorerdb_user":"root",
  "explorerdb_password":"root1234",
  "explorerdb_name":"layer2"
}
```

| Field                | Description                                       |
| -------------------- | ------------------------------------------------- |
| log_level            | System log level                                  |
| rest_port            | Rest API port                                     |
| http_max_connections | Max. number of connections allowed simultaneously |
| explorerdb_url       | Database URL                                      |
| explorerdb_user      | Database login username                           |
| explorerdb_password  | Database login password                           |
| explorerdb_name      | Database name                                     |

The database that the Layer2 server accesses is the same database as that of the Layer2 operator. The database is configured to be used to by the Operator.

## Compilation

Run the following command in the main directory where the `main.go` file is located.

```go
go build
```

## Server API

The server API can be used to fetch account and transaction related details. The following APIs have been made available.

### 1. Fetch transaction history for an account

To fetch the transaction history of `AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc`,

Method: GET

URL: `http://{{host}}/api/v1/getlayer2tx/AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc`

### 2. Fetch deposit hitory for an account

To fetch the deposit hitory of `AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc`,

Method: GET

URL: `http://{{host}}/api/v1/getlayer2deposit/AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc`

### 3. Fetch withdrawal history for an account

To fetch the withdrawal history of `AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc`,

Method: GET

URL: `http://{{host}}/api/v1/getlayer2withdraw/AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc`

