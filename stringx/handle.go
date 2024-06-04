package stringx

import (
	"github.com/gogf/gf/v2/text/gstr"
	"regexp"
	"strings"
	"unicode/utf8"
)

// TrimHtml 去除字符串的HTML
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

// Split 切分字符串
// needSpace是否保留空字符，默认false
func Split(s, sep string, needSpace ...bool) *T {
	var ns bool
	if len(needSpace) != 0 {
		ns = needSpace[0]
	}
	var arr []string
	for _, v := range strings.Split(s, sep) {
		if v == "" && !ns {
			continue
		}
		arr = append(arr, v)
	}
	return &T{arr}
}

// Join 将任何数组转为sep分割的字符串
func Join(in any, sep string) string {
	return gstr.JoinAny(in, sep)
}

// Hidden 用指定字符替换指定位置的字符
// Hidden("张三丰", 1, 1)
// 输出“张*丰”
// 第三个参数，默认为：'*'
func Hidden(src string, begin, end uint, hidden ...rune) string {
	placeholder := '*'
	if len(hidden) > 0 {
		placeholder = hidden[0]
	}
	all := []rune(src)
	pre := all[:begin]
	count := uint(len(all))
	if begin > count {
		return src
	}
	if end > count {
		end = count
	}

	for i := begin; i < end; i++ {
		pre = append(pre, placeholder)
	}

	return string(pre) + string(all[end:])
}

// HiddenUnknown 用指定长度字符替换指定位置的字符
// HiddenUnknown("0123456789", 1, 5)
// 输出“0***789”
// 第三个参数，默认为："***"
func Mask(src string, begin, end uint, hidden ...string) string {
	placeholder := "***"
	if len(hidden) > 0 {
		placeholder = hidden[0]
	}
	all := []rune(src)
	pre := all[:begin]
	count := uint(len(all))
	if begin > count {
		return src
	}
	if end > count {
		end = count
	}

	return string(pre) + placeholder + string(all[end:])
}

// Sub 截取字符串
func Sub(src string, begin, end uint) string {
	all := []rune(src)
	count := uint(len(all))
	if begin > count {
		return ""
	}
	if end > count {
		end = count
	}

	return string(all[begin:end])
}

// Trim 去除字符串中所有相关字符
// fn 返回true则删除
func Trim(in string, fn func(r rune) bool) string {
	var rt []rune
	for _, r := range in {
		if !fn(r) {
			rt = append(rt, r)
		}
	}
	return string(rt)
}

// Camel2Snake 驼峰转蛇形 XxYy to xx_yy , XxYY to xx_yy
func Camel2Snake(s string) string {
	reg := regexp.MustCompile("([a-z0-9])([A-Z])")
	snakeCase := reg.ReplaceAllString(s, "${1}_${2}")
	// Convert to lower case
	snakeCase = strings.ToLower(snakeCase)
	return snakeCase
}

// Len 以UTF-8编码获取string的长度，主要试用于中文长度
func Len(s string) int {
	return utf8.RuneCountInString(s)
}

// RemoveDuplicates 字符串去重
// exception 例外的字符，不参与去重
func RemoveDuplicates(s string, exception ...rune) string {
	unique := make(map[rune]bool)
	var result []rune
	for _, r := range s {
		if !unique[r] || (len(exception) != 0 && exception[0] == r) {
			unique[r] = true
			result = append(result, r)
		}
	}
	return string(result)
}
