package timex_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/odycenter/std-library/timex"
)

func TestParse(t *testing.T) {
	fmt.Println(int(math.Ceil(timex.Cs().Float64()-float64(1645574400)) / 86400))
	fmt.Println(timex.Parse("2021/12/19 8:36:29", "2006-01-02 15:04:05", true).Int64())
	fmt.Println(timex.Cs().Int64())
}

func TestFirstDay(t *testing.T) {
	fmt.Println(timex.FirstDay().Week())
	fmt.Println(timex.FirstDay().Month())
	fmt.Println(timex.FirstDay().Year())
}

func TestFirstNight(t *testing.T) {
	fmt.Println(timex.FirstNight().Week())
	fmt.Println(timex.FirstNight().Month())
	fmt.Println(timex.FirstNight().Year())

}

func TestLastDay(t *testing.T) {
	fmt.Println(timex.LastDay().Week())
	fmt.Println(timex.LastDay().Month())
	fmt.Println(timex.LastDay().Year())
}

func TestLastNight(t *testing.T) {
	fmt.Println(timex.LastNight().Week())
	fmt.Println(timex.LastNight().Month())
	fmt.Println(timex.LastNight().Year())

}

func BenchmarkTimeGenerator(b *testing.B) {
	var d = timex.Range(timex.ZeroTime(), timex.ZeroTime().Add(time.Hour*8760), time.Hour*24)
	for i := 0; i < b.N; i++ {
		for {
			n := d.Next()
			if n.IsZero() {
				break
			}
			//fmt.Println(n)
		}
	}
}

func TestTimeGenerator(t *testing.T) {
	d := timex.Range(timex.ZeroTime(), timex.ZeroTime().Add(time.Hour*1024), time.Hour*24)
	for {
		n := d.Next()
		if n.IsZero() {
			return
		}
		fmt.Println(n)
	}
}

func TestBetween(t *testing.T) {
	b := timex.Between(timex.ZeroTime(), timex.ZeroTime().Add(time.Hour*24*7))
	fmt.Println(b.Second())
	fmt.Println(b.Minute())
	fmt.Println(b.Day())
	fmt.Println(b.Duration())
}
