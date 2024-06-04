// Package timex 日期时间扩展
package timex

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

// T 结果类型转换结构体
type T struct {
	t int64
}

type Numeric interface {
	uint8 |
		uint16 |
		uint32 |
		uint64 |
		int8 |
		int16 |
		int32 |
		int64 |
		float32 |
		float64 |
		int |
		uint
}

func (t *T) Uint64() uint64 {
	return uint64(t.t)
}

func (t *T) Uint32() uint32 {
	return uint32(t.t)
}

func (t *T) Int64() int64 {
	return t.t
}

func (t *T) Int() int {
	return int(t.t)
}

func (t *T) String() string {
	return fmt.Sprint(t.t)
}

func (t *T) Bytes() []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(t.t))
	return buf
}

func (t *T) Float32() float32 {
	return float32(t.t)
}

func (t *T) Float64() float64 {
	return float64(t.t)
}

func (t *T) Duration() time.Duration {
	return time.Duration(t.t)
}

// CNs Current NanoSecond
func CNs() *T {
	return &T{time.Now().UnixNano()}
}

// CMs Current MicroSecond
func CMs() *T {
	return &T{time.Now().UnixNano() / 1e6}
}

// CUs Current MilliSecond
func CUs() *T {
	return &T{time.Now().UnixNano() / 1e3}
}

// Cs Current Second
func Cs() *T {
	return &T{time.Now().Unix()}
}

// After Return the Time From Now After Duration
func After(duration time.Duration) *T {
	return &T{time.Now().Add(duration).Unix()}
}

// CFormat Current time format
func CFormat(layout string) string {
	return time.Now().Format(layout)
}

// Format Any time format
func Format[T Numeric](timestamp T, layout string) string {
	return time.Unix(int64(timestamp), 0).Format(layout)
}

// Time Any timestamp to time
func Time[T Numeric](timestamp T) time.Time {
	return time.Unix(int64(timestamp), 0)
}

// AtToday Check timestamp is in today time range
func AtToday[T Numeric](timestamp T) bool {
	t := time.Now()
	today := []int64{
		time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix(),
		time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location()).Unix(),
	}
	ts := time.Unix(int64(timestamp), 0).Unix()
	return ts >= today[0] && ts <= today[1]
}

// Parse DataTime parse with layout to unix
// set compatible true can make func try some preset layouts
var presetLayouts = []string{
	"2006-01-02 15:04:05",
	"2006-01-2 15:04:05",
	"2006-1-2 15:04:05",
	"2006-1-2 15:04",
	"2006/1/2 15:04",
	"2006/1/2 15:04:05",
	"2006/1/02 15:04:05",
	"2006/01/02 15:04:05",
}

// Parse 日期转其他类型值
func Parse(dt, layout string, compatible ...bool) *T {
	var t time.Time
	if dt == "" {
		return &T{0}
	}
	t, err := time.ParseInLocation(layout, dt, LoadLocation())
	if err != nil {
		if len(compatible) != 0 && compatible[0] {
			for _, presetLayout := range presetLayouts {
				if t, e := time.ParseInLocation(presetLayout, dt, LoadLocation()); e == nil {
					return &T{t.Unix()}
				}
			}
		}
		fmt.Println(err)
		return &T{0}
	}
	return &T{t.Unix()}
}

// ParseT 日期转时间对象 time.Time
func ParseT(dt, layout string, compatible ...bool) *time.Time {
	if dt == "" {
		return nil
	}
	t, err := time.ParseInLocation(layout, dt, LoadLocation())
	if err != nil {
		if len(compatible) != 0 && compatible[0] {
			for _, presetLayout := range presetLayouts {
				if t, e := time.ParseInLocation(presetLayout, dt, LoadLocation()); e == nil {
					return &t
				}
			}
		}
		fmt.Println(err)
		return nil
	}
	return &t
}

// ParseInLoc 日期转其他类型值并使用location
func ParseInLoc(dt, layout string, loc *time.Location, compatible ...bool) *T {
	var t time.Time
	if loc == nil {
		loc = LoadLocation()
	}
	if dt == "" {
		return &T{0}
	}
	t, err := time.ParseInLocation(layout, dt, loc)
	if err != nil {
		if len(compatible) != 0 && compatible[0] {
			for _, presetLayout := range presetLayouts {
				if t, e := time.ParseInLocation(presetLayout, dt, loc); e == nil {
					return &T{t.Unix()}
				}
			}
		}
		fmt.Println(err)
		return &T{0}
	}
	return &T{t.Unix()}
}

// ParseTInLoc 日期转时间对象 time.Time 并使用location
func ParseTInLoc(dt, layout string, loc *time.Location, compatible ...bool) *time.Time {
	var t time.Time
	if loc == nil {
		loc = LoadLocation()
	}
	if dt == "" {
		return nil
	}
	t, err := time.ParseInLocation(layout, dt, loc)
	if err != nil {
		if len(compatible) != 0 && compatible[0] {
			for _, presetLayout := range presetLayouts {
				if t, e := time.ParseInLocation(presetLayout, dt, loc); e == nil {
					return &t
				}
			}
		}
		fmt.Println(err)
		return nil
	}
	return &t
}

// WeekZhou 周
var WeekZhou = map[string]string{
	"Monday":    "周一",
	"Tuesday":   "周二",
	"Wednesday": "周三",
	"Thursday":  "周四",
	"Friday":    "周五",
	"Saturday":  "周六",
	"Sunday":    "周日",
}

// WeekXingQi 星期
var WeekXingQi = map[string]string{
	"Monday":    "星期一",
	"Tuesday":   "星期二",
	"Wednesday": "星期三",
	"Thursday":  "星期四",
	"Friday":    "星期五",
	"Saturday":  "星期六",
	"Sunday":    "星期日",
}

// Week Any timestamp to chinese week
func Week(timestamp int64, weekMap map[string]string) string {
	return weekMap[time.Unix(timestamp, 0).Weekday().String()]
}

type FD struct {
	t      time.Time
	offset time.Duration
}

// Month 月
func (ft *FD) Month() time.Time {
	return ZeroTime(ft.t.AddDate(0, 0, -ft.t.Day()+1)).Add(ft.offset)
}

// Year 年
func (ft *FD) Year() time.Time {
	return ZeroTime(ft.t.AddDate(0, int(-ft.t.Month())+1, -ft.t.Day()+1)).Add(ft.offset)
}

// Week 周
func (ft *FD) Week() time.Time {
	offset := int(time.Monday - ft.t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return ZeroTime(ft.t).AddDate(0, 0, offset).Add(ft.offset)
}

// FirstDay 年/月/周 的第一天 零点
func FirstDay(in ...time.Time) *FD {
	var t time.Time
	if len(in) == 0 {
		t = time.Now()
	} else {
		t = in[0]
	}
	return &FD{t, 0}
}

// FirstNight 年/月/周 的第一天 23:59:59
func FirstNight(in ...time.Time) *FD {
	var t time.Time
	if len(in) == 0 {
		t = time.Now()
	} else {
		t = in[0]
	}
	return &FD{t, 86399 * time.Second}
}

type LD struct {
	t      time.Time
	offset time.Duration
}

// Month 月
func (lt *LD) Month() time.Time {
	return ZeroTime(lt.t.AddDate(0, 1, -lt.t.Day())).Add(lt.offset)
}

// Year 年
func (lt *LD) Year() time.Time {
	return ZeroTime(lt.t.AddDate(0, 12-int(lt.t.Month())+1, -lt.t.Day())).Add(lt.offset)
}

// Week 周
func (lt *LD) Week() time.Time {
	return ZeroTime(lt.t).AddDate(0, 0, 7-int(lt.t.Weekday())).Add(lt.offset)
}

// LastDay 年/月/周 的最后一天 零点
func LastDay(in ...time.Time) *LD {
	var t time.Time
	if len(in) == 0 {
		t = time.Now()
	} else {
		t = in[0]
	}
	return &LD{t, 0}
}

// LastNight 年/月/周 的最后一天 23:59:59
func LastNight(in ...time.Time) *LD {
	var t time.Time
	if len(in) == 0 {
		t = time.Now()
	} else {
		t = in[0]
	}
	return &LD{t, 86399 * time.Second}
}

// Zero Any(Current) timestamp's zero time at that(to-) day
func Zero(timestamp ...int64) int64 {
	var t time.Time
	if len(timestamp) == 0 {
		t = time.Now()
	} else {
		t = time.Unix(timestamp[0], 0)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}

// ZeroTime Any(Current) zero time at that(to-) day
func ZeroTime(in ...time.Time) time.Time {
	var t time.Time
	if len(in) == 0 {
		t = time.Now()
	} else {
		t = in[0]
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// Night Any(Current) timestamp's 23:59:59 time at that(to-) day
func Night(timestamp ...int64) int64 {
	var t time.Time
	if len(timestamp) == 0 {
		t = time.Now()
	} else {
		t = time.Unix(timestamp[0], 0)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix()
}

// NightTime Any(Current) timestamp's 23:59:59 time at that(to-) day
func NightTime(in ...time.Time) time.Time {
	var t time.Time
	if len(in) == 0 {
		t = time.Now()
	} else {
		t = in[0]
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

// IsFormattedDate 是否为指定模式的日期格式
func IsFormattedDate(date, format string) bool {
	_, err := time.Parse(date, format)
	return err == nil
}

var EOT = time.Time{} //结束标记

// RTG 范围时间生成器
type RTG struct {
	s    time.Time
	e    time.Time
	step time.Duration
	i    int
}

func (t *RTG) Reset() {
	t.i = 0
}

func (t *RTG) Next() time.Time {
	n := t.s.Add(t.step * time.Duration(t.i))
	if n.After(t.e) {
		return EOT
	}
	t.i++
	return n
}

// FuncNext 已回调方式处理范围时间逻辑
func (t *RTG) FuncNext(fn func(time.Time) bool) {
	for {
		n := t.s.Add(t.step * time.Duration(t.i))
		if n.Sub(t.e) > 0 {
			return
		}
		if !fn(n) {
			return
		}
		t.i++
	}
}

// Range 创建一个时间范围生成器
func Range(begin, end time.Time, step time.Duration) *RTG {
	return &RTG{begin, end, step, 0}
}

type BWT struct {
	d *time.Duration
}

// Nanoseconds 转为纳秒(ns)
func (b *BWT) Nanoseconds() int {
	if b.d == nil {
		return 0
	}
	return int(b.d.Nanoseconds())
}

// Microseconds 转为微秒(μs)
func (b *BWT) Microseconds() int {
	if b.d == nil {
		return 0
	}
	return int(b.d.Microseconds())
}

// Milliseconds 转为毫秒(ms)
func (b *BWT) Milliseconds() int {
	if b.d == nil {
		return 0
	}
	return int(b.d.Milliseconds())
}

// Second 转为秒(s)
func (b *BWT) Second() int {
	if b.d == nil {
		return 0
	}
	return int(math.Ceil(b.d.Seconds()))
}

// Minute 转为分(min.)
func (b *BWT) Minute() int {
	if b.d == nil {
		return 0
	}
	return int(math.Ceil(b.d.Minutes()))
}

// Hour 转为小时(hr.)
func (b *BWT) Hour() int {
	if b.d == nil {
		return 0
	}
	return int(math.Ceil(b.d.Hours()))
}

// Day 转为天(day)
func (b *BWT) Day() int {
	if b.d == nil {
		return 0
	}
	return int(math.Ceil(b.d.Hours() / 24))
}

// Duration 返回Dur对象
func (b *BWT) Duration() time.Duration {
	if b.d == nil {
		return 0
	}
	return *b.d
}

// Between 两个时间的差
func Between(a, b time.Time) *BWT {
	d := a.Sub(b)
	return &BWT{&d}
}

// BetweenDT 两个时间DataTime的差
func BetweenDT(adt, bdt, layout string) *BWT {
	a, err := time.Parse(layout, adt)
	if err != nil {
		return &BWT{nil}
	}
	b, err := time.Parse(layout, bdt)
	if err != nil {
		return &BWT{nil}
	}
	d := a.Sub(b)
	return &BWT{&d}
}

// LoadLocation 按照传入的IANA获取Location
// 关于IANA可以参见 IANA.go
func LoadLocation(loc ...string) *time.Location {
	if len(loc) == 0 {
		loc = append(loc, "Local")
	}
	location, err := time.LoadLocation(loc[0])
	if err != nil {
		location, _ = time.LoadLocation("Local")
	}
	return location
}

// DurSec 数字转为秒级间隔(s)
func DurSec[T Numeric](t T) time.Duration {
	return time.Duration(t) * time.Second
}

// DurMin 数字转为分级间隔(m)
func DurMin[T Numeric](t T) time.Duration {
	return time.Duration(t) * time.Minute
}

// DurHour 数字转为小时级间隔(h)
func DurHour[T Numeric](t T) time.Duration {
	return time.Duration(t) * time.Hour
}
