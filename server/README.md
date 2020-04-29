
# Layer2 Server

Layer2 Server收集Layer2的交易信息，提供查询服务，包括查询指定地址的历史交易记录，查询指定地址的deposit和withdraw记录。

## 配置说明

```
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
Layer2 Server需要访问的数据库和Layer2 Operator的数据库是同一个数据库，数据库配置为Operator使用的数据库。

## 编译

```
go build main.go
```

## API

### 1 查询历史交易记录
如查询AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc地址的历史交易记录

GET

```
http://{{host}}/api/v1/getlayer2tx/AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc
```

### 2 查询deposit记录
如查询AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc地址的deposit记录

GET
```
http://{{host}}/api/v1/getlayer2deposit/AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc
```

### 3 查询withdraw记录
如查询AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc地址的withdraw记录

GET
```
http://{{host}}/api/v1/getlayer2withdraw/AMUGPqbVJ3TG6pe7xRpxxaeh4ai4fu9ahc
```


