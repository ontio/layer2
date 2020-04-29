# Layer2合约

## Layer2合约接口

主网合约哈希为 xxxxxxxxxxx

### init(operator, stateRoot, confirmHeight)
该接口由operator节点调用，用于初始化合约

#### 参数
operator为layer2发起转账的地址，为Address类型，stateRoot为合约初始化的状态根，其类型为array，结构为[stateRootHash, height, version]，分别代表状态根，当前layer2节点高度，layer2协议的版本号，类型分别为bytearray，integer，string。

#### 返回值
如果调用成功返回True，否则返回False。

### deposit(player, amount, assetAddress)
该方法由用户调用，当用户想使用layer2的功能时调用该接口，用于将用户资产锁入合约。

#### 参数
player, amount, assetAddress 分别为玩家地址，参与的金额，资产地址

#### 返回值
如果调用成功返回True，否则返回False

#### Notify
```
DepositEvent('deposit', currentId, player, amount, height, state, assetAddress)
```
currentId, player, amount, height, state, assetAddress分别代表depositId，玩家地址，存入金额，高度，状态，资产地址。

### withdraw(withdrawId)
该方法由用户调用

#### 参数
withdrawId用于用户从合约中赎回withdrawId对应的金额

#### 返回值
如果调用成功返回True，否则返回False

#### Notify
```
WithdrawEvent('withdraw', 'withdrawId', 'amount', 'toAddress', 'height', 'status', 'assetAddress')

```
'withdrawId', 'amount', 'toAddress', 'height', 'status', 'assetAddress'分别代表withdrawId，赎回金额，出金地址，当前layer2节点高度，状态，资产地址。

### updateState(stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses)
该方法由operator地址调用，用于更新节点状态信息。

#### 参数
stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses分别代表节点状态根，高度，版本号，需要更新状态depositId，withdraw金额集合，转移金额地址集合，转移资产地址集合。

#### 返回值
调用成功返回True，否则返回False

#### Notify
```
Notify(['updateState', stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses])
Notify(['updateDepositState', depositId])
WithdrawEvent(id, withdrawAmount, toAddresse, height, status, assetAddress)
```
## 安装Layer2合约

在ontology主链安装Layer2合约包括两步：

1 在ontology主链部署合约

2 初始化合约

### 部署合约到ontology

在ontology合约开发工具smartx上按照指定合约开发部署流程将Layer2合约部署到ontology主链，在成功部署Layer2合约后，我们得到Layer2合约地址。

smartx工具web地址： https://smartx.ont.io/#/

Layer2合约代码： https://github.com/ontio/layer2/blob/master/contract/layer2.py

### 初始化合约

继续在ontology合约开发工具smartx上初始化合约，在smartx上成功部署合约到ontology后，执行该合约的init方法来初始化合约。

![](pic/init_smart.jpg)

该合约接口有三个参数：
Address数据类型的operator、Array类型的stateRoot、Integer类型的confirmHeight。

operator指定了Laye2安全守护账户，安全守护账户会将Layer2最新状态提交到该Layer2合约作为证明，当在链下的Layer2有作恶或者纠纷时，可以依据此证明在链上裁决。ontology的安全守护账户就是第一步生成的ontology账户。

stateRoot指定了Layer2的创世状态，如以下创世状态
```
[
    {
        "type": "ByteArray",
        "value": "0000000000000000000000000000000000000000000000000000000000000000"
    },
    {
        "type": "Integer",
        "value": 0
    },
    {
        "type": "String",
        "value": "1.0.0"
    }
]
```
ByteArray的"0000000000000000000000000000000000000000000000000000000000000000"是Layer2账户状态树根hash，表明了Layer2在创世状态下账户状态为空。

Integer的0是状态高度，初始为0。

String的"1.0.0"是版本，当前版本为1.0.0。

confirmHeight指定了用户在Layer2进行withdraw后，需要多少个状态高度来确认withdraw有效。