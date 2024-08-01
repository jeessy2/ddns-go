package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 密码哈希
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// PasswordOK 检查密码
func PasswordOK(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// IsHashedPassword 是否是哈希密码
func IsHashedPassword(password string) bool {
	_, err := bcrypt.Cost([]byte(password))
	return err == nil
}
