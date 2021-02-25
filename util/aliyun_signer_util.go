package util

import (
	"net/url"
	"strconv"
	"time"
)

// AliyunSigner AliyunSigner
func AliyunSigner(accessKeyID, accessSecret string, params *url.Values) {
	// 公共参数
	params.Set("SignatureMethod", "HMAC-SHA1")
	params.Set("SignatureNonce", strconv.FormatInt(time.Now().UnixNano(), 10))
	params.Set("AccessKeyId", accessKeyID)
	params.Set("SignatureVersion", "1.0")
	params.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	params.Set("Format", "JSON")
	params.Set("Version", "2015-01-09")
	params.Set("Signature", HmacSignToB64("HMAC-SHA1", "GET", accessSecret, *params))
}
