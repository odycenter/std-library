package zip

import (
	"github.com/alexmullins/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Compress(src, dst string, pwd ...string) error {
	zf, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zf.Close()
	archive := zip.NewWriter(zf)
	defer archive.Close()
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, filepath.Dir(src)+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		// 设置密码
		if len(pwd) > 0 {
			header.SetPassword(pwd[0])
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})
	return err
}
