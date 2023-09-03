package logs

import (
	"fmt"
	"path"
	"strconv"
)

var formatterMap = make(map[string]LogFormatter, 4)

type LogFormatter interface {
	Format(lm *Msg) string
}

// PatternLogFormatter
// 提供了一个快速格式化的方法例如：
// tes := &PatternLogFormatter{Pattern: "%F:%n|%w %t>> %m", WhenFormat : "2006-01-02"}
// RegisterFormatter("tes", tes )
// SetGlobalFormatter("tes")
type PatternLogFormatter struct {
	Pattern    string
	WhenFormat string
}

func (p *PatternLogFormatter) getWhenFormatter() string {
	s := p.WhenFormat
	if s == "" {
		s = "2006/01/02 15:04:05.123" // default style
	}
	return s
}

func (p *PatternLogFormatter) Format(lm *Msg) string {
	return p.ToString(lm)
}

// RegisterFormatter 注册一个格式。通常你应该使用它来扩展你的自定义格式化程序
// 例如:
// RegisterFormatter("my-fmt", &MyFormatter{})
// logs.SetFormatter(Console, `{"formatter": "my-fmt"}`)
func RegisterFormatter(name string, formatter LogFormatter) {
	formatterMap[name] = formatter
}

func GetFormatter(name string) (LogFormatter, bool) {
	res, ok := formatterMap[name]
	return res, ok
}

// ToString 将log格式化为目标模式
// 'w' when, 'm' msg,'f' filename，'F' full path，'n' line number
// 'l' level number, 't' prefix of level type, 'T' full name of level type
func (p *PatternLogFormatter) ToString(lm *Msg) string {
	s := []rune(p.Pattern)
	msg := fmt.Sprintf(lm.Msg, lm.Args...)
	m := map[rune]string{
		'w': lm.When.Format(p.getWhenFormatter()),
		'm': msg,
		'n': strconv.Itoa(lm.LineNumber),
		'l': strconv.Itoa(int(lm.Level)),
		't': levelPrefix[lm.Level],
		'T': levelNames[lm.Level],
		'F': lm.FilePath,
	}
	_, m['f'] = path.Split(lm.FilePath)
	res := ""
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '%' {
			if k, ok := m[s[i+1]]; ok {
				res += k
				i++
				continue
			}
		}
		res += string(s[i])
	}
	return res
}
