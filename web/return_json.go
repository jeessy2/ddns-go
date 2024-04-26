package web

import (
	"encoding/json"
	"net/http"
)

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
