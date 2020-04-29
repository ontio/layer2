# Layer2 Operator

Layer2 Operator是Layer2的安全守护程序，负责监听ontology主链是否有到Layer2的代币转移或者Layer2都ontology主链的代币转移交易，同时Operator还负责周期性的将Layer2的state提交到ontology主网作为证明。

## 安装Operator

### 安装Mysql

选择适合自己的平台环境来安装mysql，指导手册https://www.mysql.com/cn/products/community/

Mysql安装完毕后，初始化数据库，在Mysql上创建layer2数据库：
```
CREATE SCHEMA IF NOT EXISTS `layer2` DEFAULT CHARACTER SET utf8;
USE `layer2`;

DROP TABLE IF EXISTS `chain_info`;
CREATE TABLE `chain_info` (
 `name` VARCHAR(100) NOT NULL COMMENT '链名称',
 `id`  INT(4) NOT NULL COMMENT '链id',
 `url` VARCHAR(256) NOT NULL COMMENT '访问链的url',
 `height` INT(4) NOT NULL COMMENT '解析的区块高度',
 PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;

INSERT INTO `chain_info`(`name`,`id`,`url`,`height`) VALUES("ontology",1,"http://138.91.6.125:20336",0);
INSERT INTO `chain_info`(`name`,`id`,`url`,`height`) VALUES("layer2",2,"http://47.90.189.186:40332",0);

DROP TABLE IF EXISTS `deposit`;
CREATE TABLE `deposit` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT '交易hash',
 `tt` INT(4) NOT NULL COMMENT '交易时间',
 `state` INT(1) NOT NULL COMMENT '交易状态',
 `height` INT(4) NOT NULL COMMENT '交易的高度',
 `fromaddress` VARCHAR(256) NOT NULL COMMENT '地址',
 `amount` BIGINT(8) NOT NULL COMMENT 'deposit的金额',
 `tokenaddress` VARCHAR(256) NOT NULL COMMENT '币地址',
 `id` INT(4) NOT NULL COMMENT '交易的高度',
 `layer2txhash` VARCHAR(256) DEFAULT NULL COMMENT 'layer2交易hash',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `withdraw`;
CREATE TABLE `withdraw` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT '交易hash',
 `tt` INT(4) NOT NULL COMMENT '交易时间',
 `state` INT(1) NOT NULL COMMENT '交易状态',
 `height` INT(4) NOT NULL COMMENT '交易的高度',
 `toaddress` VARCHAR(256) NOT NULL COMMENT '地址',
 `amount` BIGINT(8) NOT NULL COMMENT 'deposit的金额',
 `tokenaddress` VARCHAR(256) NOT NULL COMMENT '币地址',
 `ontologytxhash` VARCHAR(256) DEFAULT NULL COMMENT '交易hash',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `layer2tx`;
CREATE TABLE `layer2tx` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT '交易hash',
 `state` INT(1) NOT NULL COMMENT '交易状态',
 `tt` INT(4) NOT NULL COMMENT '交易时间',
 `fee` BIGINT(8) NOT NULL COMMENT '交易手续费',
 `height` INT(4) NOT NULL COMMENT '交易的高度',
 `fromaddress` VARCHAR(256) NOT NULL COMMENT '地址',
 `tokenaddress` VARCHAR(256) NOT NULL COMMENT '执行的合约',
 `toaddress` VARCHAR(256) NOT NULL COMMENT '地址',
 `amount` BIGINT(8) NOT NULL COMMENT 'deposit的金额',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `layer2commit`;
CREATE TABLE `layer2commit` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT '交易hash',
 `state` INT(1)  DEFAULT 0 COMMENT '交易状态',
 `tt` INT(4) DEFAULT 0 COMMENT '交易时间',
 `fee` BIGINT(8) DEFAULT 0 COMMENT '交易手续费',
 `ontologyheight` INT(4) DEFAULT 0 COMMENT '交易的高度',
 `layer2height` INT(4) DEFAULT 0 COMMENT '交易的高度',
 `layer2msg` VARCHAR(1024) NOT NULL COMMENT 'laeyr2 msg',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
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
主要包括：

ontology的访问配置：节点地址、以上第二步部署的layer2合约地址，以上第一步生成的ontology钱包文件wallet_ontology.dat及其密码。

Node的访问配置：节点地址、以上第一步生成的Layer2钱包文件wallet_layer2.dat及其密码。

Mysql数据库访问配置：数据库URL、用户名和密码以及Layer2数据库名称。