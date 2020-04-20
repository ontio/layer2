package restful

import (
	"github.com/ontio/layer2/server/core"
)

//start
func GetLayer2Tx(cmd map[string]interface{}) map[string]interface{} {
	if cmd["address"] == nil {
		return ResponsePack(core.REST_PARAM_INVALID)
	}
	address, ok := cmd["address"].(string)
	if !ok {
		return ResponsePack(core.REST_PARAM_INVALID)
	}
	code, result := core.Explorer.GetLayer2Tx(address)
	if code != core.SUCCESS {
		return ResponsePack(code)
	}
	resp := ResponsePack(core.SUCCESS)
	resp["result"] = result
	return resp
}

func GetLayer2Deposit(cmd map[string]interface{}) map[string]interface{} {
	if cmd["address"] == nil {
		return ResponsePack(core.REST_PARAM_INVALID)
	}
	address, ok := cmd["address"].(string)
	if !ok {
		return ResponsePack(core.REST_PARAM_INVALID)
	}
	code, result := core.Explorer.GetLayer2Deposit(address)
	if code != core.SUCCESS {
		return ResponsePack(code)
	}
	resp := ResponsePack(core.SUCCESS)
	resp["result"] = result
	return resp
}

func GetLayer2Withdraw(cmd map[string]interface{}) map[string]interface{} {
	if cmd["address"] == nil {
		return ResponsePack(core.REST_PARAM_INVALID)
	}
	address, ok := cmd["address"].(string)
	if !ok {
		return ResponsePack(core.REST_PARAM_INVALID)
	}
	code, result := core.Explorer.GetLayer2Withdraw(address)
	if code != core.SUCCESS {
		return ResponsePack(code)
	}
	resp := ResponsePack(core.SUCCESS)
	resp["result"] = result
	return resp
}