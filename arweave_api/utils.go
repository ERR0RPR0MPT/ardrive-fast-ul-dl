package arweave_api

import (
	"mime"
	"path/filepath"
)

func getMimeType(fileName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(fileName)
	// 根据扩展名查找 MIME 类型
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream" // 默认类型
	}
	return mimeType
}
