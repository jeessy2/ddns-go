// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/detect.go

package update

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/jeessy2/ddns-go/v6/util/semver"
)

// detectLatest 尝试从源提供者获取版本信息。
func detectLatest(repo string) (latest *Latest, found bool, err error) {
	rel, err := getLatest(repo)
	if err != nil {
		return nil, false, err
	}

	asset, ver, found := findAsset(rel)
	if !found {
		return nil, false, nil
	}

	return newLatest(asset, ver), true, nil
}

// findAsset 返回最新的 asset
func findAsset(rel *Release) (*Asset, *semver.Version, bool) {
	// 将检测到的架构放在列表的末尾，对于 ARM 来说这是可以的。
	// 因为附加的架构比通用架构更准确
	for _, arch := range append(generateAdditionalArch(), runtime.GOARCH) {
		asset, version, found := findAssetForArch(arch, rel)
		if found {
			return asset, version, found
		}
	}

	return nil, nil, false
}

func findAssetForArch(arch string, rel *Release,
) (asset *Asset, version *semver.Version, found bool) {
	var release *Release

	// 从 release 列表中查找最新的版本。
	// GitHub API 返回的列表按照创建日期的顺序排列。
	if a, v, ok := findAssetFromRelease(rel, getSuffixes(arch)); ok {
		version = v
		asset = a
		release = rel
	}

	if release == nil {
		log.Printf("Cannot find any release for %s/%s", runtime.GOOS, runtime.GOARCH)
		return nil, nil, false
	}

	return asset, version, true
}

func findAssetFromRelease(rel *Release, suffixes []string) (*Asset, *semver.Version, bool) {
	if rel == nil {
		log.Print("There is no source release information")
		return nil, nil, false
	}

	// 如果无法解析版本文本，则表示该文本不符合语义化版本规范，应该跳过。
	ver, err := semver.NewVersion(rel.tagName)
	if err != nil {
		log.Printf("Cannot parse semantic version: %s", rel.tagName)
		return nil, nil, false
	}

	for _, asset := range rel.assets {
		if assetMatchSuffixes(asset.name, suffixes) {
			return &asset, ver, true
		}
	}

	log.Printf("Can't find suitable asset in release %s", rel.tagName)
	return nil, nil, false
}

func assetMatchSuffixes(name string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) { // 需要版本、架构等
			// 假设唯一的构件被匹配（或者第一个匹配将足够）
			return true
		}
	}
	return false
}

// getSuffixes 返回所有要与 asset 进行检查的候选后缀
//
// TODO: 由于缺失获取 MIPS 架构 float 的方法，所以目前无法正确获取 MIPS 架构的后缀。
func getSuffixes(arch string) []string {
	suffixes := make([]string, 0)
	for _, ext := range []string{".zip", ".tar.gz"} {
		suffix := fmt.Sprintf("%s_%s%s", runtime.GOOS, arch, ext)
		suffixes = append(suffixes, suffix)
	}
	return suffixes
}
