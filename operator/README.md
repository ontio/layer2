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
 `name` VARCHAR(100) NOT NULL COMMENT '链名称',
 `id`  INT(4) NOT NULL COMMENT '链id',
 `height` INT(4) NOT NULL COMMENT '解析的区块高度',
 PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;

INSERT INTO `chain_info`(`name`,`id`,`height`) VALUES("ontology",1,0);
INSERT INTO `chain_info`(`name`,`id`,`height`) VALUES("layer2",2,0);

DROP TABLE IF EXISTS `deposit`;
CREATE TABLE `deposit` (
 `txhash`  VARCHAR(256) NOT NULL COMMENT '交易hash',
 `tt` INT(4) NOT NULL COMMENT '交易时间',
 `state` INT(1) NOT NULL COMMENT '交易状态',
 `height` INT(4) NOT NULL COMMENT '交易的高度',
 `fromaddress` VARCHAR(256) NOT NULL COMMENT '地址',
 `amount` BIGINT(8) NOT NULL COMMENT 'deposit的金额',
 `tokenaddress` VARCHAR(256) NOT NULL COMMENT '币地址',
 `id` INT(4) NOT NULL COMMENT '交易的ID',
 `layer2txhash` VARCHAR(256) DEFAULT NULL COMMENT 'layer2交易hash',
 PRIMARY KEY (`id`),
 UNIQUE (`txhash`)
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
 `layer2height` INT(4) DEFAULT 0 COMMENT '交易的高度',
 `layer2msg` VARCHAR(1024) NOT NULL COMMENT 'laeyr2 msg',
 PRIMARY KEY (`txhash`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
```

### Compilation

Run the following command in the directory with the `main.go` file.

```go
cd $gopath/src/github.com/ontio
git clone https://github.com/ontio/layer2/
cd ./layer2/operator
go build
```