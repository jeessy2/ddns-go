package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// GetHTTPResponse 获得http response结果
func GetHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
	} else if resp.StatusCode != 200 {
		defer resp.Body.Close()
		log.Printf("请求接口%s失败! StatusCode: %d", url, resp.StatusCode)
		return errors.New("Response status code not equals 200")
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
		}

		err = json.Unmarshal(body, &result)

		if err != nil {
			log.Printf("请求接口%s解析json结果失败! ERROR: %s\n", url, err)
		}

	}
	return err
}
