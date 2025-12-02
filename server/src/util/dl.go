package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode"

	py "github.com/mozillazg/go-pinyin"
)

const owner, repo = "mageg-x", "novel-library"

// PinyinName 将小说名转为驼峰式拼音文件名，如 "一个人修仙" → "YiGeRenXiuXian.7z"
func PinyinName(name string) string {
	var parts []string
	args := py.NewArgs()
	args.Style = py.Style // 无声调，如 "yi"

	for _, r := range name {
		if r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
			// ASCII 字符直接保留（保持原样，如 DOTA）
			parts = append(parts, string(r))
		} else {
			// 转单个汉字为拼音
			pyList := py.Pinyin(string(r), args)
			if len(pyList) > 0 && len(pyList[0]) > 0 {
				p := pyList[0][0] // 获取第一个拼音（最常用读音）
				// 首字母大写，其余小写（模拟 Python pypinyin 的 Style.NORMAL + capitalize）
				if p != "" {
					capitalized := strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
					parts = append(parts, capitalized)
				} else {
					parts = append(parts, "_")
				}
			} else {
				parts = append(parts, "_")
			}
		}
	}

	result := strings.Join(parts, "")
	// 清理非法文件名字符（只保留字母、数字、下划线、点、连字符）
	cleaned := make([]rune, 0, len(result))
	for _, r := range result {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '-' {
			cleaned = append(cleaned, r)
		} else {
			cleaned = append(cleaned, '_')
		}
	}
	return string(cleaned) + ".7z"
}

// GetUrl 返回 GitHub 下载 URL
// name: 完整小说名（不含 .7z），如 "一个人修仙"
func GetUrl(name string) string {
	if name == "" {
		return ""
	}
	runes := []rune(name)
	if len(runes) == 0 {
		return ""
	}
	tag := string(runes[0])  // 首字作为 tag
	file := PinyinName(name) // 整个名字转拼音
	return fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		owner, repo,
		url.PathEscape(tag),
		file,
	)
}

// Download 下载小说文件到 dstPath
// name: 完整小说名（不含 .7z），如 "一个人修仙"
func Download(name, dstPath string) error {
	url := GetUrl(name)
	if url == "" {
		return fmt.Errorf("invalid name")
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
