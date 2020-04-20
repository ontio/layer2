
# cross chain explorer

## API

### 1 getexplorerinfo
查询信息

GET
```
http://{{host}}/api/v1/getexplorerinfo
```

#### example:

```
curl -i http://172.168.3.26:30334/api/v1/getexplorerinfo
```

```
{
    "action": "getexplorerinfo",
    "code": 1,
    "desc": "success",
    "result": "{\"chains\":[{\"id\":1,\"name\":\"multichain\",\"height\":1,\"in\":10000,\"out\":10000,\"contracts\":[{\"id\":0,\"name\":\"unkown\",\"contract\":\"xxx\"}]},{\"id\":2,\"name\":\"ontology\",\"height\":30000,\"in\":0,\"out\":10000,\"contracts\":[{\"id\":3,\"name\":\"neo\",\"contract\":\"xxx\"}]},{\"id\":3,\"name\":\"neo\",\"height\":30000,\"in\":10000,\"out\":0,\"contracts\":[{\"id\":2,\"name\":\"ontology\",\"contract\":\"xxx\"}]}],\"crosstxnumber\":10000}",
    "version": "1.0.0"
}
```

### 2 getcrosstxlist
查询跨链交易列表

POST
```
http://{{host}}/api/v1/getcrosstxlist
```

#### 参数:
```
{
    "start":"0",
    "end":"5"
}
```

#### example:

```
curl -i http://172.168.3.26:30334/api/v1/getcrosstxlist -X POST -d "{"start":"0","end":"5"}"
```

```
{
    "action": "getcrosstxlist",
    "code": 1,
    "desc": "success",
    "result": "[\"00004b359ea797697e9586896d8a9e162fd9a7cb93afbf1851a809a08a4bbdc0\",\"0000dc58a03fa300e59cf8a1bc022de513f023619ee9c8e256c7d3810f4c46cc\",\"0007f34ed7595eaeb599a30207e45b5bf0f3dbeb098eecd5b091f37c1c010320\",\"001280ce9707d5b8ad6bea7a843049b291c7bcd27e6e0679c89dbf84af2bd07c\",\"0018b75ce13823f294b92ca9b4f62a6571910b81398805f929b7185f78734094\",\"001e5dd2ff07883f279d8279892d5d6edf1119a2165646d2742dea436e1c784f\"]",
    "version": "1.0.0"
}
```

### 3 getcrosstx
查询跨链交易详细信息

GET
```
http://{{host}}/api/v1/getcrosstx/:txhash
```

#### example:

```
curl -i http://172.168.3.26:30334/api/v1/getcrosstx/00020fb0b090681648b50734dc835b26b0552adfa3186259089ed3e1ac0e7af1
```

```
{
    "action": "getcrosstx",
    "code": 1,
    "desc": "success",
    "result": "{\"fchaintr\":{\"chainid\":2,\"chainname\":\"ontology\",\"txhash\":\"980107a8ca1c2db41497391cc3487c0a4898de442036d24ecdb36553bef74ba4\",\"state\":1,\"tt\":1566463480,\"fee\":10000,\"height\":31608,\"tchainid\":3,\"tchainname\":\"neo\",\"key\":\"\",\"contract\":\"0200000000000000000000000000000000000000\",\"value\":\"AJLHY6wbJqE2j1VoRNPp3ZdacpLE58qNNy 123\",\"type\":0,\"typename\":\"unkown\",\"transfer\":{\"from\":\"AJLHY6wbJqE2j1VoRNPp3ZdacpLE58qNNy\",\"to\":\"AJLHY6wbJqE2j1VoRNPp3ZdacpLE58qNNy\",\"token\":\"0200000000000000000000000000000000000000\",\"amount\":123}},\"fchaintr_valid\":true,\"mchaintx\":{\"chainid\":1,\"chainname\":\"multichain\",\"txhash\":\"6c19f93157bc5f8ed7850b669f66c4457b3256fd1849b2be0a9f8da9aba86101\",\"state\":1,\"tt\":1566463472,\"fee\":10000,\"height\":11908,\"fchainid\":2,\"fchainname\":\"ontology\",\"ftxhash\":\"980107a8ca1c2db41497391cc3487c0a4898de442036d24ecdb36553bef74ba4\",\"tchainid\":3,\"tchainname\":\"neo\",\"key\":\"xx\"},\"mchaintx_valid\":true,\"tchaintx\":{\"chainid\":3,\"chainname\":\"neo\",\"txhash\":\"00020fb0b090681648b50734dc835b26b0552adfa3186259089ed3e1ac0e7af1\",\"state\":1,\"tt\":1566463488,\"fee\":10000,\"height\":311608,\"fchainid\":2,\"fchainname\":\"ontology\",\"mtxhash\":\"6c19f93157bc5f8ed7850b669f66c4457b3256fd1849b2be0a9f8da9aba86101\"},\"tchaintx_valid\":true}",
    "version": "1.0.0"
}
```


