/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package ont

import (
	"fmt"

	"github.com/ontio/layer2/node/common"
	"github.com/ontio/layer2/node/common/config"
	"github.com/ontio/layer2/node/common/constants"
	cstates "github.com/ontio/layer2/node/core/states"
	"github.com/ontio/layer2/node/errors"
	"github.com/ontio/layer2/node/smartcontract/event"
	"github.com/ontio/layer2/node/smartcontract/service/native"
	"github.com/ontio/layer2/node/smartcontract/service/native/utils"
)

const (
	UNBOUND_TIME_OFFSET       = "unboundTimeOffset"
	TOTAL_SUPPLY_NAME         = "totalSupply"
	INIT_NAME                 = "init"
	TRANSFER_NAME             = "transfer"
	APPROVE_NAME              = "approve"
	TRANSFERFROM_NAME         = "transferFrom"
	NAME_NAME                 = "name"
	SYMBOL_NAME               = "symbol"
	DECIMALS_NAME             = "decimals"
	TOTALSUPPLY_NAME          = "totalSupply"
	BALANCEOF_NAME            = "balanceOf"
	ALLOWANCE_NAME            = "allowance"
	TOTAL_ALLOWANCE_NAME      = "totalAllowance"
	UNBOUND_ONG_TO_GOVERNANCE = "unboundOngToGovernance"
)

func AddNotifications(native *native.NativeService, contract common.Address, state *State) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{TRANSFER_NAME, state.From.ToBase58(), state.To.ToBase58(), state.Value},
		})
}
func GetToUInt64StorageItem(toBalance, value uint64) *cstates.StorageItem {
	sink := common.NewZeroCopySink(nil)
	sink.WriteUint64(toBalance + value)
	return &cstates.StorageItem{Value: sink.Bytes()}
}

func GenTotalSupplyKey(contract common.Address) []byte {
	return append(contract[:], TOTAL_SUPPLY_NAME...)
}

func GenBalanceKey(contract, addr common.Address) []byte {
	return append(contract[:], addr[:]...)
}

func Transfer(native *native.NativeService, contract common.Address, state *State) (uint64, uint64, error) {
	// this is layer2 tx
	isLayer2Deposit := IsLayer2Addr(state.From)
	if isLayer2Deposit {
		if native.Operator == false {
			return 0, 0, errors.NewErr("only operator can use layer2 deposit!")
		}
	} else {
		if !native.ContextRef.CheckWitness(state.From) {
			return 0, 0, errors.NewErr("authentication failed!")
		}
	}

	//
	isLayer2Withdraw := IsLayer2Addr(state.To)
	fromBalance := uint64(0)
	toBalance := uint64(0)
	var err error
	if isLayer2Deposit != true {
		fromBalance, err = fromTransfer(native, GenBalanceKey(contract, state.From), state.Value)
		if err != nil {
			return 0, 0, err
		}
	}

	if isLayer2Withdraw != true {
		toBalance, err = toTransfer(native, GenBalanceKey(contract, state.To), state.Value)
		if err != nil {
			return 0, 0, err
		}
	}
	return fromBalance, toBalance, nil
}

func GenApproveKey(contract, from, to common.Address) []byte {
	temp := append(contract[:], from[:]...)
	return append(temp, to[:]...)
}

func TransferedFrom(native *native.NativeService, currentContract common.Address, state *TransferFrom) (uint64, uint64, error) {
	isLayer2Deposit := IsLayer2Addr(state.From)
	if isLayer2Deposit {
		if native.Operator == false {
			return 0, 0, errors.NewErr("only operator can use layer2 deposit!")
		}
	} else {
		if native.Time <= config.GetOntHolderUnboundDeadline()+constants.GENESIS_BLOCK_TIMESTAMP {
			if !native.ContextRef.CheckWitness(state.Sender) {
				return 0, 0, errors.NewErr("authentication failed!")
			}
		} else {
			if state.Sender != state.To && !native.ContextRef.CheckWitness(state.Sender) {
				return 0, 0, errors.NewErr("authentication failed!")
			}
		}
	}

	isLayer2Withdraw := IsLayer2Addr(state.To)

	var fromBalance uint64
	var toBalance uint64
	var err error
	if !isLayer2Deposit {
		if err := fromApprove(native, genTransferFromKey(currentContract, state), state.Value); err != nil {
			return 0, 0, err
		}

		fromBalance, err = fromTransfer(native, GenBalanceKey(currentContract, state.From), state.Value)
		if err != nil {
			return 0, 0, err
		}
	}

	if !isLayer2Withdraw {
		toBalance, err = toTransfer(native, GenBalanceKey(currentContract, state.To), state.Value)
		if err != nil {
			return 0, 0, err
		}
	}
	return fromBalance, toBalance, nil
}

func getUnboundOffset(native *native.NativeService, contract, address common.Address) (uint32, error) {
	offset, err := utils.GetStorageUInt32(native, genAddressUnboundOffsetKey(contract, address))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func getGovernanceUnboundOffset(native *native.NativeService, contract common.Address) (uint32, error) {
	offset, err := utils.GetStorageUInt32(native, genGovernanceUnboundOffsetKey(contract))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func genTransferFromKey(contract common.Address, state *TransferFrom) []byte {
	temp := append(contract[:], state.From[:]...)
	return append(temp, state.Sender[:]...)
}

func fromApprove(native *native.NativeService, fromApproveKey []byte, value uint64) error {
	approveValue, err := utils.GetStorageUInt64(native, fromApproveKey)
	if err != nil {
		return err
	}
	if approveValue < value {
		return fmt.Errorf("[TransferFrom] approve balance insufficient! have %d, got %d", approveValue, value)
	} else if approveValue == value {
		native.CacheDB.Delete(fromApproveKey)
	} else {
		native.CacheDB.Put(fromApproveKey, utils.GenUInt64StorageItem(approveValue-value).ToArray())
	}
	return nil
}

func fromTransfer(native *native.NativeService, fromKey []byte, value uint64) (uint64, error) {
	fromBalance, err := utils.GetStorageUInt64(native, fromKey)
	if err != nil {
		return 0, err
	}
	if fromBalance < value {
		addr, _ := common.AddressParseFromBytes(fromKey[20:])
		return 0, fmt.Errorf("[Transfer] balance insufficient. contract:%s, account:%s,balance:%d, transfer amount:%d",
			native.ContextRef.CurrentContext().ContractAddress.ToHexString(), addr.ToBase58(), fromBalance, value)
	} else if fromBalance == value {
		native.CacheDB.Delete(fromKey)
	} else {
		native.CacheDB.Put(fromKey, utils.GenUInt64StorageItem(fromBalance-value).ToArray())
	}
	return fromBalance, nil
}

func toTransfer(native *native.NativeService, toKey []byte, value uint64) (uint64, error) {
	toBalance, err := utils.GetStorageUInt64(native, toKey)
	if err != nil {
		return 0, err
	}
	native.CacheDB.Put(toKey, GetToUInt64StorageItem(toBalance, value).ToArray())
	return toBalance, nil
}

func IsLayer2Addr(addr common.Address) bool {
	if addr.ToHexString() == common.ADDRESS_EMPTY.ToHexString() {
		return true
	} else {
		return false
	}
}

func genAddressUnboundOffsetKey(contract, address common.Address) []byte {
	temp := append(contract[:], UNBOUND_TIME_OFFSET...)
	return append(temp, address[:]...)
}

func genGovernanceUnboundOffsetKey(contract common.Address) []byte {
	temp := append(contract[:], UNBOUND_TIME_OFFSET...)
	return temp
}
