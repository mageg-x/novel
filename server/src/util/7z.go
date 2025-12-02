// 7z.go
package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
)

// Extract7z 解压 .7z 文件到目标目录
// archivePath: .7z 文件路径
// destDir:     解压目标目录（会自动创建）
// password:    可选密码，无密码传空字符串 ""
func Extract7z(archivePath, destDir, password string) error {
	// 打开 7z 文件
	var r *sevenzip.ReadCloser
	var err error

	if password != "" {
		r, err = sevenzip.OpenReaderWithPassword(archivePath, password)
	} else {
		r, err = sevenzip.OpenReader(archivePath)
	}
	if err != nil {
		return fmt.Errorf("failed to open 7z file: %w", err)
	}
	defer r.Close()

	// 确保目标目录存在
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 遍历所有文件/目录
	for _, file := range r.File {
		targetPath := filepath.Join(destDir, filepath.FromSlash(file.Name))

		// 防止路径穿越（安全加固）
		if !isSubPath(destDir, targetPath) {
			fmt.Printf("skipping invalid path: %s\n", file.Name)
			continue
		}

		if file.FileInfo().IsDir() {
			// 创建目录
			if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
				fmt.Printf("failed to create directory %s: %v\n", targetPath, err)
			}
		} else {
			// 创建父目录
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				fmt.Printf("failed to create parent directory %s: %v\n", filepath.Dir(targetPath), err)
				continue
			}

			// 打开文件流
			rc, err := file.Open()
			if err != nil {
				fmt.Printf("failed to open file %s: %v\n", file.Name, err)
				continue
			}

			// 创建目标文件
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
			if err != nil {
				rc.Close()
				fmt.Printf("failed to create target file %s: %v\n", targetPath, err)
				continue
			}

			// 拷贝内容
			_, err = io.Copy(outFile, rc)
			rc.Close()
			outFile.Close()

			if err != nil {
				fmt.Printf("failed to write file %s: %v\n", targetPath, err)
			}
		}
	}

	return nil
}

// isSubPath 确保 target 在 base 目录内（防止 ../ 路径穿越）
func isSubPath(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

// // ===== Example usage =====
// func main() {
// 	archive := "/tmp/novel.7z"
// 	dest := "./extracted_novel"
// 	password := "" // enter password here if needed

// 	fmt.Printf("extracting %s to %s...\n", archive, dest)
// 	if err := Extract7z(archive, dest, password); err != nil {
// 		fmt.Printf("extraction failed: %v\n", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println("extraction completed!")
// }
