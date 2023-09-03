package imgex

import (
	"bytes"
	orientation "github.com/takumakei/exif-orientation"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

// Thumbnail 缩略图
// filePath 图片地址(/data/static/upload/image/6/5/1/2/565fdsa56daasd456asd65sda)
// suffix 后缀(normal,small,medium...)
// ext 扩展名(.jpg...)
// size 缩略图尺寸大小
func Thumbnail(filePath, suffix, ext string, size int) (Image, error) {
	image := Image{}
	image.ImgPath = filePath + suffix + ext
	srcFile := filePath + ext
	log.Println("srcFile : ", srcFile)
	//打开图片
	p, err := os.Open(srcFile)
	if err != nil {
		return image, err
	}
	defer p.Close()
	//区分jpg和jpeg
	if ext == ".jpg" || ext == ".jpeg" {
		//将jpeg解码为image.Image
		img, err := jpeg.Decode(p)
		if err != nil {
			return image, err
		}
		b, err := os.ReadFile(srcFile)
		if err != nil {
			return image, err
		}
		r := bytes.NewReader(b)
		r.Reset(b)
		or, err := orientation.Read(r)
		if err == nil {
			img = ReverseOrientation(img, or.String())
		}
		if err = image.resizeImage(img, image.ImgPath, ext, size); err != nil {
			return image, err
		}

	} else if ext == ".png" {
		//将jpeg解码为image.Image
		img, err := png.Decode(p)
		if err != nil {
			return image, err
		}
		if err = image.resizeImage(img, image.ImgPath, ext, size); err != nil {
			return image, err
		}
	}
	return image, nil
}

type ImageThumbnail struct {
	Url    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}

func SaveImageThumbnail(imageData []byte, basePath, filePathBase, tempPath string, smallSize, mediumSize int, needDecode bool) (map[string]ImageThumbnail, error) {
	normalImage, err := SaveImage(imageData, basePath, filePathBase, tempPath, needDecode)
	if err != nil {
		return nil, err
	}
	// 对必要信息进行赋值
	imgExt := normalImage.Ext

	index := strings.LastIndex(normalImage.ImgPath, ".")
	filePath := normalImage.ImgPath[:index]
	//小缩略图
	sThumbnailInf, err := Thumbnail(filePath, "_small", imgExt, smallSize)
	if err != nil {
		log.Println("DebugSaveImageThumbnail err:", err)
		sThumbnailInf = *normalImage
	}
	//中缩略图
	mThumbnailInf, err := Thumbnail(filePath, "_medium", imgExt, mediumSize)
	if err != nil {
		log.Println("DebugSaveImageThumbnail err:", err)
		mThumbnailInf = *normalImage
	}
	sThumbnailInf.ImgPath = strings.ReplaceAll(sThumbnailInf.ImgPath, string(os.PathSeparator), "/")
	mThumbnailInf.ImgPath = strings.ReplaceAll(mThumbnailInf.ImgPath, string(os.PathSeparator), "/")

	normalPath := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(normalImage.ImgPath, basePath+"\\", ""), basePath+"/", ""), "\\", "/")
	smallPath := strings.ReplaceAll(strings.ReplaceAll(sThumbnailInf.ImgPath, basePath+"/", ""), "\\", "/")
	mediumPath := strings.ReplaceAll(strings.ReplaceAll(mThumbnailInf.ImgPath, basePath+"/", ""), "\\", "/")

	imageInfo := map[string]ImageThumbnail{
		"normal": {Url: normalPath, Width: normalImage.Width, Height: normalImage.Height},
		"small":  {Url: smallPath, Width: sThumbnailInf.Width, Height: sThumbnailInf.Height},
		"medium": {Url: mediumPath, Width: mThumbnailInf.Width, Height: mThumbnailInf.Height},
	}
	return imageInfo, nil
}
