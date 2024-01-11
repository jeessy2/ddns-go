// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/github_release.go
// and https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/github_source.go

package update

import (
	"fmt"

	"github.com/jeessy2/ddns-go/v6/util"
)

type Release struct {
	tagName string
	assets  []Asset
}

type Asset struct {
	name string
	url  string
}

// ReleaseResp 表示仓库中的 GitHub release 和 asset。
type ReleaseResp struct {
	TagName string `json:"tag_name,omitempty"`
	Assets  []struct {
		Name               string `json:"name,omitempty"`
		BrowserDownloadURL string `json:"browser_download_url,omitempty"`
	} `json:"assets,omitempty"`
}

// getLatest 列出仓库的最新 release 并返回包装过的 Release
//
// GitHub API 文档：https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
func getLatest(repo string) (*Release, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	client := util.CreateHTTPClient()
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}

	var result ReleaseResp
	err = util.GetHTTPResponse(resp, err, &result)
	if err != nil {
		util.Log("异常信息: %s", err)
		return nil, err
	}

	return newRelease(&result), err
}

func newRelease(from *ReleaseResp) *Release {
	release := &Release{
		tagName: from.TagName,
		assets:  make([]Asset, len(from.Assets)),
	}
	for i, fromAsset := range from.Assets {
		release.assets[i] = Asset{
			name: fromAsset.Name,
			url:  fromAsset.BrowserDownloadURL,
		}
	}
	return release
}
