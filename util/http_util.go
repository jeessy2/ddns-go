package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const LastModifyTimeEnv = "DDNS_GO_LAST_MODIFY_TIME"

// GetHTTPResponse 处理HTTP结果，返回序列化的json
func GetHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
	body, err := GetHTTPResponseOrg(resp, url, err)

	if err == nil {
		// log.Println(string(body))
		if len(body) != 0 {
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Printf("请求接口%s解析json结果失败! ERROR: %s\n", url, err)
			}
		}
	}

	return err

}

// GetHTTPResponseOrg 处理HTTP结果，返回byte
func GetHTTPResponseOrg(resp *http.Response, url string, err error) ([]byte, error) {
	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
	}

	// 300及以上状态码都算异常
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("请求接口 %s 失败! 返回内容: %s ,返回状态码: %d\n", url, string(body), resp.StatusCode)
		log.Println(errMsg)
		err = fmt.Errorf(errMsg)
	}

	return body, err
}

// CheckStaticCache 检查静态文件缓存
func CheckStaticCache(writer http.ResponseWriter, request *http.Request) bool {
	// 获取请求头中的If-Modified-Since
	ifModifiedSince := request.Header.Get("If-Modified-Since")
	lastModifyTime, err := http.ParseTime(os.Getenv(LastModifyTimeEnv))
	if ifModifiedSince != "" && err == nil {
		// 将时间字符串转为时间戳
		ifModifiedSinceTime, err := http.ParseTime(ifModifiedSince)
		if err == nil && lastModifyTime.Unix() <= ifModifiedSinceTime.Unix() {
			// 设置状态码
			writer.WriteHeader(http.StatusNotModified)
			return true
		}
	}

	// 设置响应头
	writer.Header().Set("Last-Modified", lastModifyTime.UTC().Format(http.TimeFormat))
	writer.Header().Set("Cache-Control", "max-age=60")
	return false
}
