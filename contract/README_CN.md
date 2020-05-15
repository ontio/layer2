# Layer2合约

中文|[English](README.md)

## Layer2合约接口

主网合约哈希为 xxxxxxxxxxx

## 接口列表

|                                                  接口名称                                                 | 描述                             |
| :-----------------------------------------------------------------------------------------------------------: | --------------------------------------- |
|                                 [init](#initoperator-stateroot-confirmheight)                                 | 初始化layer2合约         |
|                                 [deposit](#depositplayer-amount-assetaddress)                                 | 锁定用户资产到合约，用于在layer2释放资产给用户 |
| [updateState](#updatestatestateroothash-height-version-depositids-withdrawamounts-toaddresses-assetaddresses) | 更新layer2的最新状态信息|

## init(operator, stateRoot, confirmHeight)
该接口由operator节点调用，用于初始化合约

|   Parameter   |  Type   | Description                                            |
| :-----------: | :-----: | ------------------------------------------------------ |
|   operator    | address | 提交layer2的状态都部署在Ontology主网的layer2合约的账户 |
|   stateRoot   |  array  | layer2初始状态根信息            |
| confirmHeight | integer | 状态确认的区块数                      |

stateRoot包含以下内容:

|     Field     |   Type    | Description                    |
| :-----------: | :-------: | ------------------------------ |
| stateRootHash | bytearray | 状态根hash                |
|    height     |  integer  | 当前状态的区块高度     |
|    version    |  string   | Layer2 协议版本号 |

如果调用成功返回True，否则返回False。

## deposit(player, amount, assetAddress)
该方法由用户调用，当用户想使用layer2的功能时调用该接口，用于将用户资产锁入合约。

|  Parameter   | Type      | Description    |
| :----------: | --------- | -------------- |
|    player    | address   | 使用layer2的用户  |
|    amount    | integer   | 入金到layer2的金额|
| assetAddress | bytearray | 资产地址  |


如果调用成功返回True，否则返回False

### Notify
```
DepositEvent('deposit', currentId, player, amount, height, state, assetAddress)
```
|    Field     | Description                                |
| :----------: | ------------------------------------------ |
|  currentId   | 入金操作生成的唯一ID                               |
|    player    | 使用layer2的用户                            |
|    amount    | 入金到layer2的金额                             |
|    height    | 当前交易的高度 |
|    state     | 当前交易的状态              |
| assetAddress | 资产地址                     |


## updateState(stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses)
该方法由operator地址调用，用于更新节点状态信息。

|    Parameter    | Decsription                                        |
| :-------------: | -------------------------------------------------- |
|  stateRootHash  | Layer2新区块的状态信息                               |
|     height      | Layer2新区块的高度                                     |
|     version     | Layer2协议版本                            |
|   depositIds    | 在Layer2已经入金到账户的deposit |
| withdrawAmounts | 在layer2已经提现的金额 |
|   toAddresses   | 在layer2已经提现的账户                      |
| assetAddresses  | 在layer2已经提现的资产   

调用成功返回True，否则返回False

### Notify
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