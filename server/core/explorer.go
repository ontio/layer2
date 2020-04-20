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
