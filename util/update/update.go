// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/update.go

package update

import (
	"io"
	"path/filepath"
)

func decompressAndUpdate(src io.Reader, assetName, cmdPath string) error {
	_, cmd := filepath.Split(cmdPath)
	asset, err := decompressCommand(src, assetName, cmd)
	if err != nil {
		return err
	}

	return apply(asset, cmdPath)
}
