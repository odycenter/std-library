// Package csv csv文档读取操作
package csv

import (
	"encoding/csv"
	"os"
)

// ReadAll 读取整个CSV文件，返回由行列组成的数组[row][column]string
func ReadAll(csvName string) ([][]string, error) {
	fs, err := os.Open(csvName)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(fs)
	return reader.ReadAll()
}

// Read 读取CSV文件，返回一行数据[row]string
func Read(csvName string) ([]string, error) {
	fs, err := os.Open(csvName)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(fs)
	return reader.Read()
}
