package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"net/url"
	"strconv"
)

// MakeHmacSha256Key sha256 签名，并返回 k1=v1&k2=v2 格式的值
// mode: 1.钉钉模式 2.飞书模式
func MakeHmacSha256Key(timestamp int64, secret string, mode int) string {
	params := url.Values{}
	if secret != "" {
		sign := HmacSha256(timestamp, secret, mode)
		params.Add("timestamp", strconv.FormatInt(timestamp, 10))
		params.Add("sign", sign)
	}
	return params.Encode()
}

// HmacSha256 sha256 签名
// mode: 1.钉钉模式 2.飞书模式
func HmacSha256(timestamp int64, secret string, mode int) string {
	var h hash.Hash
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	if mode == 1 {
		h = hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(stringToSign))
	} else {
		h = hmac.New(sha256.New, []byte(stringToSign))
		h.Write(nil)
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
