OntCversion = '2.0.0'
from ontology.builtins import state, concat
from ontology.interop.Ontology.Native import Invoke
from ontology.interop.System.Action import RegisterAction
from ontology.interop.System.App import DynamicAppCall
from ontology.interop.System.Blockchain import GetHeight
from ontology.interop.System.ExecutionEngine import GetExecutingScriptHash
from ontology.interop.System.Runtime import CheckWitness, Serialize, Deserialize, Notify
from ontology.interop.System.Storage import GetContext, Get, Put
from ontology.libont import bytearray_reverse

ONGAddress = bytearray(b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02')
ONTAddress = bytearray(b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01')
ContractAddress = GetExecutingScriptHash()

ADMIN = 'admin'

DepositEvent = RegisterAction('deposit', 'depositId', 'fromAddress', 'amount', 'height', 'status', 'assetAddress')

WithdrawEvent = RegisterAction('withdraw', 'withdrawId', 'amount', 'toAddress', 'height', 'status', 'assetAddress')

DEPOSIT_PREFIX = 'deposit'

WITHDRAW_PREFIX = 'withdraw'

CURRENT_DEPOSIT_ID = 'currentDepositId'

INITED = 'Initialized'

CURRENT_WITHDRAW_ID = 'currentWithdrawId'

CURRENT_HEIGHT = 'currentHeight'

Current_STATE_PREFIX = 'stateRoot'

CONFRIM_HEIGHT = 'confirmHeight'

OPERATOR_ADDRESS = 'operator'


def Main(operation, args):
    ## FOR OPERATOR INVOkE ONLY
    if operation == 'init':
        assert (len(args) == 3)
        operator = args[0]
        stateRoot = args[1]
        confirmHeight = args[2]
        return init(operator, stateRoot, confirmHeight)

    if operation == 'deposit':
        assert (len(args) == 3)
        player = args[0]
        amount = args[1]
        assetAddress = args[2]
        return deposit(player, amount, assetAddress)

    if operation == 'updateState':
        assert (len(args) == 7)
        stateRootHash = args[0]
        height = args[1]
        version = args[2]
        depositIds = args[3]
        withdrawAmounts = args[4]
        toAddresses = args[5]
        assetAddresses = args[6]
        return updateState(stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses)

    if operation == 'getStateRootByHeight':
        assert (len(args) == 1)
        height = args[0]
        return getStateRootByHeight(height)
    if operation == 'getCurrentHeight':
        assert (len(args) == 0)
        return getCurrentHeight()
    return True


## 初始化operator和status信息 [stateRootHash, height, version]
def init(operator, stateRoot, confirmHeight):
    inited = Get(GetContext(), INITED)
    if inited:
        Notify(["idiot admin, you have initialized the contract"])
        return False
    else:
        assert (len(operator) == 20)
        assert (CheckWitness(operator))
        assert (len(stateRoot) == 3)
        Put(GetContext(), INITED, 1)
        assert (confirmHeight > 0)
        Put(GetContext(), CONFRIM_HEIGHT, confirmHeight)
        Put(GetContext(), OPERATOR_ADDRESS, operator)
        Put(GetContext(), CURRENT_HEIGHT, stateRoot[1])
        stateRootInfo = Serialize(stateRoot)
        Put(GetContext(), concatKey(Current_STATE_PREFIX, stateRoot[1]), stateRootInfo)
        Notify(["Initialized contract successfully"])
    return True


## 用户为了使用Layer2发起交易将资金质押在合约中
def deposit(player, amount, assetAddress):
    assert (CheckWitness(player))

    currentId = Get(GetContext(), CURRENT_DEPOSIT_ID)
    if not currentId:
        currentId = 1
    height = Get(GetContext(), CURRENT_HEIGHT)

    if assetAddress == ONTAddress:
        assert (_transferONT(player, ContractAddress, amount))
    elif assetAddress == ONGAddress:
        assert (_transferONG(player, ContractAddress, amount))
    else:
        reverseAssetAddress = bytearray_reverse(assetAddress)
        assert (_transferOEP4(reverseAssetAddress, player, ContractAddress, amount))

    depositRecord = [player, amount, height, 0]
    depositRecordInfo = Serialize(depositRecord)

    Put(GetContext(), CURRENT_DEPOSIT_ID, currentId + 1)
    Put(GetContext(), concatKey(DEPOSIT_PREFIX, currentId), depositRecordInfo)
    DepositEvent(currentId, player, amount, height, 0, assetAddress)
    return True

## 用户将质押在合约中用于Layer2交易的资产赎回
def withdraw(withdrawId):
    withdrawStatusInfo = Get(GetContext(), concatKey(WITHDRAW_PREFIX, withdrawId))
    assert (withdrawStatusInfo)
    withdrawStatus = Deserialize(withdrawStatusInfo)
    currentHeight = GetHeight()
    confirmHeight = Get(GetContext(), CONFRIM_HEIGHT)
    assert (currentHeight - withdrawStatusInfo[3] >= confirmHeight)
    assetAddress = withdrawStatus[5]

    assert (withdrawStatus[4] == 0)
    if assetAddress == ONTAddress:
        assert (_transferONTFromContact(withdrawStatus[2], withdrawStatus[1]))
    elif assetAddress == ONGAddress:
        assert (_transferONGFromContact(withdrawStatus[2], withdrawStatus[1]))
    else:
        reverseAssetAddress = bytearray_reverse(assetAddress)
        assert (_transferOEP4FromContact(reverseAssetAddress, withdrawStatus[2], withdrawStatus[1]))

    WithdrawEvent(withdrawStatus[0], withdrawStatus[1], withdrawStatus[2], withdrawStatus[3], 1, withdrawStatus[5])
    return True


## 更新全局的状态根，合约需要验证签名的有效性
def updateState(stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses):
    operator = Get(GetContext(), OPERATOR_ADDRESS)
    assert (CheckWitness(operator))
    preHeight = Get(GetContext(), CURRENT_HEIGHT)
    assert (preHeight + 1 == height)

    Put(GetContext(), CURRENT_HEIGHT, height)
    stateRoot = [stateRootHash, height, version]
    stateRootInfo = Serialize(stateRoot)
    Put(GetContext(), concatKey(Current_STATE_PREFIX, height), stateRootInfo)
    # 返回满足条件用户的钱
    currentWithDrawId = Get(GetContext(), CURRENT_WITHDRAW_ID)
    confirmHeight = Get(GetContext(), CONFRIM_HEIGHT)
    if currentWithDrawId:
        while currentWithDrawId > 1:
            withdrawStatusInfo = Get(GetContext(), concatKey(WITHDRAW_PREFIX, currentWithDrawId - 1))
            withdrawStatus = Deserialize(withdrawStatusInfo)
            if (height - withdrawStatus[3] == confirmHeight):
                assert (withdraw(withdrawStatus[0]))
            elif height - withdrawStatus[3] > confirmHeight:
                break
            currentWithDrawId = currentWithDrawId - 1
    # 更新deposit状态
    _updateDepositState(depositIds)
    # 更新withdraw状态
    _createWithdrawState(height, withdrawAmounts, toAddresses, assetAddresses)
    Notify(['updateState', stateRootHash, height, version, depositIds, withdrawAmounts, toAddresses, assetAddresses])
    return True


## 根据高度获取状态根信息
def getStateRootByHeight(height):
    stateRootInfo = Get(GetContext(), concatKey(Current_STATE_PREFIX, height))
    if stateRootInfo:
        stateRoot = Deserialize(stateRootInfo)
        return stateRoot
    return []


def getCurrentHeight():
    height = Get(GetContext(), CURRENT_HEIGHT)
    return height


def _updateDepositState(depositIds):
    for i in range(len(depositIds)):
        depositStatusInfo = Get(GetContext(), concatKey(DEPOSIT_PREFIX, depositIds[i]))
        if depositStatusInfo:
            depositStatus = Deserialize(depositStatusInfo)
            assert (depositStatus[3] == 0)
            depositStatus[3] = 1
            depositStatusInfo = Serialize(depositStatus)
            Put(GetContext(), concatKey(DEPOSIT_PREFIX, depositIds[i]), depositStatusInfo)
            Notify(['updateDepositState', depositIds[i]])
    return True


def _createWithdrawState(height, withdrawAmounts, toAddresses, assetAddresses):
    assert (len(withdrawAmounts) == len(toAddresses))
    length = len(withdrawAmounts)
    id = Get(GetContext(), CURRENT_WITHDRAW_ID)
    if not id:
        id = 1
    for i in range(length):
        assert (withdrawAmounts[i] > 0)
        assert (len(toAddresses[i]) == 20)
        withdrawStatus = [id, withdrawAmounts[i], toAddresses[i], height, 0, assetAddresses[i]]
        withdrawStatusInfo = Serialize(withdrawStatus)
        Put(GetContext(), concatKey(WITHDRAW_PREFIX, id), withdrawStatusInfo)
        WithdrawEvent(id, withdrawAmounts[i], toAddresses[i], height, 0, assetAddresses[i])
        id = id + 1
    Put(GetContext(), CURRENT_WITHDRAW_ID, id)
    return True


### 内部调用方法
def concatKey(str1, str2):
    return concat(concat(str1, '_'), str2)


def _transferONT(fromAcct, toAcct, amount):
    """
    transfer ONT
    :param fromacct:
    :param toacct:
    :param amount:
    :return:
    """
    assert (CheckWitness(fromAcct))
    param = state(fromAcct, toAcct, amount)
    res = Invoke(0, ONTAddress, 'transfer', [param])
    if res and res == b'\x01':
        return True
    else:
        return False


def _transferONG(fromAcct, toAcct, amount):
    """
    transfer ONG
    :param fromacct:
    :param toacct:
    :param amount:
    :return:
    """
    assert (CheckWitness(fromAcct))
    param = state(fromAcct, toAcct, amount)
    res = Invoke(0, ONGAddress, 'transfer', [param])
    if res and res == b'\x01':
        return True
    else:
        return False


def _transferOEP4(oep4ReverseAddr, fromAcct, toAcct, amount):
    """
    transfer _transferOEP4
    :param fromacct:
    :param toacct:
    :param amount:
    :return:
    """
    assert (CheckWitness(fromAcct))
    params = [fromAcct, toAcct, amount]
    res = DynamicAppCall(oep4ReverseAddr, 'transfer', params)
    if res and res == b'\x01':
        return True
    else:
        return False


def _transferONTFromContact(toAcct, amount):
    param = state(ContractAddress, toAcct, amount)
    res = Invoke(0, ONTAddress, 'transfer', [param])
    if res and res == b'\x01':
        return True
    else:
        return False


def _transferONGFromContact(toAcct, amount):
    param = state(ContractAddress, toAcct, amount)
    res = Invoke(0, ONGAddress, 'transfer', [param])
    if res and res == b'\x01':
        return True
    else:
        return False


def _transferOEP4FromContact(oep4ReverseAddr, toAcct, amount):
    params = [ContractAddress, toAcct, amount]
    res = DynamicAppCall(oep4ReverseAddr, 'transfer', params)
    if res and res == b'\x01':
        return True
    else:
        return False