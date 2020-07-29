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
	"fmt"
	ontology_common "github.com/ontio/ontology/common"
)

type Layer2State struct {
	Version    byte
	Height     uint32
	StatesRoot ontology_common.Uint256
	SigData [][]byte
}

type ChainInfo struct {
	Name        string
	Id          uint32
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
	Layer2State       *Layer2State
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