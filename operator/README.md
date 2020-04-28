# layer2_operator

Layer2 Operator负责监听ontology主链是否有到Layer2的代币转移或者Layer2都ontology主链的代币转移交易，同时Operator还负责周期性的将Layer2的state提交到ontology主网作为证明。

## 安装Operator

### 安装Mysql

选择适合自己的平台环境来安装mysql，指导手册https://www.mysql.com/cn/products/community/

### 下载源码

```
git clone https://ontio/Layer2/layer2-operator.git
```

### 编译

```
go build main.go
```

### 配置

在源码目录下有config.json配置文件，是operator启动的配置文件。
```
{
  "OntologyConfig":{
    "RestURL":"http://polaris1.ont.io:20336",
    "Layer2ContractAddress":"4229a92d90d446d1598e12e35698b681ae4d4642",
    "WalletFile":"./wallet_ontology.dat",
    "WalletPwd":"1",
    "GasPrice":0,
    "GasLimit":2000000
  },
  "Layer2Config":{
    "RestURL":"http://localhost:40336",
    "WalletFile":"./wallet_layer2.dat",
    "WalletPwd":"1",
    "GasPrice":0,
    "GasLimit":2000000
  },
  "DBConfig":{
    "ProjectDBUrl":"127.0.0.1:3306",
    "ProjectDBUser":"root",
    "ProjectDBPassword":"root1234",
    "ProjectDBName":"layer2"
  }
}
```
同时需要两个钱包文件，分别为ontology主网的钱包文件和layer2节点的钱包文件。