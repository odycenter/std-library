// Package gzip gzip压缩与解压缩
package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Compress 使用gzip进行压缩
func Compress(data []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(data)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

// Decompress 使用gzip进行解压
func Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
