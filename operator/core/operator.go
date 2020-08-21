/*
 * Copyright (C) 2020 The ontology Authors
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


package core

import (
	"encoding/hex"
	"fmt"
	"github.com/ontio/layer2/operator/config"
	"github.com/ontio/layer2/operator/log"
	ontology_sdk "github.com/ontio/ontology-go-sdk"
	ontology_sdk_common "github.com/ontio/ontology-go-sdk/common"
	ontology_common "github.com/ontio/ontology/common"
	ontology_types "github.com/ontio/ontology/core/types"
	"time"
)

type Layer2Operator struct {
	config             *config.ServiceConfig

	ontologySdk        *ontology_sdk.OntologySdk
	ontologyAccount    *ontology_sdk.Account
	ontologyChainInfo  *ChainInfo

	layer2Sdk          *ontology_sdk.Layer2Sdk
	layer2Account      *ontology_sdk.Account
	layer2ChainInfo    *ChainInfo

	depositChain        chan []*Deposit
	msgChan             chan []*Layer2CommitMsg
	exitChan            chan int
}

func NewLayer2Operator(servCfg *config.ServiceConfig) (*Layer2Operator, error) {
	ontologySdk := ontology_sdk.NewOntologySdk()
	ontologySdk.NewRpcClient().SetAddress(servCfg.OntologyConfig.RestURL)
	layer2Sdk := ontology_sdk.NewLayer2Sdk()
	layer2Sdk.NewRpcClient().SetAddress(servCfg.Layer2Config.RestURL)
	return &Layer2Operator{
		exitChan:           make(chan int),
		depositChain:       make(chan []*Deposit),
		msgChan:            make(chan []*Layer2CommitMsg),
		config:             servCfg,
		ontologySdk:        ontologySdk,
		layer2Sdk:          layer2Sdk,
	}, nil
}

func (this *Layer2Operator) getOntologyAccount() (*ontology_sdk.Account, error) {
	var wallet *ontology_sdk.Wallet
	var err error
	if !ontology_common.FileExisted(this.config.OntologyConfig.WalletFile) {
		wallet, err = this.ontologySdk.CreateWallet(this.config.OntologyConfig.WalletFile)
		if err != nil {
			return nil, err
		}
	} else {
		wallet, err = this.ontologySdk.OpenWallet(this.config.OntologyConfig.WalletFile)
		if err != nil {
			log.Errorf("ontologyAccount - wallet open error: %s", err.Error())
			return nil, err
		}
	}
	signer, err := wallet.GetDefaultAccount([]byte(this.config.OntologyConfig.WalletPwd))
	if err != nil || signer == nil {
		signer, err = wallet.NewDefaultSettingAccount([]byte(this.config.OntologyConfig.WalletPwd))
		if err != nil {
			log.Errorf("ontologyAccount - wallet password error")
			return nil, err
		}

		err = wallet.Save()
		if err != nil {
			return nil, err
		}
	}
	log.Infof("ontologyAccount - ont account address: %s, %s", signer.Address.ToBase58(), signer.Address.ToHexString())
	return signer, nil
}

func (this *Layer2Operator) getLyer2Account() (*ontology_sdk.Account, error) {
	var wallet *ontology_sdk.Wallet
	var err error
	if !ontology_common.FileExisted(this.config.Layer2Config.WalletFile) {
		wallet, err = this.layer2Sdk.CreateWallet(this.config.Layer2Config.WalletFile)
		if err != nil {
			return nil, err
		}
	} else {
		wallet, err = this.layer2Sdk.OpenWallet(this.config.Layer2Config.WalletFile)
		if err != nil {
			log.Errorf("layer2Account - wallet open error: %s", err.Error())
			return nil, err
		}
	}
	signer, err := wallet.GetDefaultAccount([]byte(this.config.Layer2Config.WalletPwd))
	if err != nil || signer == nil {
		signer, err = wallet.NewDefaultSettingAccount([]byte(this.config.Layer2Config.WalletPwd))
		if err != nil {
			log.Errorf("layer2Account - wallet password error")
			return nil, err
		}

		err = wallet.Save()
		if err != nil {
			return nil, err
		}
	}
	log.Infof("layer2Account - layer2 account address: %s, %s", signer.Address.ToBase58(), signer.Address.ToHexString())
	return signer, nil
}

func (this *Layer2Operator) Start() error {
	// try to connect db
	dberr := ConnectDB(this.config.DBConfig.ProjectDBUser, this.config.DBConfig.ProjectDBPassword, this.config.DBConfig.ProjectDBUrl, this.config.DBConfig.ProjectDBName)
	if dberr != nil {
		return fmt.Errorf(dberr.Error())
	}

	//  try to load all chains
	ontologyChain := LoadChainInfo("ontology")
	if ontologyChain == nil {
		return fmt.Errorf("load multichain info error")
	}
	this.ontologyChainInfo = ontologyChain

	layer2Chain := LoadChainInfo("layer2")
	if layer2Chain == nil {
		return fmt.Errorf("load ontology chain info error")
	}
	this.layer2ChainInfo = layer2Chain
	
	ontologyAccount, err := this.getOntologyAccount()
	if err != nil {
		return err
	}
	layer2Account, err := this.getLyer2Account()
	if err != nil {
		return err
	}
	this.ontologyAccount = ontologyAccount
	this.layer2Account = layer2Account
	//
	{
		currentHeight, err := this.ontologySdk.GetCurrentBlockHeight()
		if err != nil {
			log.Errorf("get ontology current block heigh err: %s", err.Error())
		} else {
			if this.ontologyChainInfo.Height <= 0 {
				this.ontologyChainInfo.Height = currentHeight
			}
		}
		log.Infof("ontology current height: %d", this.ontologyChainInfo.Height)
	}
	{
		currentHeight, err := this.GetLayer2CommitHeight()
		if err != nil {
			log.Errorf("get layer2 current block heigh err: %s", err.Error())
		} else {
			this.layer2ChainInfo.Height = currentHeight
		}
		log.Infof("layer2 current height: %d", this.ontologyChainInfo.Height)
	}
	go this.MonitorOntologyChain()
	go this.MonitorLayer2Chain()
	go this.depositLoop()
	go this.commitMsgLoop()
	return nil
}

func (this *Layer2Operator) Stop() {
	this.exitChan <- 1
	this.exitChan <- 1
	close(this.exitChan)
	CloseDB()
	log.Infof("multi chain manager exit.")
}

func (this *Layer2Operator) MonitorOntologyChain() {
	log.Infof("start MonitorOntologyChain")
	updateTicker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <- updateTicker.C:
			currentHeight, err := this.ontologySdk.GetCurrentBlockHeight()
			if err != nil {
				log.Errorf("get ontology chain current height err: %s", err.Error())
				continue
			}
			log.Infof("chain %s current height: %d, parser height: %d", this.ontologyChainInfo.Name, currentHeight, this.ontologyChainInfo.Height)
			if currentHeight <= this.ontologyChainInfo.Height {
				continue
			}
			for currentHeight > this.ontologyChainInfo.Height {
				this.ontologyChainInfo.Height ++
				deposits, err := this.parseOntologyChainBlock(this.ontologyChainInfo)
				if err != nil {
					log.Errorf("parse ontology chain block err: %s", err.Error())
					this.ontologyChainInfo.Height --
					break
				}
				SetChainParseHeight(this.ontologyChainInfo.Id, this.ontologyChainInfo.Height)
				this.depositChain <- deposits
			}
		case <- this.exitChan:
			updateTicker.Stop()
			log.Infof("chain %s, exit!", this.ontologyChainInfo.Name)
			return
		}
	}
}

func (this *Layer2Operator) parseOntologyChainBlock(chain *ChainInfo) ([]*Deposit, error) {
	block, err := this.ontologySdk.GetBlockByHeight(chain.Height)
	if err != nil {
		return nil, err
	}
	tt := block.Header.Timestamp
	events, err := this.ontologySdk.GetSmartContractEventByBlock(chain.Height)
	if err != nil {
		return nil, err
	}
	deposits := make([]*Deposit, 0)
	for _, event := range events {
		//log.Infof("tx hash: %s, state:%d, gas: %d", event.TxHash, event.State, event.GasConsumed)
		for _, notify := range event.Notify {
			if notify.ContractAddress != this.config.OntologyConfig.Layer2ContractAddress {
				continue
			}
			//
			states := notify.States.([]interface{})
			method, _ := hex.DecodeString(states[0].(string))
			log.Infof("find layer2 transaction: %s, method: %s", event.TxHash, string(method))
			if string(method) == "deposit" {
				id, _ := hex.DecodeString(states[1].(string))
				player := revertHexString(states[2].(string))
				playerAddr, _ := ontology_common.AddressFromHexString(player)
				amount, _ := hex.DecodeString(states[3].(string))

				deposit := &Deposit{}
				deposit.TxHash = event.TxHash
				deposit.TT = tt
				deposit.Height = chain.Height
				deposit.State = DEPOSIT_EVENT
				deposit.FromAddress = playerAddr.ToBase58()
				deposit.Amount = BytesToInt(amount)
				deposit.TokenAddress = states[6].(string)
				deposit.ID = BytesToInt(id)
				err = SaveDeposit(deposit)
				if err != nil {
					log.Warnf("save deposit tx error: %v", err)
				}
				if deposit.TokenAddress == ONG_CONTRACT_ADDRESS && deposit.Amount < this.config.Layer2Config.MinOngLimit {
					log.Warnf("the deposit ong is too small!")
					continue
				}
				deposits = append(deposits, deposit)
			}
		}
	}
	return deposits, nil
}

func (this *Layer2Operator) depositLoop() {
	log.Infof("start depositLoop")
	for {
		select {
		case deposits := <-this.depositChain:
			for _, deposit := range deposits {
				err := this.commitDeposit2Layer2(deposit)
				if err != nil {
					deposit.State = DEPOSIT_FAILED
					formatStr := "2006-01-02 15:04:05"
					timehash := time.Now().Format(formatStr)
					UpdateDepositByID(deposit.ID, deposit.State, timehash)
					log.Errorf("commit deposit to layer2, from : %s, to : %s, hash: %s, err: %s", ontology_common.ADDRESS_EMPTY.ToBase58(), deposit.FromAddress, timehash, err.Error())
				}
			}
		}
	}
}

func (this *Layer2Operator) commitDeposit2Layer2(deposit *Deposit) error {
	log.Infof("commit deposit to layer2: %s", deposit.Dump())
	var tx *ontology_types.MutableTransaction
	var err error
	toAddr, err := ontology_common.AddressFromBase58(deposit.FromAddress)
	if err != nil {
		return fmt.Errorf("to address is not right!")
	}
	if deposit.TokenAddress == ONT_CONTRACT_ADDRESS {
		tx, err = this.layer2Sdk.Native.Ont.NewTransferTransaction(0, 20000, ontology_common.ADDRESS_EMPTY, toAddr, deposit.Amount)
		if err != nil {
			return err
		}
	} else if deposit.TokenAddress == ONG_CONTRACT_ADDRESS {
		tx, err = this.layer2Sdk.Native.Ong.NewTransferTransaction(0, 20000, ontology_common.ADDRESS_EMPTY, toAddr, deposit.Amount)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("token is not supported!")
	}
	this.layer2Sdk.SetPayer(tx, this.layer2Account.Address)
	err = this.layer2Sdk.SignToTransaction(tx, this.layer2Account)
	if err != nil {
		return err
	}
	var hash ontology_common.Uint256
	counter := 0
	for true {
		hash, err = this.layer2Sdk.SendTransaction(tx)
		if err != nil {
			log.Errorf("send transaction err when commit deposit 2 layer2, err: %s, try again......", err.Error())
			if counter == LAYER_TRANSACTION_RETRY {
				break
			}
			time.Sleep(time.Second * 1)
			counter ++
			// send error, we cannot send again, so ignore this error
		} else {
			break
		}
	}
	if counter == LAYER_TRANSACTION_RETRY {
		return fmt.Errorf("commit deposit 2 layer2 error!")
	}
	deposit.State = DEPOSIT_COMMIT
	UpdateDepositByID(deposit.ID, deposit.State, hash.ToHexString())
	log.Infof("commit deposit to layer2, from : %s, to : %s, tx hash: %s", ontology_common.ADDRESS_EMPTY.ToBase58(), toAddr.ToBase58(), hash.ToHexString())
	return nil
}

func (this *Layer2Operator) MonitorLayer2Chain() {
	log.Infof("start MonitorLayer2Chain")
	updateTicker := time.NewTicker(time.Second * time.Duration(this.config.Layer2Config.BlockDuration))
	for {
		select {
		case <- updateTicker.C:
			currentHeight, err := this.layer2Sdk.GetCurrentBlockHeight()
			if err != nil {
				log.Errorf("get layer2 current block height err: %s", err.Error())
				continue
			}
			log.Infof("chain %s current height: %d, parser height: %d", this.layer2ChainInfo.Name, currentHeight, this.layer2ChainInfo.Height)
			commitHeight, err := this.GetLayer2CommitHeight()
			if err != nil {
				log.Errorf("get layer2 commit height err: %s", err.Error())
				continue
			}
			this.layer2ChainInfo.Height = commitHeight
			if this.layer2ChainInfo.Height >= currentHeight {
				continue
			}
			commitMsgs := make([]*Layer2CommitMsg, 0)
			for this.layer2ChainInfo.Height < currentHeight {
				this.layer2ChainInfo.Height ++
				layer2CommitMsg, err := this.parseLayer2ChainBlock(this.layer2ChainInfo)
				if err != nil {
					log.Errorf("parser layer2 chain block err: %s", err.Error())
					break
				} else {
					commitMsgs = append(commitMsgs, layer2CommitMsg)
				}
			}
			SetChainParseHeight(this.layer2ChainInfo.Id, this.layer2ChainInfo.Height)
			this.msgChan <- commitMsgs
		case <- this.exitChan:
			updateTicker.Stop()
			log.Infof("chain %s, exit!", this.layer2ChainInfo.Name)
			return
		}
	}
}

func (this *Layer2Operator) parseLayer2ChainBlock(chain *ChainInfo) (*Layer2CommitMsg, error) {
	block, err := this.layer2Sdk.GetLayer2BlockByHeight(chain.Height)
	if err != nil {
		return nil, err
	}
	tt := block.Header.Timestamp

	events, err := this.layer2Sdk.GetSmartContractEventByBlock(chain.Height)
	if err != nil {
		return nil, err
	}
	msg := &Layer2CommitMsg{}
	insertLayer2TxBatch := NewMysqlInsertBatch(DefDB, 9, "(?,?,?,?,?,?,?,?,?)", "insert into layer2tx(txhash, tt, state, fee, height, fromaddress, tokenaddress, toaddress, amount)")
	insertLayer2TxArgs := make([]interface{}, 9)
	updateDepositBatch := NewMysqlUpdateBatch(DefDB, 9, "(?,?,?,?,?,?,?,?,?)", "insert into deposit(txhash, tt, state, height, fromaddress, amount, tokenaddress, id, layer2txhash)", "ON DUPLICATE KEY UPDATE state=VALUES(state)")
	updateDepositArgs := make([]interface{}, 9)
	insertWithdrawBatch := NewMysqlInsertBatch(DefDB, 7, "(?,?,?,?,?,?,?)", "insert into withdraw(txhash, tt, state, height, toaddress, amount, tokenaddress)")
	insertWithdrawArgs := make([]interface{}, 7)
	// log.Infof("chain: %s, block height: %d, events num: %d\n", chain.Name, chain.Height, len(events))
	for _, event := range events {
		//log.Infof("tx hash: %s, state:%d, gas: %d\n", event.TxHash, event.State, event.GasConsumed)
		for _, notify := range event.Notify {
			if notify.ContractAddress != ONT_REV_CONTRACT_ADDRESS && notify.ContractAddress != ONG_REV_CONTRACT_ADDRESS  {
				continue
			}
			states := notify.States.([]interface{})
			if len(states) != 4 {
				continue
			}
			if states[0] != NOTIFY_TRANSFER {
				continue
			}
			transferFrom, ok := states[1].(string)
			if !ok {
				continue
			}
			transferTo, ok := states[2].(string)
			if !ok {
				continue
			}
			transferAmount, ok := states[3].(uint64)
			if !ok {
				continue
			}

			layer2Tx := &Layer2Tx{}
			layer2Tx.TxHash = event.TxHash
			layer2Tx.TT = tt
			layer2Tx.Fee = 0
			layer2Tx.Height = chain.Height
			layer2Tx.State = 1
			layer2Tx.FromAddress = transferFrom
			layer2Tx.Amount = transferAmount
			layer2Tx.TokenAddress = revertHexString(notify.ContractAddress)
			layer2Tx.ToAddress = transferTo
			insertLayer2TxArgs[0] = layer2Tx.TxHash
			insertLayer2TxArgs[1] = layer2Tx.TT
			insertLayer2TxArgs[2] = layer2Tx.State
			insertLayer2TxArgs[3] = layer2Tx.Fee
			insertLayer2TxArgs[4] = layer2Tx.Height
			insertLayer2TxArgs[5] = layer2Tx.FromAddress
			insertLayer2TxArgs[6] = layer2Tx.TokenAddress
			insertLayer2TxArgs[7] = layer2Tx.ToAddress
			insertLayer2TxArgs[8] = layer2Tx.Amount
			err := insertLayer2TxBatch.Insert(insertLayer2TxArgs)
			if err != nil {
				log.Warnf("insert layer2 tx err: %s", err.Error())
			}
			//
			if isLayer2Tx(layer2Tx.FromAddress) {
				deposit := LoadDepositByLayer2TxHash(layer2Tx.TxHash)
				msg.Deposits = append(msg.Deposits, deposit.ID)
				updateDepositArgs[0] = ""
				updateDepositArgs[1] = 0
				updateDepositArgs[2] = DEPOSIT_FINISH
				updateDepositArgs[3] = 0
				updateDepositArgs[4] = ""
				updateDepositArgs[5] = 0
				updateDepositArgs[6] = ""
				updateDepositArgs[7] = deposit.ID
				updateDepositArgs[8] = ""
				updateDepositBatch.Insert(updateDepositArgs)
			}
			//
			if isLayer2Tx(layer2Tx.ToAddress) {
				withdraw := &Withdraw{}
				withdraw.TxHash = event.TxHash
				withdraw.TT = tt
				withdraw.Height = chain.Height
				withdraw.State = WITHDRAW_INIT
				withdraw.ToAddress = transferFrom
				withdraw.Amount = transferAmount
				withdraw.TokenAddress = revertHexString(notify.ContractAddress)
				insertWithdrawArgs[0] = withdraw.TxHash
				insertWithdrawArgs[1] = withdraw.TT
				insertWithdrawArgs[2] = withdraw.State
				insertWithdrawArgs[3] = withdraw.Height
				insertWithdrawArgs[4] = withdraw.ToAddress
				insertWithdrawArgs[5] = withdraw.Amount
				insertWithdrawArgs[6] = withdraw.TokenAddress
				insertWithdrawBatch.Insert(insertWithdrawArgs)
				msg.WithDraws = append(msg.WithDraws, withdraw)
			}
		}
	}
	insertLayer2TxBatch.Close()
	updateDepositBatch.Close()
	insertWithdrawBatch.Close()
	//
	layer2State := &Layer2State{
		Height: block.Header.Height,
		Version: 1,
		StatesRoot: block.Header.StateRoot,
	}
	msg.Layer2State = layer2State
	return msg, nil
}

func (this *Layer2Operator) commitMsgLoop() {
	log.Infof("start commitMsgLoop")
	for {
		select {
		case msgs := <-this.msgChan:
			start := 0
			end := 0
			for start < len(msgs) {
				end = start + int(this.config.Layer2Config.MaxBatchSize)
				if end > len(msgs) {
					end = len(msgs)
				 }
				msgBatch := msgs[start : end]
				for true {
					err := this.commitLayer2State2Ontology(msgBatch)
					if err != nil {
						log.Errorf("commit layer2 state to ontology err: %s", err.Error())
						time.Sleep(time.Second * 1)
					} else {
						break
					}
				}
				start = end
			}
		}
	}
}

func (this *Layer2Operator) commitLayer2State2Ontology(msgs []*Layer2CommitMsg) error {
	contractAddress, _ := ontology_common.AddressFromHexString(this.config.OntologyConfig.Layer2ContractAddress)
	stateRootsBatch := make([]string, 0)
	heightsBatch := make([]uint32, 0)
	versionsBatch := make([]string, 0)
	depositidsBatch := make([][]uint64, 0)
	withdrawAmountsBatch := make([][]uint64, 0)
	toAddressesBatch := make([][]ontology_common.Address, 0)
	assetAddressBatch := make([][][]byte, 0)
	cacheMsgs := make([]*Layer2CommitMsg, 0)
	for _, msg := range msgs {
		layer2Msg := msg.Dump()
		log.Infof("commit layer2 state to ontology: %s", layer2Msg)
		//
		cacheMsgs = append(cacheMsgs, msg)
		depositids := make([]uint64, 0)
		for _, id := range msg.Deposits {
			depositids = append(depositids, id)
		}
		withdrawAmounts := make([]uint64, 0)
		toAddresses := make([]ontology_common.Address, 0)
		assetAddress := make([][]byte, 0)
		for _, withdraw := range msg.WithDraws {
			withdrawAmounts = append(withdrawAmounts, withdraw.Amount)
			toAddress, _ := ontology_common.AddressFromBase58(withdraw.ToAddress)
			toAddresses = append(toAddresses, toAddress)
			tokenAddress, _ := hex.DecodeString(withdraw.TokenAddress)
			assetAddress = append(assetAddress, tokenAddress)
		}
		stateRootsBatch = append(stateRootsBatch, msg.Layer2State.StatesRoot.ToHexString())
		heightsBatch = append(heightsBatch, msg.Layer2State.Height)
		versionsBatch = append(versionsBatch, string(msg.Layer2State.Version))
		depositidsBatch = append(depositidsBatch, depositids)
		withdrawAmountsBatch = append(withdrawAmountsBatch, withdrawAmounts)
		toAddressesBatch = append(toAddressesBatch, toAddresses)
		assetAddressBatch = append(assetAddressBatch, assetAddress)
		/*
		result, err := this.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"updateState", []interface{}{
			msg.Layer2State.StatesRoot.ToHexString(), msg.Layer2State.Height, string(msg.Layer2State.Version),
			depositids, withdrawAmounts, toAddresses, assetAddress}})
		*/
	}
	result, err := this.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"updateStates", []interface{}{
		stateRootsBatch, heightsBatch, versionsBatch, depositidsBatch, withdrawAmountsBatch, toAddressesBatch, assetAddressBatch}})
	var gasLimit uint64
	if err != nil {
		log.Warnf("Pre exec contract err: %s", err.Error())
		gasLimit = this.config.OntologyConfig.GasLimit
	} else {
		gasLimit = result.Gas
	}
	tx, err := this.ontologySdk.NeoVM.NewNeoVMInvokeTransaction(this.config.OntologyConfig.GasPrice, gasLimit, contractAddress, []interface{}{"updateStates", []interface{}{
		stateRootsBatch, heightsBatch, versionsBatch, depositidsBatch, withdrawAmountsBatch, toAddressesBatch, assetAddressBatch}})
	if err != nil {
		return fmt.Errorf("new layer2 state commit transaction failed! err: %s", err.Error())
	}
	this.ontologySdk.SetPayer(tx, this.ontologyAccount.Address)
	err = this.ontologySdk.SignToTransaction(tx, this.ontologyAccount)
	if err != nil {
		return fmt.Errorf("sign layer2 state commit transaction failed! err: %s", err.Error())
	}
	var txHash ontology_common.Uint256
	for true {
		txHash, err = this.ontologySdk.SendTransaction(tx)
		if err != nil {
			log.Errorf("send layer2 state commit transaction failed! err: %s, try again......", err.Error())
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}
	log.Infof("layer2 state commit transaction hash: %s", txHash.ToHexString())

	//
	for _, msg := range cacheMsgs {
		for _, id := range msg.Deposits {
			UpdateDepositByID2(id, DEPOSIT_NOTIFY)
		}
		for _, withdraw := range msg.WithDraws {
			UpdateWithdraw(withdraw.TxHash, WITHDRAW_COMMIT, txHash.ToHexString())
		}
		SaveLayer2Commit(txHash.ToHexString(), msg.Dump1(), uint64(msg.Layer2State.Height))
	}
	return nil
}

func (this *Layer2Operator) PreExecInvokeNeoVMContract(contractAddress ontology_common.Address, params []interface{}) (*ontology_sdk_common.PreExecResult, error) {
	tx, err := this.ontologySdk.NeoVM.NewNeoVMInvokeTransaction(0, 0, contractAddress, params)
	if err != nil {
		return nil, err
	}
	this.ontologySdk.SetPayer(tx, this.ontologyAccount.Address)
	err = this.ontologySdk.SignToTransaction(tx, this.ontologyAccount)
	if err != nil {
		return nil, fmt.Errorf("sign layer2 state commit transaction failed! err: %s", err.Error())
	}
	return this.ontologySdk.PreExecTransaction(tx)
}

func (this *Layer2Operator) GetLayer2CommitHeight() (uint32, error) {
	contractAddress, _ := ontology_common.AddressFromHexString(this.config.OntologyConfig.Layer2ContractAddress)
	tx, err := this.ontologySdk.NeoVM.NewNeoVMInvokeTransaction(0, 0, contractAddress, []interface{}{"getCurrentHeight", []interface{}{}})
	if err != nil {
		return 0, err
	}
	result, err := this.ontologySdk.PreExecTransaction(tx)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, fmt.Errorf("can not find the result")
	}
	height, _ := result.Result.ToInteger()
	return uint32(height.Uint64()), nil
}

func isLayer2Tx(addr string) bool {
	newAddr,_ := ontology_common.AddressFromBase58(addr)
	if newAddr.ToHexString() == ontology_common.ADDRESS_EMPTY.ToHexString() {
		return true
	} else {
		return false
	}
}