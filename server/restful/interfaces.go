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