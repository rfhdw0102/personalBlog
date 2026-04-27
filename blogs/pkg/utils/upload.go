package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var allowedImageExt = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
}

func ValidateImage(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedImageExt[ext]; !ok {
		return errors.New("只允许上传 jpg, jpeg, png, gif 格式")
	}
	return nil
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func BuildUploadFilename(userID int, ext string) string {
	return fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
}
