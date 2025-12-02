package util

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func NumToChinese(num uint) string {
	if num == 0 {
		return "零"
	}
	digits := []string{"", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	units := []string{"", "十", "百", "千"}
	bigUnits := []string{"", "万", "亿"}

	res := ""
	for i := 0; num > 0; i++ {
		section := num % 10000
		if section != 0 {
			secStr := ""
			for j := 0; section > 0; j++ {
				d := section % 10
				if d != 0 {
					secStr = digits[d] + units[j] + secStr
				} else if secStr != "" && !strings.HasPrefix(secStr, "零") {
					secStr = "零" + secStr
				}
				section /= 10
			}
			res = secStr + bigUnits[i] + res
		}
		num /= 10000
	}

	if strings.HasPrefix(res, "一十") {
		res = res[1:]
	}
	return res
}

// GetFirstChar 获取目录名的第一个有效字符（跳过空格、标点等）
// 优先返回第一个字母、数字或中文字符；若全为符号，则返回原字符串首字符（非空时），否则返回 "_"
func GetFirstChar(dirname string) string {
	trimmed := strings.TrimSpace(dirname)
	for _, char := range trimmed {
		// 判断是否为字母、数字，或中文字符（Unicode 范围 \u4e00-\u9fff）
		if unicode.IsLetter(char) || unicode.IsDigit(char) || (char >= '\u4e00' && char <= '\u9fff') {
			return string(char)
		}
	}

	// 如果字符串为空，返回 "_"
	if len(dirname) == 0 {
		return "_"
	}

	// 否则返回原字符串（未 trim）的第一个字符
	for _, char := range dirname {
		return string(char)
	}

	// 理论上不会走到这里，但为了安全返回 "_"
	return "_"
}

// GetDirSize 计算目录的总大小（以字节为单位），优化支持多级目录
func GetDirSize(path string) (int64, error) {
	var size int64
	// 使用filepath.WalkDir替代Walk，性能更好且内存占用更低
	err := filepath.WalkDir(path, func(path string, entry os.DirEntry, err error) error {
		// 如果访问当前路径出错，记录错误但继续遍历（跳过这个路径）
		if err != nil {
			return filepath.SkipDir
		}

		// 检查是否为符号链接，如果是则跳过（避免循环引用）
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}

		// 如果是文件，累加其大小
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return nil // 忽略无法获取信息的文件
			}
			size += info.Size()
		}
		return nil
	})
	return size, err
}
