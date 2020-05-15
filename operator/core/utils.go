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
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/ontio/layer2/go-sdk/common"
)

const (
	NOTIFY_TRANSFER = "transfer"

	ONT_CONTRACT_ADDRESS               = "0000000000000000000000000000000000000001"
	ONT_REV_CONTRACT_ADDRESS               = "0100000000000000000000000000000000000000"
	ONT_CONTRACT_ADDRESS_BASE58        = "AFmseVrdL9f9oyCzZefL9tG6UbvhUMqNMV"
	ONG_CONTRACT_ADDRESS               = "0000000000000000000000000000000000000002"
	ONG_REV_CONTRACT_ADDRESS               = "0200000000000000000000000000000000000000"
	ONG_CONTRACT_ADDRESS_BASE58        = "AFmseVrdL9f9oyCzZefL9tG6UbvhfRZMHJ"
	GOVERNANCE_CONTRACT_ADDRESS        = "0700000000000000000000000000000000000000"
	GOVERNANCE_CONTRACT_ADDRESS_BASE58 = "AFmseVrdL9f9oyCzZefL9tG6UbviEH9ugK"
)

const (
	DEPOSIT_EVENT = iota
	DEPOSIT_COMMIT
	DEPOSIT_FINISH
	DEPOSIT_NOTIFY
)

const (
	WITHDRAW_INIT = iota
	WITHDRAW_COMMIT
)

const (
	LAYER2MSG_COMMIT = iota
	LAYER2MSG_FINISH
	LAYER2MSG_FAILED
)

type ChainInfo struct {
	Name        string
	Id          uint32
	Url         string
	Height      uint32
}

type Deposit struct {
	TxHash          string
	TT              uint32
	State           int
	Height          uint32
	FromAddress     string
	Amount          uint64
	TokenAddress    string
	ID              uint64
	Layer2TxHash    string
}

func (this *Deposit) Dump() string {
	dumpStr := ""
	dumpStr += fmt.Sprintf("Deposit: TxHash: %s, TT: %d, State: %d, Height: %d, FromAddress: %s, Amount: %d, TokenAddress: %s",
		this.TxHash, this.TT, this.State, this.Height, this.FromAddress, this.Amount, this.TokenAddress)
	return dumpStr
}

type Withdraw struct {
	TxHash          string
	TT              uint32
	State           int
	Height          uint32
	ToAddress       string
	Amount          uint64
	TokenAddress    string
	OntologyTxHash  string
}

func (this *Withdraw) Dump() string {
	dumpStr := ""
	dumpStr += fmt.Sprintf("Withdraw: TxHash: %s, TT: %d, State: %d, Height: %d, ToAddress: %s, Amount: %d, TokenAddress: %s",
		this.TxHash, this.TT, this.State, this.Height, this.ToAddress, this.Amount, this.TokenAddress)
	return dumpStr
}

type Layer2Tx struct {
	TxHash           string
	State            int
	TT               uint32
	Fee              uint64
	Height           uint32
	FromAddress      string
	TokenAddress     string
	ToAddress        string
	Amount           uint64
}

func (this *Layer2Tx) Dump() string {
	dumpStr := ""
	dumpStr += fmt.Sprintf("Layer2Tx: TxHash: %s, TT: %d, State: %d, Fee: %d, Height: %d, FromAddress: %s, ToAddress: %s, Amount: %d, TokenAddress: %s",
		this.TxHash, this.TT, this.State, this.Fee, this.Height, this.FromAddress, this.ToAddress, this.Amount, this.TokenAddress)
	return dumpStr
}

type Layer2CommitMsg struct {
	Layer2State       *common.Layer2State
	Deposits          []uint64
	WithDraws         []*Withdraw
}

func (this *Layer2CommitMsg) Dump() string {
	dumpStr := "Layer2 commit msg: \n"
	dumpStr += fmt.Sprintf("layer2 state, Version: %d, Height: %d, StatesRoot: %s\n",
		this.Layer2State.Version, this.Layer2State.Height, this.Layer2State.StatesRoot.ToHexString())
	dumpStr += "deposits, ["
	for _, deposit := range this.Deposits {
		dumpStr += fmt.Sprintf(" %d ", deposit)
	}
	dumpStr += "]\n"
	for _, withdraw := range this.WithDraws {
		dumpStr += withdraw.Dump()
		dumpStr += "\n"
	}
	return dumpStr
}

func (this *Layer2CommitMsg) Dump1() string {
	dumpStr := "Layer2 commit msg: "
	dumpStr += fmt.Sprintf("layer2 state, Version: %d, Height: %d, StatesRoot: %s",
		this.Layer2State.Version, this.Layer2State.Height, this.Layer2State.StatesRoot.ToHexString())
	return dumpStr
}

func revertHexString(a string) string {
	b, _ := hex.DecodeString(a)
	c := make([]byte, 0)
	for i := len(b) - 1;i >= 0;i -- {
		c = append(c, b[i])
	}
	return hex.EncodeToString(c)
}

func revertHex(a []byte) []byte {
	c := make([]byte, 0)
	for i := len(a) - 1;i >= 0;i -- {
		c = append(c, a[i])
	}
	return c
}

func BytesToInt(bys []byte) uint64 {
	for i := len(bys);i <= 8;i ++ {
		bys = append(bys, 0)
	}
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.LittleEndian, &data)
	return uint64(data)
}