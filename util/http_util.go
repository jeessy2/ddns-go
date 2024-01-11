package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetHTTPResponse 处理HTTP结果，返回序列化的json
func GetHTTPResponse(resp *http.Response, err error, result interface{}) error {
	body, err := GetHTTPResponseOrg(resp, err)

	if err == nil {
		// log.Println(string(body))
		if len(body) != 0 {
			err = json.Unmarshal(body, &result)
		}
	}

	return err

}

// GetHTTPResponseOrg 处理HTTP结果，返回byte
func GetHTTPResponseOrg(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	lr := io.LimitReader(resp.Body, 1024000)
	body, err := io.ReadAll(lr)

	if err != nil {
		return nil, err
	}

	// 300及以上状态码都算异常
	if resp.StatusCode >= 300 {
		err = fmt.Errorf(LogStr("返回内容: %s ,返回状态码: %d", string(body), resp.StatusCode))
	}

	return body, err
}
