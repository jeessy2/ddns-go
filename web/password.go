package web

import passwordvalidator "github.com/wagslane/go-password-validator"

// validate 检查密码强度是否大于最低要求。如果不是则返回错误并说明如何加强密码。向客户端显示此错误是安全的。
func validate(password string) error {
	return passwordvalidator.Validate(password, 60)
}
