package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ontio/layer2/operator/log"
)

const (
	ETH_MONITOR_INTERVAL     = 3 * time.Second
	ONT_MONITOR_INTERVAL     = 3 * time.Second
	KEY_UNLOCK_TIME          = 30 * time.Second

	ETH_USEFUL_BLOCK_NUM      = 3
	ETH_PROOF_USERFUL_BLOCK   = 25
	ONT_USEFUL_BLOCK_NUM      = 1
	ETH_CHAIN_ID               = 2
	DEFAULT_CONFIG_FILE_NAME  = "./config.json"
	Version                   = "1.0"

	DEFAULT_LOG_LEVEL = log.InfoLog
)

//type ETH struct {
//	Chain             string // eth or etc
//	ChainId           uint64
//	RpcAddress        string
//	ConfirmedBlockNum uint
//	//Tokens            []*Token
//}

type ServiceConfig struct {
	OntologyConfig         *OntologyConfig
	DBConfig               *DBConfig
	Layer2Config           *Layer2Config
}

type OntologyConfig struct {
	RestURL                 string
	Layer2ContractAddress   string
	WalletFile              string
	WalletPwd               string
	GasPrice                uint64
	GasLimit                uint64
}

type Layer2Config struct {
	RestURL                 string
	WalletFile              string
	WalletPwd               string
	GasPrice                uint64
	GasLimit                uint64
}

type DBConfig struct {
	ProjectDBUrl       string
	ProjectDBUser      string
	ProjectDBPassword  string
	ProjectDBName      string
}

func ReadFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: open file %s error %s", fileName, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Errorf("ReadFile: File %s close error %s", fileName, err)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: ioutil.ReadAll %s error %s", fileName, err)
	}
	return data, nil
}

func NewServiceConfig(configFilePath string) *ServiceConfig {
	fileContent, err := ReadFile(configFilePath)
	if err != nil {
		log.Errorf("NewServiceConfig: failed, err: %s", err)
		return nil
	}
	servConfig := &ServiceConfig{}
	err = json.Unmarshal(fileContent, servConfig)
	if err != nil {
		log.Errorf("NewServiceConfig: failed, err: %s", err)
		return nil
	}

	return servConfig
}
