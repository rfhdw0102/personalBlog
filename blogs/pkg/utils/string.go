package utils

import (
	"errors"
	"regexp"
	"strings"
)

// TrimSpace 去除首尾空格
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return TrimSpace(s) == ""
}

// ExtractImages 提取内容中的图片地址
func ExtractImages(content string) []string {
	// 匹配 Markdown 图片格式 ![alt](url)
	mdRegex := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	mdMatches := mdRegex.FindAllStringSubmatch(content, -1)

	// 匹配 HTML 图片格式 <img src="url" ...>
	htmlRegex := regexp.MustCompile(`<img\s+[^>]*src=["']([^"']+)["'][^>]*>`)
	htmlMatches := htmlRegex.FindAllStringSubmatch(content, -1)

	var images []string
	for _, match := range mdMatches {
		if len(match) > 1 {
			images = append(images, match[1])
		}
	}
	for _, match := range htmlMatches {
		if len(match) > 1 {
			images = append(images, match[1])
		}
	}
	return images
}

// ValidatePassword 验证密码格式（5-18位，必须包含数字和字母，且只能是数字和字母）
func ValidatePassword(pwd string) error {
	if len(pwd) < 5 || len(pwd) > 18 {
		return errors.New("密码长度必须在5-18位之间")
	}
	hasLetter := false
	hasNumber := false
	for _, c := range pwd {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		} else if c >= '0' && c <= '9' {
			hasNumber = true
		} else {
			return errors.New("密码只能包含数字和字母")
		}
	}
	if !hasLetter || !hasNumber {
		return errors.New("密码必须同时包含数字和字母")
	}
	return nil
}
