
<h1 align="center">Ontology Layer2 Node</h1>
<h4 align="center">Version 1.0.0</h4>

English | [中文](README_CN.md)

Ontology Layer2致力于创建一个链下扩展方案，来满足用户低延迟、低费用的交易需求，其更好的应用扩展特性可以更好的支持大型复杂的应用。

Layer2 Node是ontology的Layer2交易收集器，它负责收集用户的Layer2交易，验证交易并周期性的提交Layer2的状态到ontology主链。可以理解为Node打包一批交易然后一次提交到主链来执行交易。

## 安装Node

### 前提

* Golang版本在1.14及以上
* 正确的Go语言开发环境
* Linux操作系统

### 获取Layer2 Node

克隆Layer2仓库到 **$GOPATH/src/github.com/ontio** 目录

```
$ git clone https://github.com/Layer2/node.git
```

### 编译
用make编译源码

```shell
$ make all
```

成功编译后会生成两个可以执行程序

* `node`: 节点程序/以命令行方式提供的节点控制程序

### 运行Layer Node

直接启动Ontology

   ```
	./node
   ```

## 许可证

Ontology遵守GNU Lesser General Public License, 版本3.0。 详细信息请查看项目根目录下的LICENSE文件。
