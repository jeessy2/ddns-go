// Based on https://github.com/Masterminds/semver/blob/v3.2.1/version.go

package semver

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 在 init() 中创建的正则表达式的编译版本被缓存在这里，这样
// 它只需要被创建一次。
var versionRegex *regexp.Regexp

// semVerRegex 是用于解析语义化版本的正则表达式。
const semVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

// Version 表示单独的语义化版本。
type Version struct {
	major, minor, patch uint64
}

func init() {
	versionRegex = regexp.MustCompile("^" + semVerRegex + "$")
}

// NewVersion 解析给定的版本并返回 Version 实例，如果
// 无法解析该版本则返回错误。如果版本是类似于 SemVer 的版本，则会
// 尝试将其转换为 SemVer。
func NewVersion(v string) (*Version, error) {
	m := versionRegex.FindStringSubmatch(v)
	if m == nil {
		return nil, fmt.Errorf("the %s, it's not a semantic version", v)
	}

	sv := &Version{}

	var err error
	sv.major, err = strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("解析版本号时出错：%s", err)
	}

	if m[2] != "" {
		sv.minor, err = strconv.ParseUint(strings.TrimPrefix(m[2], "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("解析版本号时出错：%s", err)
		}
	} else {
		sv.minor = 0
	}

	if m[3] != "" {
		sv.patch, err = strconv.ParseUint(strings.TrimPrefix(m[3], "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("解析版本号时出错：%s", err)
		}
	} else {
		sv.patch = 0
	}

	return sv, nil
}

// String 将 Version 对象转换为字符串。
// 注意，如果原始版本包含前缀 v，则转换后的版本将不包含 v。
// 根据规范，语义版本不包含前缀 v，而在实现上则是可选的。
func (v Version) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%d.%d.%d", v.major, v.minor, v.patch)

	return buf.String()
}

// GreaterThan 测试一个版本是否大于另一个版本。
func (v *Version) GreaterThan(o *Version) bool {
	return v.compare(o) > 0
}

// GreaterThanOrEqual 测试一个版本是否大于或等于另一个版本。
func (v *Version) GreaterThanOrEqual(o *Version) bool {
	return v.compare(o) >= 0
}

// compare 比较当前版本与另一个版本。如果当前版本小于另一个版本则返回 -1；如果两个版本相等则返回 0；如果当前版本大于另一个版本，则返回 1。
//
// 版本比较是基于 X.Y.Z 格式进行的。
func (v *Version) compare(o *Version) int {
	// 比较主版本号、次版本号和修订号。如果
	// 发现差异则返回比较结果。
	if d := compareSegment(v.major, o.major); d != 0 {
		return d
	}
	if d := compareSegment(v.minor, o.minor); d != 0 {
		return d
	}
	if d := compareSegment(v.patch, o.patch); d != 0 {
		return d
	}

	return 0
}

func compareSegment(v, o uint64) int {
	if v < o {
		return -1
	}
	if v > o {
		return 1
	}

	return 0
}
