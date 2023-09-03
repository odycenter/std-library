// Package imgex 图片处理方法封装
package imgex

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/odycenter/std-library/crypto/md5"
	"github.com/odycenter/std-library/file"
)

// 算法
const (
	// NearestNeighbor 最近邻插值
	NearestNeighbor resize.InterpolationFunction = iota
	// Bilinear 插值
	Bilinear
	// Bicubic 插值（使用三次 Hermite 样条）
	Bicubic
	// MitchellNetravali Mitchell-Netravali 插值法
	MitchellNetravali
	// Lanczos2 Lanczos 插值 (a=2)
	Lanczos2
	// Lanczos3 Lanczos 插值 (a=3)
	Lanczos3
)

// Compress 压缩图片
// fs:源文件位置
// height:高度
// t:指定的压缩算法
func Compress(fs string, height uint, t ...resize.InterpolationFunction) (newImg image.Image, err error) {
	reg, _ := regexp.Compile(`^.*\.((png)|(jpg)|(jpeg))$`)
	if !reg.MatchString(fs) {
		return nil, fmt.Errorf("support png|jpg|jpeg only")
	}
	var f *os.File
	if f, err = os.Open(fs); err != nil {
		return nil, err
	}
	defer f.Close()
	var img image.Image
	switch {
	case strings.HasSuffix(f.Name(), "png"):
		if img, err = png.Decode(f); err != nil {
			return nil, err
		}
	case strings.HasSuffix(f.Name(), "jpg") || strings.HasSuffix(f.Name(), "jpeg"):
		if img, err = jpeg.Decode(f); err != nil {
			return nil, err
		}
	}
	t = append(t, resize.Lanczos3)
	newImg = resize.Resize(0, height, img, t[0])
	return newImg, nil
}

// ReverseOrientation 换向
func ReverseOrientation(img image.Image, o string) *image.NRGBA {
	switch o {
	case "1":
		return imaging.Clone(img)
	case "2":
		return imaging.FlipV(img)
	case "3":
		return imaging.Rotate180(img)
	case "4":
		return imaging.Rotate180(imaging.FlipV(img))
	case "5":
		return imaging.Rotate270(imaging.FlipV(img))
	case "6":
		return imaging.Rotate270(img)
	case "7":
		return imaging.Rotate90(imaging.FlipV(img))
	case "8":
		return imaging.Rotate90(img)
	}
	return imaging.Clone(img)
}

type Image struct {
	ImgPath string
	Width   uint
	Height  uint
	Ext     string
}

// 按size大小等比缩放
func (i *Image) resizeImage(img image.Image, storageImg string, ext string, size int) error {
	c := img.Bounds()
	if c.Dx()+c.Dy() > 200 {
		sideLen := math.Max(float64(c.Dx()), float64(c.Dy()))
		ratio := float64(size) / sideLen
		i.Width = uint(float64(c.Dx()) * ratio)
		i.Height = uint(float64(c.Dy()) * ratio)
	} else {
		i.Width = uint(c.Dx())
		i.Height = uint(c.Dy())
		return nil
	}
	m := resize.Resize(i.Width, i.Height, img, resize.MitchellNetravali)
	out, err := os.Create(storageImg)
	if err != nil {
		return err
	}
	i.ImgPath = storageImg
	defer out.Close()
	//将新图像写入文件
	if ext == ".jpg" || ext == ".jpeg" {
		err = jpeg.Encode(out, m, nil)
	} else if ext == ".png" {
		err = png.Encode(out, m)
	}

	if err != nil {
		return err
	}
	return nil
}

func (i *Image) GetImageSize(picture string) error {
	p1, err := os.Open(picture)
	if err != nil {
		return err
	}
	defer p1.Close()

	// 先用jpeg解码，失败后换为png
	img, err := jpeg.Decode(p1)
	if err == nil {
		c := img.Bounds()
		i.Width = uint(c.Dx())
		i.Height = uint(c.Dy())
		i.Ext = ".jpg"
		return nil
	}

	p2, err := os.Open(picture)
	if err != nil {
		return err
	}
	defer p2.Close()

	img, err = png.Decode(p2)
	if err == nil {
		c := img.Bounds()
		i.Width = uint(c.Dx())
		i.Height = uint(c.Dy())
		i.Ext = ".png"
		return nil
	}
	return err
}

func SaveImage(imageData []byte, basePath, filePathBase, tempPath string, needDecode bool) (*Image, error) {
	if needDecode {
		var dst []byte
		base64.StdEncoding.Encode(dst, imageData)
		if len(imageData) == 0 {
			return nil, errors.New("Base64DecodeByte fail")
		}
	}
	imgExt := ".jpg"
	v := fmt.Sprint(time.Now().UnixNano() / 1000)
	hashPath := file.Hash(md5.Sum([]byte(v)).Hex(), 3)
	path := fmt.Sprintf(filePathBase, hashPath)            // upload/image/6/5/1/2/
	fullPath := basePath + string(os.PathSeparator) + path // /data/static/upload/image/6/5/1/2/
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return nil, err
	}
	fileName := v                   // 565fdsa56daasd456asd65sda
	filePath := fullPath + fileName // /data/static/upload/image/5/6/5/565fdsa56daasd456asd65sda
	// /tmp/im_tmp/565fdsa56daasd456asd65sda.xxx
	if err := os.MkdirAll(tempPath, os.ModePerm); err != nil {
		return nil, err
	}
	tempFullName := tempPath + string(os.PathSeparator) + fileName + imgExt
	// err = this.SaveToFile("img", tempFullName)
	err := os.WriteFile(tempFullName, imageData, 0644)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFullName)

	img := &Image{}
	img.ImgPath = string(os.PathSeparator) + path + fileName + imgExt
	if err := img.GetImageSize(tempFullName); err != nil {
		return nil, err
	}

	// 重命名
	fullName := filePath + img.Ext
	err = file.Copy(tempFullName, fullName)
	if err != nil {
		return nil, err
	}

	// 这里给本地的决对路径，上层逻辑层需要过滤
	img.ImgPath = fullName
	return img, nil
}
