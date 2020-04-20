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
	"encoding/json"
	"fmt"
	"github.com/ontio/layer2/server/config"
	"github.com/ontio/layer2/server/log"
)

var Explorer *explorer

type explorer struct {

}

func NewExplorer() *explorer {
	Explorer = &explorer{
	}
	return Explorer
}

func (self *explorer) Start() (code int64, err error) {
	// try to connect db
	dberr := ConnectDB(config.DefConfig.ProjectDBUser, config.DefConfig.ProjectDBPassword, config.DefConfig.ProjectDBUrl, config.DefConfig.ProjectDBName)
	if dberr != nil {
		return DB_CONNECTTION_FAILED, fmt.Errorf(dberr.Error())
	}
	//
	return SUCCESS, nil
}

func (self *explorer) Stop() error {
	// try to close db connection
	CloseDB()

	return nil
}

func (self *explorer) GetLayer2Tx(address string) (int64,string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("recover info:", r)
		}
	}()

	layer2Tx := LoadLayer2Tx(address)
	newLayer2Tx := make([]*Layer2Tx, 0)
	for _, tx := range layer2Tx {
		if tx.FromAddress == "AFmseVrdL9f9oyCzZefL9tG6UbvhPbdYzM" {
			continue
		}
		if tx.ToAddress == "AFmseVrdL9f9oyCzZefL9tG6UbvhPbdYzM" {
			continue
		}
		newLayer2Tx = append(newLayer2Tx, tx)
	}
	json_crosstx, _ := json.Marshal(newLayer2Tx)
	return SUCCESS, string(json_crosstx)
}

func (self *explorer) GetLayer2Deposit(address string) (int64,string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("recover info:", r)
		}
	}()

	layer2Tx := LoadLayer2Tx(address)
	newLayer2Tx := make([]*Layer2Tx, 0)
	for _, tx := range layer2Tx {
		if tx.FromAddress == "AFmseVrdL9f9oyCzZefL9tG6UbvhPbdYzM" {
			newLayer2Tx = append(newLayer2Tx, tx)
		}
	}
	json_crosstx, _ := json.Marshal(newLayer2Tx)
	return SUCCESS, string(json_crosstx)
}

func (self *explorer) GetLayer2Withdraw(address string) (int64,string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("recover info:", r)
		}
	}()

	layer2Tx := LoadLayer2Tx(address)
	newLayer2Tx := make([]*Layer2Tx, 0)
	for _, tx := range layer2Tx {
		if tx.ToAddress == "AFmseVrdL9f9oyCzZefL9tG6UbvhPbdYzM" {
			newLayer2Tx = append(newLayer2Tx, tx)
		}
	}
	json_crosstx, _ := json.Marshal(newLayer2Tx)
	return SUCCESS, string(json_crosstx)
}
