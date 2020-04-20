# layer2合约接口
# 主网合约哈希为 xxxxxxxxxxx
## init(operator, stateRoot, confirmHeight)
该接口由operator节点调用，用于初始化合约
## 参数
operator为layer2发起转账的地址，为Address类型，stateRoot为合约初始化的状态根，其类型为array，结构为[stateRootHash, height, version]，分别代表状态根，当前layer2节点高度，layer2协议的版本号，类型分别为bytearray，integer，string。
## 返回值
如果调用成功返回True，否则返回False。
## deposit(player, amount, assetAddress)
该方法由用户调用，当用户想使用layer2的功能时调用该接口，用于将用户资产锁入合约。
## 参数
player, amount, assetAddress 分别为玩家地址，参与的金额，资产地址
## 返回值
如果调用成功返回True，否则返回False
## Notify
```
DepositEvent('deposit', currentId, player, amount, height, state, assetAddress)
```
currentId, player, amount, height, state, assetAddress分别代表depositId，玩家地址，存入金额，高度，状态，资产地址。
## withdraw(withdrawId)
该方法由用户调用
## 参数
withdrawId用于用户从合约中赎回withdrawId对应的金额
## 返回值
如果调用成功返回True，否则返回False
## Notify
```
WithdrawEvent('withdraw', 'withdrawId', 'amount', 'toAddress', 'height', 'status', 'assetAddress')

```
'withdrawId', 'amount', 'toAddress', 'height', 'status', 'assetAddress'分别代表withdrawId，赎回金额，出金地址，当前layer2节点高度，状态，资产地址。
## updateState(stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses)
该方法由operator地址调用，用于更新节点状态信息。
## 参数
stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses分别代表节点状态根，高度，版本号，需要更新状态depositId，withdraw金额集合，转移金额地址集合，转移资产地址集合。
## 返回值
调用成功返回True，否则返回False
## Notify
```
Notify(['updateState', stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses])
Notify(['updateDepositState', depositId])
WithdrawEvent(id, withdrawAmount, toAddresse, height, status, assetAddress)
```