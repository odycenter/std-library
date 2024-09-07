// Package imgex 图片处理方法封装
package imgex

import (
	"bytes"
	"embed"
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/odycenter/std-library/crypto/md5"
	"github.com/odycenter/std-library/file"
	"github.com/odycenter/std-library/imgex/resize"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
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
	if img, _, err = image.Decode(f); err != nil {
		return nil, err
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

func SaveImageCompress(imageData []byte, basePath, filePathBase, tempPath string, needDecode bool) (*Image, error) {
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

	imgData, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}
	// 创建一个缓冲区用于存储压缩后的图像数据
	var compressedImageBuffer bytes.Buffer
	// 设置 JPEG 编码参数
	options := &jpeg.Options{Quality: 80}
	errImg := jpeg.Encode(&compressedImageBuffer, imgData, options)
	if errImg != nil {
		return nil, errImg
	}
	imageData = compressedImageBuffer.Bytes()

	err = os.WriteFile(tempFullName, imageData, 0644)
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

// Draw 绘制对象
type Draw struct {
	imgType string
	img     image.Image
	weight  int
	height  int
}

// Weight 返回原图的宽度
func (d *Draw) Weight() int {
	return d.weight
}

// Height 返回原图的高度
func (d *Draw) Height() int {
	return d.height
}

var (
	//go:embed simsun.ttc
	ttf            embed.FS
	defaultFont, _ = ttf.ReadFile("simsun.ttc") //默认字体
)

// DrawOption 绘制对象配置
type DrawOption struct {
	ExtWidth  int //设置扩展区域宽
	ExtHeight int //设置扩展区域高
}

// NewDraw 创建绘画对象
// Draw 实现了在图片或图片扩展区域叠加绘制图形或文字
func NewDraw(path string, opt *DrawOption) (*Draw, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	img, name, err := image.Decode(fs)
	if err != nil {
		return nil, err
	}
	weight := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if opt != nil {
		img = extend(img, opt.ExtWidth, opt.ExtHeight)
	}
	return &Draw{name, img, weight, height}, nil
}

func extend(origin image.Image, width, height int) image.Image {
	dst := imaging.New(origin.Bounds().Dx()+width, origin.Bounds().Dy()+height, color.White)
	draw.Draw(dst, dst.Bounds(), origin, image.Point{}, draw.Src)
	return dst
}

// DrawTextOption 文字绘制配置
type DrawTextOption struct {
	Font        string      //字体文件路径,为空则使用默认字体 仿宋
	font        []byte      //读出的字体
	FontSize    float64     //字号
	DPI         float64     //图片分辨率 像素/英寸
	StartX      int         //x坐标偏移起点
	X           int         //当 MatchX = FitX 时， X 将表示起点在绘制区域以外的偏移坐标
	StartY      int         //y坐标偏移起点
	Y           int         //当 MatchY = FitY 时， Y 将表示起点在绘制区域以外的偏移坐标
	Color       color.Color //颜色color.RGBA{R: 255, G: 0, B: 0, A: 255}
	Output      string      //结果图片输出路径 {path}/{name}.png
	LineSpacing float64     //多行文字行间距，1.2为1.2倍行距
}

func (o *DrawTextOption) getFont() []byte {
	if o.Font == "" {
		o.font = defaultFont
	}

	return o.font
}

func (o *DrawTextOption) getX() int {
	return o.StartX + o.X
}

func (o *DrawTextOption) getY() int {
	return o.StartY + o.Y
}

// DrawLine 叠加单行文字到图片
func (d *Draw) DrawLine(line string, opt *DrawTextOption) error {
	bounds := d.img.Bounds()
	img := image.NewRGBA(bounds)
	font, err := freetype.ParseFont(opt.getFont())
	if err != nil {
		return err
	}
	fc := freetype.NewContext()
	fc.SetDPI(opt.DPI)
	fc.SetFont(font)
	fc.SetFontSize(opt.FontSize)
	fc.SetClip(bounds)
	fc.SetDst(img)
	fc.SetSrc(image.NewUniform(opt.Color))

	draw.Draw(img, bounds, d.img, image.Point{}, draw.Src)

	pt := freetype.Pt(opt.getX(), opt.getY())
	_, err = fc.DrawString(line, pt)
	if err != nil {
		return err
	}
	err = d.save(opt.Output, img)
	return err
}

// DrawText 叠加多行文字到图片
// 文字换行使用 \n 分割
func (d *Draw) DrawText(text string, opt *DrawTextOption) error {
	bounds := d.img.Bounds()
	img := image.NewRGBA(bounds)
	font, err := freetype.ParseFont(opt.getFont())
	if err != nil {
		return err
	}
	fc := freetype.NewContext()
	fc.SetDPI(opt.DPI)
	fc.SetFont(font)
	fc.SetFontSize(opt.FontSize)
	fc.SetClip(bounds)
	fc.SetDst(img)
	fc.SetSrc(image.NewUniform(opt.Color))

	draw.Draw(img, bounds, d.img, image.Point{}, draw.Src)

	lines := strings.Split(text, "\n")
	y := opt.getY()
	for _, line := range lines {
		pt := fixed.Point26_6{
			X: fixed.I(opt.getX()),
			Y: fixed.I(y),
		}
		_, err = fc.DrawString(line, pt)
		if err != nil {
			return err
		}
		y += fc.PointToFixed(opt.FontSize * opt.LineSpacing).Ceil()
	}
	err = d.save(opt.Output, img)
	return err
}

func (d *Draw) save(output string, img image.Image) (err error) {
	fs, err := os.Create(output)
	if err != nil {
		return err
	}
	defer fs.Close()
	switch d.imgType {
	case "png":
		err = png.Encode(fs, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(fs, img, &jpeg.Options{Quality: 100})
	}
	return err
}
