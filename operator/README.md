# Layer2 Operator

English|[中文](README_CN.md)

Layer2 operator is a security daemon that runs on Layer2. It monitors the Layer2 and Ontology main chain token transfer transactions and periodically sends the Layer2 state to the Ontology main chain as proof.

## Operator Installation

### Installing MySQL

Install MySQL on a suitable platform. For details on how to download and install MySQL please refer to https://www.mysql.com/products/community/

After successfully installing and initializing the database system, create the Layer2 database in the following manner:

```sql
CREATE SCHEMA IF NOT EXISTS `layer2` DEFAULT CHARACTER SET utf8;
USE `layer2`;

DROP TABLE IF EXISTS `chain_info`;
CREATE TABLE `chain_info` (
 `name` VARCHAR(100) NOT NULL COMMENT 'Chain Name',
 `id`  INT(4) NOT NULL COMMENT 'Chain ID',
 `url` VARCHAR(256) NOT NULL COMMENT 'URL of accessing chain',
 `height` INT(4) NOT NULL COMMENT 'Block height',
 PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;

INSERT INTO `chain_info`(`name`,`id`,`url`,`height`) VALUES("ontology",1,"http://138.91.6.125:20336",0);
INSERT INTO `chain_info`(`name`,`id`,`url`,`height`) VALUES("layer2",2,"http://47.90.189.186:40332",0);

DROP TABLE IF EXISTS `deposit`;
CREATE TABLE `deposit` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT 'Transaction hash',
 `tt` INT(4) NOT NULL COMMENT 'Transaction time',
 `state` INT(1) NOT NULL COMMENT 'Transaction state',
 `height` INT(4) NOT NULL COMMENT 'Transaction block height',
 `fromaddress` VARCHAR(256) NOT NULL COMMENT 'Source address',
 `amount` BIGINT(8) NOT NULL COMMENT 'Deposit amount',
 `tokenaddress` VARCHAR(256) NOT NULL COMMENT 'Token contract address',
 `id` INT(4) NOT NULL COMMENT 'ID',
 `layer2txhash` VARCHAR(256) DEFAULT NULL COMMENT 'Layer2 transaction hash',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `withdraw`;
CREATE TABLE `withdraw` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT 'Transaction hash',
 `tt` INT(4) NOT NULL COMMENT 'Transaction time',
 `state` INT(1) NOT NULL COMMENT 'Transaction state',
 `height` INT(4) NOT NULL COMMENT 'Transaction block height',
 `toaddress` VARCHAR(256) NOT NULL COMMENT 'Destination address',
 `amount` BIGINT(8) NOT NULL COMMENT 'Deposit amount',
 `tokenaddress` VARCHAR(256) NOT NULL COMMENT 'Token contract address',
 `ontologytxhash` VARCHAR(256) DEFAULT NULL COMMENT 'Transaction hash',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `layer2tx`;
CREATE TABLE `layer2tx` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT 'Transaction hash',
 `state` INT(1) NOT NULL COMMENT 'Transaction state',
 `tt` INT(4) NOT NULL COMMENT 'Transaction time',
 `fee` BIGINT(8) NOT NULL COMMENT 'Transaction fee',
 `height` INT(4) NOT NULL COMMENT 'Transaction block height',
 `fromaddress` VARCHAR(256) NOT NULL COMMENT 'Source address',
 `tokenaddress` VARCHAR(256) NOT NULL COMMENT 'Token contract address',
 `toaddress` VARCHAR(256) NOT NULL COMMENT 'Destination address',
 `amount` BIGINT(8) NOT NULL COMMENT 'Deposit amount',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `layer2commit`;
CREATE TABLE `layer2commit` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT 'Transaction hash',
 `state` INT(1)  DEFAULT 0 COMMENT 'Transaction state',
 `tt` INT(4) DEFAULT 0 COMMENT 'Transaction time',
 `fee` BIGINT(8) DEFAULT 0 COMMENT 'Transaction fee',
 `ontologyheight` INT(4) DEFAULT 0 COMMENT 'Transaction block height',
 `layer2height` INT(4) DEFAULT 0 COMMENT 'Transaction block height',
 `layer2msg` VARCHAR(1024) NOT NULL COMMENT 'layer2 msg',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
```

### Compilation

Run the following command in the directory with the `main.go` file.

```go
go build main.go
```

### Configuration

The `config.json` configuration file in the source code directory is used to start the operator.

```json
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

As illustrated by the above sample configuration, the `config.json` file contains access parameters to:

- **Ontology:** Node address, Layer2 contract address, Ontology `.dat` wallet file, and the wallet password.
- **Node:** Node address, Layer2 `.dat` wallet file, and the wallet password.
- **MySQL:** Database URL, username, password, and database name.