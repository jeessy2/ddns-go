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

	"golang.org/x/crypto/scrypt"
)

const secretPrefix = "ENC:"
const secretPrefixV2 = "ENCv2:"

func deriveSecretKey(password string) []byte {
	sum := sha256.Sum256([]byte(password))
	return sum[:]
}

func deriveSecretKeyV2(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
}

func EncryptSecretWithPassword(hashedPassword, secret string) (string, error) {
	if secret == "" {
		return "", nil
	}
	if strings.HasPrefix(secret, secretPrefix) || strings.HasPrefix(secret, secretPrefixV2) {
		return secret, nil
	}
	if hashedPassword == "" {
		return secret, nil
	}

	// V2 Encryption
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key, err := deriveSecretKeyV2(hashedPassword, salt)
	if err != nil {
		return "", err
	}

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
	// salt + nonce + ciphertext
	buf := append(salt, nonce...)
	buf = append(buf, cipherText...)
	return secretPrefixV2 + base64.StdEncoding.EncodeToString(buf), nil
}

func DecryptSecretWithPassword(hashedPassword, stored string) (string, error) {
	if stored == "" {
		return "", nil
	}

	// V2 Decryption
	if strings.HasPrefix(stored, secretPrefixV2) {
		if hashedPassword == "" {
			return "", errors.New("missing password for secret decryption")
		}
		raw := strings.TrimPrefix(stored, secretPrefixV2)
		data, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return "", err
		}
		if len(data) < 16 {
			return "", errors.New("invalid secret data")
		}
		salt := data[:16]
		rest := data[16:]

		key, err := deriveSecretKeyV2(hashedPassword, salt)
		if err != nil {
			return "", err
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			return "", err
		}
		aead, err := cipher.NewGCM(block)
		if err != nil {
			return "", err
		}
		nonceSize := aead.NonceSize()
		if len(rest) < nonceSize {
			return "", errors.New("invalid secret data")
		}
		nonce := rest[:nonceSize]
		cipherText := rest[nonceSize:]
		plain, err := aead.Open(nil, nonce, cipherText, nil)
		if err != nil {
			return "", err
		}
		return string(plain), nil
	}

	// V1 Decryption (Legacy)
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

