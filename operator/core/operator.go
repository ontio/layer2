package core

import (
	"encoding/hex"
	"fmt"
	layer2_sdk "github.com/ontio/layer2/go-sdk"
	"github.com/ontio/layer2/operator/config"
	"github.com/ontio/layer2/operator/log"
	layer2_common "github.com/ontio/layer2/node/common"
	layer2_types "github.com/ontio/layer2/node/core/types"
	ontology_sdk "github.com/ontio/ontology-go-sdk"
	ontology_common "github.com/ontio/ontology/common"
	"math/rand"
	"time"
)

type Layer2Operator struct {
	config             *config.ServiceConfig

	ontologySdk        *ontology_sdk.OntologySdk
	ontologyAccount    *ontology_sdk.Account
	ontologyChainInfo  *ChainInfo

	layer2Sdk          *layer2_sdk.OntologySdk
	layer2Account      *layer2_sdk.Account
	layer2ChainInfo    *ChainInfo

	depositChain        chan *Deposit
	msgChan             chan *Layer2CommitMsg
	exitChan            chan int

	// use for test
	fortest              int
	deposit              int
	withdraw             int
	depositHeight        uint32
}

func NewLayer2Operator(servCfg *config.ServiceConfig) (*Layer2Operator, error) {
	ontologySdk := ontology_sdk.NewOntologySdk()
	ontologySdk.NewRpcClient().SetAddress(servCfg.OntologyConfig.RestURL)
	layer2Sdk := layer2_sdk.NewOntologySdk()
	layer2Sdk.NewRpcClient().SetAddress(servCfg.Layer2Config.RestURL)
	return &Layer2Operator{
		exitChan:           make(chan int),
		depositChain:       make(chan *Deposit),
		msgChan:            make(chan *Layer2CommitMsg),
		config:             servCfg,
		ontologySdk:        ontologySdk,
		layer2Sdk:          layer2Sdk,
		fortest:            0,
		deposit:            0,
		withdraw:           0,
		depositHeight:      0,
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

func (this *Layer2Operator) getLyer2Account() (*layer2_sdk.Account, error) {
	var wallet *layer2_sdk.Wallet
	var err error
	if !layer2_common.FileExisted(this.config.Layer2Config.WalletFile) {
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
			fmt.Println(err)
		} else {
			if this.ontologyChainInfo.Height <= 0 {
				this.ontologyChainInfo.Height = currentHeight
			}
		}
		log.Infof("ontology current height: %d", this.ontologyChainInfo.Height)
	}
	/*
	{
		currentHeight, err := this.layer2Sdk.GetCurrentBlockHeight()
		if err != nil {
			fmt.Println(err)
		} else {
			if this.layer2ChainInfo.Height <= 0 {
				this.layer2ChainInfo.Height = currentHeight
			}
		}
	}
	 */
	{
		currentHeight := GetLastestLayer2Commit()
		this.layer2ChainInfo.Height = currentHeight
		log.Infof("layer2 current height: %d", this.layer2ChainInfo.Height)
	}

	go this.MonitorOntologyChain()
	go this.MonitorLayer2Chain()
	go this.depositLoop()
	go this.commitMsgLoop()
	go this.checkMsgLoop()
	if this.fortest == 1 {
		go this.testLoop()
	}
	return nil
}

func (this *Layer2Operator) Stop() {
	this.exitChan <- 1
	close(this.exitChan)
	log.Infof("multi chain manager exit.")
}

func (this *Layer2Operator) MonitorOntologyChain() {
	updateTicker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <- updateTicker.C:
			currentHeight, err := this.ontologySdk.GetCurrentBlockHeight()
			if err != nil {
				fmt.Println(err)
				continue
			}
			log.Infof("chain %s current height: %d, parser height: %d", this.ontologyChainInfo.Name, currentHeight, this.ontologyChainInfo.Height)
			if currentHeight <= this.ontologyChainInfo.Height {
				continue
			}
			for currentHeight > this.ontologyChainInfo.Height {
				this.ontologyChainInfo.Height ++
				err = this.parseOntologyChainBlock(this.ontologyChainInfo)
				if err != nil {
					fmt.Println(err)
					this.ontologyChainInfo.Height --
					break
				}
				SetChainParseHeight(this.ontologyChainInfo.Id, this.ontologyChainInfo.Height)
			}
		case <- this.exitChan:
			updateTicker.Stop()
			log.Infof("chain %s, exit!", this.ontologyChainInfo.Name)
			return
		}
	}
}

func (this *Layer2Operator) parseOntologyChainBlock(chain *ChainInfo) error {
	block, err := this.ontologySdk.GetBlockByHeight(chain.Height)
	if err != nil {
		return err
	}
	tt := block.Header.Timestamp

	events, err := this.ontologySdk.GetSmartContractEventByBlock(chain.Height)
	if err != nil {
		return err
	}

	//log.Infof("chain: %s, block height: %d, events num: %d", chain.Name, chain.Height, len(events))
	for _, event := range events {
		//log.Infof("tx hash: %s, state:%d, gas: %d", event.TxHash, event.State, event.GasConsumed)
		for _, notify := range event.Notify {
			if notify.ContractAddress != this.config.OntologyConfig.Layer2ContractAddress {
				continue
			}
			// todo
			states := notify.States.([]interface{})
			method, _ := hex.DecodeString(states[0].(string))
			log.Infof("find layer2 deposit transaction: %s, method: %s", event.TxHash, string(method))
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
					log.Errorf("save deposit tx error: %v", err)
				}
				//
				this.depositChain <- deposit
			}
		}
	}

	if this.fortest == 1 {
		rand.Seed(time.Now().UnixNano())
		t := rand.Intn(20)
		if t == 0 && this.deposit == 0 {
			this.deposit = 1
			this.depositHeight = this.layer2ChainInfo.Height
			{
				deposit := &Deposit{}
				deposit.TxHash = fmt.Sprintf("%d", time.Now().Unix())
				deposit.TT = uint32(time.Now().Unix())
				deposit.Height = 0
				deposit.State = DEPOSIT_EVENT
				deposit.FromAddress = this.layer2Account.Address.ToBase58()
				deposit.Amount = 100000
				deposit.TokenAddress = ONT_CONTRACT_ADDRESS
				deposit.ID = uint64(time.Now().Unix())
				err = SaveDeposit(deposit)
				if err != nil {
					log.Errorf("save deposit tx error: %v", err)
				}
				//
				this.depositChain <- deposit
			}
			{
				deposit := &Deposit{}
				deposit.TxHash = fmt.Sprintf("%d", time.Now().Unix()+1)
				deposit.TT = uint32(time.Now().Unix()) + 1
				deposit.Height = 0
				deposit.State = DEPOSIT_EVENT
				deposit.FromAddress = this.layer2Account.Address.ToBase58()
				deposit.Amount = 100000
				deposit.TokenAddress = ONG_CONTRACT_ADDRESS
				deposit.ID = uint64(time.Now().Unix()) + 1
				err = SaveDeposit(deposit)
				if err != nil {
					log.Errorf("save deposit tx error: %v", err)
				}
				//
				this.depositChain <- deposit
			}
		}
	}
	return nil
}

func (this *Layer2Operator) depositLoop() {
	for {
		select {
		case deposit := <-this.depositChain:
			for true {
				err := this.commitDeposit2Layer2(deposit)
				if err != nil {
					log.Errorf("commit deposit 2 alyer2 error: %s", err.Error())
				} else {
					break
				}
			}
		}
	}
}

func (this *Layer2Operator) commitDeposit2Layer2(deposit *Deposit) error {
	log.Infof("commit deposit to layer2: %s", deposit.Dump())
	toAddr, _ := layer2_common.AddressFromBase58(deposit.FromAddress)
	var tx *layer2_types.MutableTransaction
	var err error
	if deposit.TokenAddress == ONT_CONTRACT_ADDRESS {
		tx, err = this.layer2Sdk.Native.Ont.NewTransferTransaction(0, 20000, layer2_common.ADDRESS_EMPTY, toAddr, deposit.Amount)
		if err != nil {
			return err
		}
	} else if deposit.TokenAddress == ONG_CONTRACT_ADDRESS {
		tx, err = this.layer2Sdk.Native.Ong.NewTransferTransaction(0, 20000, layer2_common.ADDRESS_EMPTY, toAddr, deposit.Amount)
		if err != nil {
			return err
		}
	}

	this.layer2Sdk.SetPayer(tx, this.layer2Account.Address)
	err = this.layer2Sdk.SignToTransaction(tx, this.layer2Account)
	if err != nil {
		return err
	}
	var hash layer2_common.Uint256
	for true {
		hash, err = this.layer2Sdk.SendTransaction(tx)
		if err != nil {
			log.Errorf("send transaction err when commit deposit 2 layer2, err: %s, try again......", err.Error())
			time.Sleep(time.Second * 1)
			// send error, we cannot send again, so ignore this error
		} else {
			break
		}
	}
	deposit.State = DEPOSIT_COMMIT
	UpdateDeposit(deposit.TxHash, deposit.State, hash.ToHexString())
	log.Infof("commit deposit to layer2, from : %s, to : %s, tx hash: %s", layer2_common.ADDRESS_EMPTY.ToBase58(), toAddr.ToBase58(), hash.ToHexString())
	return nil
}

func (this *Layer2Operator) MonitorLayer2Chain() {
	updateTicker := time.NewTicker(time.Second * 1)
	counter := 0
	for {
		select {
		case <- updateTicker.C:
			currentHeight, err := this.layer2Sdk.GetCurrentBlockHeight()
			if err != nil {
				fmt.Println(err)
				continue
			}

			log.Infof("chain %s current height: %d, parser height: %d", this.layer2ChainInfo.Name, currentHeight, this.layer2ChainInfo.Height)
			if this.layer2ChainInfo.Height >= currentHeight {
				continue
			}
			for this.layer2ChainInfo.Height < currentHeight-1 {
				if counter > 600 {
					this.layer2ChainInfo.Height --
				}
				commitHeight := GetLastestLayer2Commit()
				if commitHeight < this.layer2ChainInfo.Height {
					counter ++
					break
				}

				this.layer2ChainInfo.Height ++
				counter = 0
				err = this.parseLayer2ChainBlock(this.layer2ChainInfo)
				if err != nil {
					fmt.Println(err)
					this.layer2ChainInfo.Height --
					break
				}
				SetChainParseHeight(this.layer2ChainInfo.Id, this.layer2ChainInfo.Height)
			}
		case <- this.exitChan:
			updateTicker.Stop()
			log.Infof("chain %s, exit!", this.layer2ChainInfo.Name)
			return
		}
	}
}

func (this *Layer2Operator) parseLayer2ChainBlock(chain *ChainInfo) error {
	block, err := this.layer2Sdk.GetBlockByHeight(chain.Height)
	if err != nil {
		return err
	}
	tt := block.Header.Timestamp

	events, err := this.layer2Sdk.GetSmartContractEventByBlock(chain.Height)
	if err != nil {
		return err
	}
	msg := &Layer2CommitMsg{}
	log.Infof("chain: %s, block height: %d, events num: %d\n", chain.Name, chain.Height, len(events))
	for _, event := range events {
		log.Infof("tx hash: %s, state:%d, gas: %d\n", event.TxHash, event.State, event.GasConsumed)
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
			err = SaveLayer2Tx(layer2Tx)
			if err != nil {
				log.Errorf("save layer2 tx error: %v", err)
			}

			//
			if isLayer2Tx(layer2Tx.FromAddress) {
				UpdateDepositByLayer2TxHash(layer2Tx.TxHash, DEPOSIT_FINISH)
				deposit := LoadDepositByLayer2TxHash(layer2Tx.TxHash)
				msg.Deposits = append(msg.Deposits, deposit.ID)
			}

			if isLayer2Tx(layer2Tx.ToAddress) {
				withdraw := &Withdraw{}
				withdraw.TxHash = event.TxHash
				withdraw.TT = tt
				withdraw.Height = chain.Height
				withdraw.State = WITHDRAW_INT
				withdraw.ToAddress = transferFrom
				withdraw.Amount = transferAmount
				withdraw.TokenAddress = revertHexString(notify.ContractAddress)
				err = SaveWithdraw(withdraw)
				if err != nil {
					log.Errorf("save withdraw tx error: %v", err)
				}
				msg.WithDraws = append(msg.WithDraws, withdraw)
			}
		}
	}

	//
	layer2State, _, _ := this.layer2Sdk.GetLayer2State(chain.Height)
	msg.Layer2State = layer2State

	this.msgChan <- msg
	return nil
}

func (this *Layer2Operator) commitMsgLoop() {
	for {
		select {
		case msg := <-this.msgChan:
			for true {
				err := this.commitLayer2State2Ontology(msg)
				if err != nil {
					log.Errorf("commit layer2 state to ontology err: %s", err.Error())
				} else {
					break
				}
			}
		}
	}
}

func (this *Layer2Operator) commitLayer2State2Ontology(msg *Layer2CommitMsg) error {
	layer2Msg := msg.Dump()
	log.Infof("commit layer2 state to ontology: %s", layer2Msg)
	//
	contractAddress, _ := ontology_common.AddressFromHexString(this.config.OntologyConfig.Layer2ContractAddress)
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
		toAddresses = append(toAddresses,toAddress)
		tokenAddress, _ := hex.DecodeString(withdraw.TokenAddress)
		assetAddress = append(assetAddress, tokenAddress)
	}
	tx, err := this.ontologySdk.NeoVM.NewNeoVMInvokeTransaction(500, 40000, contractAddress, []interface{}{"updateState", []interface{}{
		msg.Layer2State.StatesRoot.ToHexString(), msg.Layer2State.Height, string(msg.Layer2State.Version),
		depositids, withdrawAmounts,toAddresses,assetAddress}})
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
		} else {
			break
		}
	}
	log.Infof("layer2 state commit transaction hash: %s", txHash.ToHexString())

	//
	for _, id := range msg.Deposits {
		UpdateDepositById(id, DEPOSIT_NOTIFY)
	}
	for _, withdraw := range msg.WithDraws {
		UpdateWithdraw(withdraw.TxHash, WITHDRAW_COMMIT, txHash.ToHexString())
	}
	SaveLayer2Commit(txHash.ToHexString(), layer2Msg, uint64(msg.Layer2State.Height))
	return nil
}

func (this *Layer2Operator) checkMsgLoop() {
	updateTicker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-updateTicker.C:
			this.checkLayer2State()
		}
	}
}

func (this *Layer2Operator) checkLayer2State() {
	txHashs := LoadLayer2Commit()
	for _, txHash := range txHashs {
		event, err := this.ontologySdk.GetSmartContractEvent(txHash)
		if err != nil {
			log.Errorf("get smart contract event failed! hash: %s, err: %s", txHash, err.Error())
			continue
		}
		heigth, err := this.ontologySdk.GetBlockHeightByTxHash(txHash)
		if err != nil {
			log.Errorf("get tx height failed! hash: %s, err: %s", txHash, err.Error())
			continue
		}
		if event == nil {
			log.Infof("layer2 commit: %s is not finished.", txHash)
			continue
		}
		for _, notify := range event.Notify {
			states := notify.States.([]interface{})
			method, _ := hex.DecodeString(states[0].(string))
			if string(method) == "updateDepositState" {

			} else if string(method) == "withdraw" {

			} else if string(method) == "updateState" {
				commitState := LAYER2MSG_COMMIT
				if event.State == 1 {
					commitState = LAYER2MSG_FINISH
				}
				UpdateLayer2Commit(event.TxHash, uint64(heigth), commitState)
				log.Infof("layer2 commit: %s is finished.", txHash)
			} else {

			}
		}

	}
}

func isLayer2Tx(addr string) bool {
	newAddr,_ := layer2_common.AddressFromBase58(addr)
	if newAddr == layer2_common.ADDRESS_EMPTY {
		return true
	} else {
		return false
	}
}

func (this *Layer2Operator) testLoop() {
	updateTicker := time.NewTicker(time.Second * 100)
	for {
		select {
		case <- updateTicker.C:
			this.test()
		}
	}
}

func (this *Layer2Operator) test() error {
	if this.deposit != 1 || this.withdraw != 0 || this.layer2ChainInfo.Height <= this.depositHeight {
		return nil
	}

	this.withdraw = 1
	ontAddr, _ := layer2_common.AddressFromHexString(ONT_CONTRACT_ADDRESS)
	ongAddr, _ := layer2_common.AddressFromHexString(ONG_CONTRACT_ADDRESS)

	txhash, err := this.transfer(this.layer2Account, ontAddr, this.layer2Account.Address, layer2_common.ADDRESS_EMPTY, 1000)
	if err != nil {
		log.Errorf("test - withdraw ont err: %s", err.Error())
	} else {
		log.Infof("test - withdraw ont from %s to %s, hash: %s", this.layer2Account.Address.ToBase58(), layer2_common.ADDRESS_EMPTY.ToBase58(), txhash.ToHexString())
	}
	txhash, err = this.transfer(this.layer2Account, ongAddr, this.layer2Account.Address, layer2_common.ADDRESS_EMPTY, 1000)
	if err != nil {
		log.Errorf("test - withdraw ong err: %s", err.Error())
	} else {
		log.Infof("test - withdraw ong from %s to %s, hash: %s", this.layer2Account.Address.ToBase58(), layer2_common.ADDRESS_EMPTY.ToBase58(), txhash.ToHexString())
	}

	txhash, err = this.transfer(this.layer2Account, ontAddr, this.layer2Account.Address, this.layer2Account.Address, 1000)
	if err != nil {
		log.Errorf("test - transfer ont err: %s", err.Error())
	} else {
		log.Infof("test - transfer ont from %s to %s, hash: %s", this.layer2Account.Address.ToBase58(), this.ontologyAccount.Address.ToBase58(), txhash.ToHexString())
	}
	txhash, err = this.transfer(this.layer2Account, ongAddr, this.layer2Account.Address, this.layer2Account.Address, 1000)
	if err != nil {
		log.Errorf("test - transfer ong err: %s", err.Error())
	} else {
		log.Infof("test - transfer ong from %s to %s, hash: %s", this.layer2Account.Address.ToBase58(), this.ontologyAccount.Address.ToBase58(), txhash.ToHexString())
	}

	return nil
}

func (this *Layer2Operator) transfer(payer *layer2_sdk.Account, token layer2_common.Address, from layer2_common.Address, to layer2_common.Address, amount uint64) (layer2_common.Uint256, error) {
	var tx *layer2_types.MutableTransaction
	var err error
	if token.ToHexString() == ONT_CONTRACT_ADDRESS {
		tx, err = this.layer2Sdk.Native.Ont.NewTransferTransaction(0, 20000, from, to, amount)
		if err != nil {
			return layer2_common.UINT256_EMPTY, err
		}
	} else if token.ToHexString() == ONG_CONTRACT_ADDRESS {
		tx, err = this.layer2Sdk.Native.Ong.NewTransferTransaction(0, 20000, from, to, amount)
		if err != nil {
			return layer2_common.UINT256_EMPTY, err
		}
	}
	if payer != nil {
		this.layer2Sdk.SetPayer(tx, payer.Address)
		err = this.layer2Sdk.SignToTransaction(tx, payer)
		if err != nil {
			return layer2_common.UINT256_EMPTY, err
		}
	}
	return this.layer2Sdk.SendTransaction(tx)
}
