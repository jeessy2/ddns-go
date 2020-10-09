package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// GetHTTPResponse 获得http response结果
func GetHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
		}

		if resp.StatusCode != 200 && resp.StatusCode != 202 {
			log.Printf("请求接口%s失败! %s\n", url, string(body))
			err = fmt.Errorf("请求接口%s失败! %s", url, string(body))
		}

		// log.Println(string(body))
		err = json.Unmarshal(body, &result)

		if err != nil {
			log.Printf("请求接口%s解析json结果失败! ERROR: %s\n", url, err)
		}

	}
	return err
}
