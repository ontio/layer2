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
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/layer2/server/log"
)

var DefDB *sql.DB

func ConnectDB(dbuser string, dbpassword string, dburl string, dbname string) error {
	db, dberr := sql.Open("mysql",
		dbuser+
			":"+dbpassword+
			"@tcp("+dburl+
			")/"+dbname+
			"?charset=utf8")
	if dberr != nil {
		return dberr
	}
	err := db.Ping()
	if err != nil {
		return err
	}
	DefDB = db
	return nil
}

func CloseDB() {
	if DefDB != nil {
		DefDB.Close()
	}
}


func LoadChainInfo(name string) *ChainInfo {
	strsql := "select id,height from chain_info where name = ?"
	stmt, err := DefDB.Prepare(strsql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil
	}
	rows, err := stmt.Query(name)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil
	}

	var height,id uint32
	var chain *ChainInfo
	for rows.Next() {
		if err = rows.Scan(&id, &height); err != nil {
			return nil
		} else {
			chain = &ChainInfo{
				Id : id,
				Name : name,
				Height: height,
			}
			break
		}
	}
	return chain
}

func SetChainParseHeight(id uint32, height uint32) error {
	strSql := "update chain_info set height = ? where id = ?"
	stmt, dberr := DefDB.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if dberr != nil {
		return dberr
	}
	_, dberr = stmt.Exec(height, id)
	return dberr
}

func SaveDeposit(deposit *Deposit) error {
	strSql := "insert into deposit(txhash, tt, state, height, fromaddress, amount, tokenaddress, id) values (?,?,?,?,?,?,?,?)"
	stmt, dberr := DefDB.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if dberr != nil {
		return dberr
	}
	_, dberr = stmt.Exec(deposit.TxHash, deposit.TT, deposit.State,deposit.Height, deposit.FromAddress, deposit.Amount, deposit.TokenAddress, deposit.ID)
	return dberr
}

func UpdateDepositByID(id uint64, state int, layer2TxHash string) error {
	strSql := "update deposit set layer2txhash = ?, state = ? where id = ?"
	stmt, dberr := DefDB.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if dberr != nil {
		return dberr
	}
	_, dberr = stmt.Exec(layer2TxHash, state, id)
	return dberr
}

func UpdateDepositByID2(id uint64, state int) error {
	strSql := "update deposit set state = ? where id = ?"
	stmt, dberr := DefDB.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if dberr != nil {
		return dberr
	}
	_, dberr = stmt.Exec(state, id)
	return dberr
}

func LoadDepositByLayer2TxHash(layer2TxHash string) *Deposit {
	strsql := "select txhash,tt,state,height,fromaddress,amount,tokenaddress,id,layer2txhash from deposit where layer2txhash = ?"
	stmt, err := DefDB.Prepare(strsql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil
	}
	rows, err := stmt.Query(layer2TxHash)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil
	}

	var height,tt uint32
	var state int
	var txhash, fromaddress,tokenaddress string
	var amount,id uint64
	var deposit *Deposit
	for rows.Next() {
		if err = rows.Scan(&txhash, &tt, &state, &height, &fromaddress, &amount, &tokenaddress, &id, &layer2TxHash); err != nil {
			return nil
		} else {
			deposit = &Deposit{
				TxHash : txhash,
				TT: tt,
				State: state,
				Height: height,
				FromAddress: fromaddress,
				TokenAddress: tokenaddress,
				ID: id,
				Layer2TxHash: layer2TxHash,
			}
			break
		}
	}
	return deposit
}

func SaveWithdraw(withdraw *Withdraw) error {
	strSql := "insert into withdraw(txhash, tt, state, height, toaddress, amount, tokenaddress) values (?,?,?,?,?,?,?)"
	stmt, dberr := DefDB.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	log.Infof("insert into withdraw: %s", withdraw.Dump())
	if dberr != nil {
		return dberr
	}
	_, dberr = stmt.Exec(withdraw.TxHash, withdraw.TT, withdraw.State,withdraw.Height, withdraw.ToAddress, withdraw.Amount, withdraw.TokenAddress)
	return dberr
}

func UpdateWithdraw(txHash string, state int, ontologyTxHash string) error {
	strSql := "update withdraw set ontologytxhash = ?, state = ? where txhash = ?"
	stmt, dberr := DefDB.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if dberr != nil {
		return dberr
	}
	_, dberr = stmt.Exec(ontologyTxHash, state, txHash)
	return dberr
}


func SaveLayer2Tx(layer2Tx *Layer2Tx) error {
	strSql := "insert into layer2tx(txhash, tt, state, fee, height, fromaddress, tokenaddress, toaddress, amount) values (?,?,?,?,?,?,?,?,?)"
	stmt, dberr := DefDB.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if dberr != nil {
		return dberr
	}
	_, dberr = stmt.Exec(layer2Tx.TxHash, layer2Tx.TT, layer2Tx.State,layer2Tx.Fee,layer2Tx.Height, layer2Tx.FromAddress,layer2Tx.TokenAddress,layer2Tx.ToAddress, layer2Tx.Amount)
	return dberr
}

func LoadLayer2Tx(address string) []*Layer2Tx {
	strsql := "select txhash, state, tt, fee, height, fromaddress, tokenaddress, toaddress, amount from layer2tx where fromaddress = ? or toaddress = ? order by height"
	stmt, err := DefDB.Prepare(strsql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil
	}
	rows, err := stmt.Query(address, address)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil
	}

	var tt, height uint32
	var state int
	var fee, amount uint64
	var txhash, fromaddress, tokenaddress,toaddress string
	layer2Txs := make([]*Layer2Tx, 0)
	for rows.Next() {
		if err = rows.Scan(&txhash, &state, &tt, &fee, &height, &fromaddress, &tokenaddress, &toaddress, &amount); err != nil {
			return nil
		} else {
			layer2Txs = append(layer2Txs, &Layer2Tx{
				TxHash: txhash,
				State: state,
				TT: tt,
				Fee: fee,
				Height: height,
				FromAddress: fromaddress,
				TokenAddress: tokenaddress,
				ToAddress: toaddress,
				Amount: amount,
			})
		}
	}
	return layer2Txs
}

func SaveLayer2Commit(txHash string, layer2Msg string, layer2Height uint64) error {
	strSql := "insert into layer2commit(txhash, layer2msg, layer2height) values (?,?,?)"
	stmt, dberr := DefDB.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if dberr != nil {
		return dberr
	}
	_, dberr = stmt.Exec(txHash, layer2Msg, layer2Height)
	return dberr
}

