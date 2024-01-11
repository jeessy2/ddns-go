// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/release.go

package update

import "github.com/jeessy2/ddns-go/v6/util/semver"

// Latest 表示当前操作系统和架构的最新 release asset。
type Latest struct {
	//Name 是 asset 的文件名
	Name string
	// URL 是 release 上传文件的 URL
	URL string
	// version 是解析后的 *Version
	Version *semver.Version
}

func newLatest(asset *Asset, ver *semver.Version) *Latest {
	latest := &Latest{
		Name:    asset.name,
		URL:     asset.url,
		Version: ver,
	}

	return latest
}
