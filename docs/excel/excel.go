// Package excel Excel文档生成操作
package excel

import (
	"errors"
	"github.com/tealeg/xlsx"
	"io"
)

type Base struct {
	file  *xlsx.File
	sheet *xlsx.Sheet
}

// New 创建excel文件对象
func New(excelName ...string) *Base {
	res := new(Base)
	res.file = xlsx.NewFile()
	// 自定义读取库
	if len(excelName) <= 0 {
		excelName[0] = "Sheet1"
	}
	sheet, err := res.file.AddSheet(excelName[0])
	if err != nil {
		return nil
	}
	res.sheet = sheet
	return res
}

// AddTitle 添加Excel标题
func (b *Base) AddTitle(titles []string) error {
	if len(titles) <= 0 {
		return errors.New("标题不能为空")
	}
	row := b.sheet.AddRow()
	for _, title := range titles {
		cell := row.AddCell()
		cell.Value = title
	}
	return nil
}
func (b *Base) AddRow() *xlsx.Row {
	return b.sheet.AddRow()
}

// SaveFile 保存文件
func (b *Base) SaveFile(fileName string) error {
	if fileName == "" {
		return errors.New("文件名不能为空")
	}
	err := b.file.Save(fileName)
	return err
}

// SaveWrite 保存文件
func (b *Base) SaveWrite(writer io.Writer) error {
	err := b.file.Write(writer)
	return err
}
