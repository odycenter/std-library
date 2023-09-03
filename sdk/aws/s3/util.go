package s3

import "strings"

func ContextType(fileName string) string {
	fileName = strings.ToLower(fileName)
	if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(fileName, ".png") {
		return "image/png"
	} else if strings.HasSuffix(fileName, ".html") || strings.HasSuffix(fileName, ".htm") {
		return "text/html"
	} else if strings.HasSuffix(fileName, ".css") {
		return "text/css"
	} else if strings.HasSuffix(fileName, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(fileName, ".json") {
		return "application/json" // text/plain
	} else {
		return "application/octet-stream" // 下载
	}
}
