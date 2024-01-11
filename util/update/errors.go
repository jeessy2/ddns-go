// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/errors.go

package update

import "errors"

var (
	errCannotDecompressFile        = errors.New("failed to decompress")
	errExecutableNotFoundInArchive = errors.New("executable not found")
)
