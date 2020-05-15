# Layer2 Contract

English|[中文](README_CN.md)

## Layer2 Contract API

Say the main net contract hash is `xxxxxxxxxxx`.

## Method List

|                                                  Method Name                                                  | Description                             |
| :-----------------------------------------------------------------------------------------------------------: | --------------------------------------- |
|                                 [init](#initoperator-stateroot-confirmheight)                                 | Initializes the Layer2 contract         |
|                                 [deposit](#depositplayer-amount-assetaddress)                                 | Locks the user's assets in the contract |
| [updateState](#updatestatestateroothash-height-version-depositids-withdrawamounts-toaddresses-assetaddresses) | Updates the layer2 node's current state        |

## init(operator, stateRoot, confirmHeight)

This interface method is invoked by the operator to initialize the Layer2 contract.

|   Parameter   |  Type   | Description                                            |
| :-----------: | :-----: | ------------------------------------------------------ |
|   operator    | address | Address of the account that commit layer2 state to layer2 contract deployed on Ontology main net |
|   stateRoot   |  array  | The state root of layer2 node used to initialize the contract             |
| confirmHeight | integer | Height to confirm state finalty                        |

The `stateRoot` array contains the following fields:

|     Field     |   Type    | Description                    |
| :-----------: | :-------: | ------------------------------ |
| stateRootHash | bytearray | State root hash                |
|    height     |  integer  | Current Layer2 node height     |
|    version    |  string   | Layer2 protocol version number |

The interface method returns a `True` upon successful invocation, else `False`.

## deposit(player, amount, assetAddress)

This method is invoked by the user. This method is invoked when a user wants to use a Layer2 feature, and it locks the user's assets into the contract.

**Method Parameters**

|  Parameter   | Type      | Description    |
| :----------: | --------- | -------------- |
|    player    | address   | Payer address  |
|    amount    | integer   | Deposit amount |
| assetAddress | bytearray | Asset address  |

The method returns `True` upon successful invocation, else returns `False`.

The notification event for the deposit event is as follows:

```py
DepositEvent('deposit', currentId, player, amount, height, state, assetAddress)
```

The data obtained from the notification event is:

|    Field     | Description                                |
| :----------: | ------------------------------------------ |
|  currentId   | Deposit ID                                 |
|    player    | Player address                             |
|    amount    | Deposit amount                             |
|    height    | Block height for respective deposit action |
|    state     | The state of this transaction              |
| assetAddress | Asset contract address                     |


## updateState(stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses)

This method is invoked using the operator address and is used to update the node state information.

**Method Parameters**


|    Parameter    | Decsription                                        |
| :-------------: | -------------------------------------------------- |
|  stateRootHash  | Node state root hash                               |
|     height      | Block height                                       |
|     version     | Layer2 protocol version                            |
|   depositIds    | Deposit IDs whose state has updated on layer2 node |
| withdrawAmounts | Amount has withdrawn on layer2 node |
|   toAddresses   | Destination account addresses has withdrawn on layer2 node                      |
| assetAddresses  | Asset addresses has withdrawn on layer2 node                             |

The method returns `True` upon successful invocation, else returns `False`.


The notification event for the respective events are as follows: 

```py
Notify(['updateState', stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses])
Notify(['updateDepositState', depositId])
WithdrawEvent(id, withdrawAmount, toAddresse, height, status, assetAddress)
```

## Setting up Layer2 Contract

The process involves two major steps:

1. Deploy the contract on the Ontology main net
2. Initialize the deployed contract

### Deploying a Contract on Ontology 

Ontology's online smart contract development IDE SmartX can be used to write, compile, deploy the Layer2 smart contract on Ontology main net. Once the Layer2 contract is deployed successfully, we will obtain the contract address.

**SmartX web address:** https://smartx.ont.io/#/

> Please note that SmartX is currently supported on Chrome browser only.

**Layer2 contract template:** https://github.com/ontio/layer2/blob/master/contract/layer2.py


### Initializing the Contract

After successful deployment, the contract can be initialized using SmartX by executing the `init()` method.

![](pic/init_smart.jpg)

The `init()` method takes three parameters.

|   Parameter   |  Type   | Description                                                                                                                                                                                            |
| :-----------: | :-----: | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
|   Operator    | Address | Layer2 security daemon account address, sends latest Layer2 state to the contract as proof to facilitate on-chain arbitration in case of off-chain malice by Layer2, Ontology main net account address |
|   stateRoot   |  Array  | Defines Layer2 genesis state                                                                                                                                                                           |
| confirmHeight | Integer | Height to determine state finalty and validity of withdraw action after the transaction has taken place                                                                                                |

Sample `stateRoot`:

```json
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
| Field Type | Value Description                                                                                                                               |
| :--------: | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| Bytearray  | Layer2 account state root hash, `0000000000000000000000000000000000000000000000000000000000000000` indicates account state was empty at genesis |
|  Integer   | Block height at the respective state, initialized with `0`                                                                                      |
|   String   | Version no., current version number is `1.0.0`                                                                                                  |
