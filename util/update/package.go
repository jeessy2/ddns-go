package update

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/jeessy2/ddns-go/v6/util"
	"github.com/jeessy2/ddns-go/v6/util/semver"
)

// Self 更新 ddns-go 到最新版本（如果可用）。
func Self(version string) {
	// 如果不为语义化版本立即退出
	v, err := semver.NewVersion(version)
	if err != nil {
		log.Printf("Cannot update because: %v", err)
		return
	}

	latest, found, err := detectLatest("jeessy2/ddns-go")
	if err != nil {
		log.Printf("Error happened when detecting latest version: %v", err)
		return
	}
	if !found {
		log.Printf("Cannot find any release for %s/%s", runtime.GOOS, runtime.GOARCH)
		return
	}
	if v.GreaterThanOrEqual(latest.Version) {
		log.Printf("Current version (%s) is the latest", version)
		return
	}

	exe, err := os.Executable()
	if err != nil {
		log.Printf("Cannot find executable path: %v", err)
		return
	}

	if err = to(latest.URL, latest.Name, exe); err != nil {
		log.Printf("Error happened when updating binary: %v", err)
		return
	}

	log.Printf("Success update to v%s", latest.Version.String())
}

// to 从 assetURL 下载可执行文件，并用下载的文件替换当前的可执行文件。
// 这个函数是用于更新二进制文件的低级 API。因为它不使用源提供者，而是直接通过 HTTP 从 URL 下载 asset 。
// 所以这个函数不能用于更新私有仓库的 release。
// cmdPath 是命令可执行文件的文件路径。
func to(assetURL, assetFileName, cmdPath string) error {
	src, err := downloadAssetFromURL(assetURL)
	if err != nil {
		return err
	}
	defer src.Close()
	return decompressAndUpdate(src, assetFileName, assetURL, cmdPath)
}

func downloadAssetFromURL(url string) (rc io.ReadCloser, err error) {
	client := util.CreateHTTPClient()
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not download release from %s: %v", url, err)
	}
	if resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("could not download release from %s. Response code: %d", url, resp.StatusCode)
	}

	return resp.Body, nil
}
