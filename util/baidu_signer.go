package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// https://cloud.baidu.com/doc/Reference/s/Njwvz1wot

const (
	BaiduDateFormat  = "2006-01-02T15:04:05Z"
	expirationPeriod = "1800"
)

func HmacSha256Hex(secret, message string) string {
	key := []byte(secret)

	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func BaiduCanonicalURI(r *http.Request) string {
	pattens := strings.Split(r.URL.Path, "/")
	var uri []string
	for _, v := range pattens {
		uri = append(uri, escape(v))
	}
	urlpath := strings.Join(uri, "/")
	if len(urlpath) == 0 || urlpath[len(urlpath)-1] != '/' {
		urlpath = urlpath + "/"
	}
	return urlpath[0 : len(urlpath)-1]
}

// BaiduSigner set Authorization header
func BaiduSigner(accessKeyID, accessSecret string, r *http.Request) {
	//format: bce-auth-v1/{accessKeyId}/{timestamp}/{expirationPeriodInSeconds}
	authStringPrefix := "bce-auth-v1/" + accessKeyID + "/" + time.Now().UTC().Format(BaiduDateFormat) + "/" + expirationPeriod
	baiduCanonicalURL := BaiduCanonicalURI(r)

	//format: HTTP Method + "\n" + CanonicalURI + "\n" + CanonicalQueryString + "\n" + CanonicalHeaders
	//由于仅仅调用三个POST接口且不会更改，这里CanonicalQueryString和CanonicalHeaders直接写死
	CanonicalReq := fmt.Sprintf("%s\n%s\n%s\n%s", r.Method, baiduCanonicalURL, "", "host:bcd.baidubce.com")

	signingKey := HmacSha256Hex(accessSecret, authStringPrefix)
	signature := HmacSha256Hex(signingKey, CanonicalReq)

	//format: authStringPrefix/{signedHeaders}/{signature}
	authString := authStringPrefix + "/host/" + signature
	r.Header.Set(HeaderAuthorization, authString)
}
