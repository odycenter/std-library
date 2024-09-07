package imgex_test

import (
	"fmt"
	"github.com/odycenter/std-library/imgex"
	"image/color"
	"image/jpeg"
	"os"
	"path"
	"testing"
	"time"
)

func TestCompress(t *testing.T) {
	s := `1.jpg`
	img, _ := imgex.Compress(s, 128)
	if f, err := os.Create(path.Join(path.Dir(s), "1"+"_compress"+path.Ext(s))); err != nil {
		fmt.Println(err)
		return
	} else {
		defer f.Close()
		jpeg.Encode(f, img, nil)
	}
}

func TestDrawLine(t *testing.T) {
	d, err := imgex.NewDraw("1.png", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = d.DrawLine("ABCDEFG\nabcdefg\n测试测试\n123456\n!@#$%^&*()", &imgex.DrawTextOption{
		Font:     "",
		FontSize: 24,
		DPI:      72,
		X:        10,
		Y:        50,
		Color:    color.RGBA{255, 0, 0, 255},
		Output:   "new.png",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestDrawText(t *testing.T) {
	d, err := imgex.NewDraw("1.png", &imgex.DrawOption{ExtWidth: 1000})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = d.DrawText("ABCDEFG\nabcdefg\n测试测试\n123456\n!@#$%^&*()", &imgex.DrawTextOption{
		Font:        "",
		FontSize:    100,
		DPI:         72,
		StartX:      d.Weight(),
		X:           50,
		StartY:      0,
		Y:           0,
		Color:       color.RGBA{255, 0, 0, 255},
		Output:      time.Now().Format("20060102150405.000000") + ".png",
		LineSpacing: 1.2,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
