package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func tencentCloudHmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

// TencentCloudSigner 腾讯云签名方法 v3 https://cloud.tencent.com/document/api/1427/56189#Golang
func TencentCloudSigner(secretId string, secretKey string, r *http.Request, action string, payload string) {
	algorithm := "TC3-HMAC-SHA256"
	service := "dnspod"
	host := WriteString(service, ".tencentcloudapi.com")
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	// step 1: build canonical request string
	canonicalHeaders := WriteString("content-type:application/json\nhost:", host, "\nx-tc-action:", strings.ToLower(action), "\n")
	signedHeaders := "content-type;host;x-tc-action"
	hashedRequestPayload := sha256hex(payload)
	canonicalRequest := WriteString("POST\n/\n\n", canonicalHeaders, "\n", signedHeaders, "\n", hashedRequestPayload)

	// step 2: build string to sign
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := WriteString(date, "/", service, "/tc3_request")
	hashedCanonicalRequest := sha256hex(canonicalRequest)
	string2sign := WriteString(algorithm, "\n", timestampStr, "\n", credentialScope, "\n", hashedCanonicalRequest)

	// step 3: sign string
	secretDate := tencentCloudHmacsha256(date, WriteString("TC3", secretKey))
	secretService := tencentCloudHmacsha256(service, secretDate)
	secretSigning := tencentCloudHmacsha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(tencentCloudHmacsha256(string2sign, secretSigning)))

	// step 4: build authorization
	authorization := WriteString(algorithm, " Credential=", secretId, "/", credentialScope, ", SignedHeaders=", signedHeaders, ", Signature=", signature)

	r.Header.Add("Authorization", authorization)
	r.Header.Set("Host", host)
	r.Header.Set("X-TC-Action", action)
	r.Header.Add("X-TC-Timestamp", timestampStr)
}
