package restful

import "github.com/ontio/layer2/server/core"

func ResponsePack(code int64) map[string]interface{} {
	resp := map[string]interface{}{
		"action":  "",
		"result":  "",
		"code":   code,
		"desc":    core.ErrMap[code],
		"version": "1.0.0",
	}
	return resp
}
