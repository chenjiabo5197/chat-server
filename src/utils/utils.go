package utils

import (
	"encoding/json"
	logger "github.com/shengkehua/xlog4go"
)

func Struct2String(object interface{}) string {
	result, err := json.Marshal(object)
	if err != nil {
		logger.Error("marshal err, err=%s", err.Error())
		return ""
	}
	return string(result)
}
