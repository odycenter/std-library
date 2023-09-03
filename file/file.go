// Package file 封装的文件操作类
package file

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Exists 判断文件是否存在  存在返回 true 不存在返回false
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

// Hash 文件夹散列   6512bd43d9caa6e02c990b0a82652dca => 6/5/1/2/
func Hash(hash string, split int) (path string) {
	if split < 0 || split > 32 {
		return ""
	}
	rs := []rune(hash)
	for i := 0; i < split; i++ {
		end := i + 1
		path += string(rs[i:end]) + string(os.PathSeparator)
	}
	return
}

// Copy 复制文件
func Copy(src, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	return err
}

// Walk 遍历path路径并返回目录中的所有文件路径
func Walk(path string) (ret []string, err error) {
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d == nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ret = append(ret, path)
		return nil
	})
	return
}
