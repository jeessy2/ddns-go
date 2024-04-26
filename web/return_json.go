package web

import (
	"encoding/json"
	"net/http"
)

// Result Result
type Result struct {
	Code int         // 状态
	Msg  string      // 消息
	Data interface{} // 数据
}

// returnError 返回错误信息
func returnError(w http.ResponseWriter, msg string) {
	result := &Result{}

	result.Code = http.StatusInternalServerError
	result.Msg = msg

	json.NewEncoder(w).Encode(result)
}

// returnOK	返回成功信息
func returnOK(w http.ResponseWriter, msg string, data interface{}) {
	result := &Result{}

	result.Code = http.StatusOK
	result.Msg = msg
	result.Data = data

	json.NewEncoder(w).Encode(result)
}
