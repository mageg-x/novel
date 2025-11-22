// src/assets/assets.go
package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:dict
var dictEmbed embed.FS

// DictTempDir 是解压后的字典根目录，供其他包使用
var DictTempDir string

func init() {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "jieba_dict")
	if err != nil {
		fmt.Print("failed to create temp dir for jieba dict: " + err.Error())
		return
	}
	DictTempDir = tmpDir

	// 递归解压 embed 的 dict/ 到临时目录
	err = fs.WalkDir(dictEmbed, "dict", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 去掉前缀 "dict/"，得到相对路径
		relPath, _ := filepath.Rel("dict", path)
		targetPath := filepath.Join(DictTempDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := dictEmbed.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, data, 0644)
	})
	if err != nil {
		fmt.Print("failed to extract jieba dict: " + err.Error())
		return
	}
	fmt.Println("jieba dict extracted to: " + DictTempDir)
}

// Cleanup 删除临时字典目录（供 main 调用）
func Cleanup() {
	if DictTempDir != "" {
		os.RemoveAll(DictTempDir)
	}
}
