package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

// GenerateToken 生成Token
func GenerateToken(username string) string {
	key := []byte(generateRandomKey())
	h := hmac.New(sha256.New, key)
	msg := fmt.Sprintf("%s%d", username, time.Now().Unix())
	h.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// generateRandomKey 生成随机密钥
func generateRandomKey() string {
	// 设置随机种子
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// 生成随机的64位整数
	randomNumber := random.Uint64()

	return fmt.Sprint(randomNumber)
}
