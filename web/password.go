package web

import (
	"errors"
	"fmt"
	"strings"

	passwordvalidator "github.com/wagslane/go-password-validator"
)

const (
	replaceChars      = `!@$&*`
	sepChars          = `_-., `
	otherSpecialChars = `"#%'()+/:;<=>?[\]^{|}~`
	lowerChars        = `abcdefghijklmnopqrstuvwxyz`
	upperChars        = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	digitsChars       = `0123456789`
)

// validate 检查密码强度是否大于最低要求（50）。如果不是则返回错误并说明如何加强密码。向客户端显示此错误是安全的。
func validate(password string) error {
	return validatePassword(password, 50)
}

// validatePassword 在密码大于或等于 minEntropy 时返回 nil。如果不是则返回错误。
// 这解释了如何加强密码。向客户端显示此错误是安全的。
//
// https://github.com/wagslane/go-password-validator/blob/v0.3.0/validate.go#L13
func validatePassword(password string, minEntropy float64) error {
	entropy := passwordvalidator.GetEntropy(password)
	if entropy >= minEntropy {
		return nil
	}

	hasReplace := false
	hasSep := false
	hasOtherSpecial := false
	hasLower := false
	hasUpper := false
	hasDigits := false
	for _, c := range password {
		if strings.ContainsRune(replaceChars, c) {
			hasReplace = true
			continue
		}
		if strings.ContainsRune(sepChars, c) {
			hasSep = true
			continue
		}
		if strings.ContainsRune(otherSpecialChars, c) {
			hasOtherSpecial = true
			continue
		}
		if strings.ContainsRune(lowerChars, c) {
			hasLower = true
			continue
		}
		if strings.ContainsRune(upperChars, c) {
			hasUpper = true
			continue
		}
		if strings.ContainsRune(digitsChars, c) {
			hasDigits = true
			continue
		}
	}

	allMessages := []string{}

	if !hasOtherSpecial || !hasSep || !hasReplace {
		allMessages = append(allMessages, "包含更多特殊字符")
	}
	if !hasLower {
		allMessages = append(allMessages, "使用小写字母")
	}
	if !hasUpper {
		allMessages = append(allMessages, "使用大写字母")
	}
	if !hasDigits {
		allMessages = append(allMessages, "使用数字")
	}

	if len(allMessages) > 0 {
		return fmt.Errorf(
			"密码不安全！尝试%v或使用更长的密码",
			strings.Join(allMessages, "，"),
		)
	}

	return errors.New("密码不安全！尝试使用更长的密码")
}
