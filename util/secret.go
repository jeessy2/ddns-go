package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

const secretPrefix = "ENC:"

func deriveSecretKey(password string) []byte {
	sum := sha256.Sum256([]byte(password))
	return sum[:]
}

func EncryptSecretWithPassword(hashedPassword, secret string) (string, error) {
	if secret == "" {
		return "", nil
	}
	if strings.HasPrefix(secret, secretPrefix) {
		return secret, nil
	}
	if hashedPassword == "" {
		return secret, nil
	}
	key := deriveSecretKey(hashedPassword)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := aead.Seal(nil, nonce, []byte(secret), nil)
	buf := append(nonce, cipherText...)
	return secretPrefix + base64.StdEncoding.EncodeToString(buf), nil
}

func DecryptSecretWithPassword(hashedPassword, stored string) (string, error) {
	if stored == "" {
		return "", nil
	}
	if !strings.HasPrefix(stored, secretPrefix) {
		return stored, nil
	}
	if hashedPassword == "" {
		return "", errors.New("missing password for secret decryption")
	}
	raw := strings.TrimPrefix(stored, secretPrefix)
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", err
	}
	key := deriveSecretKey(hashedPassword)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aead.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid secret data")
	}
	nonce := data[:nonceSize]
	cipherText := data[nonceSize:]
	plain, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

