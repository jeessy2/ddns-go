// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/errors.go

package update

import "errors"

var (
	errCannotDecompressFile        = errors.New("无法解压")
	errExecutableNotFoundInArchive = errors.New("找不到可执行文件")
)
