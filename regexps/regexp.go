/*
Package regexps 是一个 regexp 包，它的接口类似于 Go 的框架 regexp 引擎，但使用
更多功能在幕后完整的正则表达式引擎。

它没有恒定的时间保证，但它允许回溯并且与 Perl5 和 .NET 兼容。
使用 regexp 包中的 RE2 引擎可能会更好，并且只有在以下情况下才应该使用它：
需要编写非常复杂的模式或需要与 .NET 兼容。
*/
package regexps

import (
	"errors"
	"math"
	"std-library/regexps/syntax"
	"strconv"
	"sync"
	"time"
)

// DefaultMatchTimeout 运行正则表达式匹配时使用的默认超时 - “forever”
var DefaultMatchTimeout = time.Duration(math.MaxInt64)

// Regexp 是编译后的正则表达式的表示。
// 正则表达式对于多个 goroutine 并发使用是安全的。
type Regexp struct {
	// 如果匹配时间（大约）超过
	// MatchTimeout。这是一项安全检查，以防匹配
	// 遇到灾难性的回溯。默认值
	// (DefaultMatchTimeout) 导致所有超时检查
	// 抑制。
	MatchTimeout time.Duration

	// read-only after Compile
	pattern string       // as passed to Compile
	options RegexOptions // options

	caps     map[int]int    // capnum->index
	capNames map[string]int //capture group name -> index
	capsList []string       //sorted list of capture group names
	capsize  int            // size of the capture array

	code *syntax.Code // compiled program

	// cache of machines for running regexp
	muRun  sync.Mutex
	runner []*runner
}

// Compile 解析正则表达式并返回，如果成功，
// 可用于匹配文本的 Regexp 对象。
func Compile(expr string, opt RegexOptions) (*Regexp, error) {
	// parse it
	tree, err := syntax.Parse(expr, syntax.RegexOptions(opt))
	if err != nil {
		return nil, err
	}

	// translate it to code
	code, err := syntax.Write(tree)
	if err != nil {
		return nil, err
	}

	// return it
	return &Regexp{
		pattern:      expr,
		options:      opt,
		caps:         code.Caps,
		capNames:     tree.Capnames,
		capsList:     tree.Caplist,
		capsize:      code.Capsize,
		code:         code,
		MatchTimeout: DefaultMatchTimeout,
	}, nil
}

// MustCompile 与 Compile 类似，但如果无法解析表达式，则会出现恐慌。
// 它简化了保存已编译常规的全局变量的安全初始化
// 表达式。
func MustCompile(str string, opt RegexOptions) *Regexp {
	regexp, err := Compile(str, opt)
	if err != nil {
		panic(`regexps: Compile(` + quote(str) + `): ` + err.Error())
	}
	return regexp
}

// Escape 转义符为输入字符串中的任何特殊字符添加反斜杠
func Escape(input string) string {
	return syntax.Escape(input)
}

// Unescape 删除输入字符串中之前转义的特殊字符中的所有反斜杠
func Unescape(input string) (string, error) {
	return syntax.Unescape(input)
}

// SetTimeoutCheckPeriod 是一个调试函数，用于设置超时 goroutine 的睡眠周期的频率。
// 默认为 100 毫秒。设置较低的唯一好处是 1 个后台 goroutine 管理
// 在所有超时到期后，超时可能会稍早退出。请参阅 Github 问题 #63
func SetTimeoutCheckPeriod(d time.Duration) {
	clockPeriod = d
}

// StopTimeoutClock 只能在单元测试中使用，以防止时钟 goroutine 超时
// 避免看起来像泄漏的 goroutine
func StopTimeoutClock() {
	stopClock()
}

// String 返回用于编译正则表达式的源文本。
func (re *Regexp) String() string {
	return re.pattern
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

// RegexOptions 影响运行时和解析行为
// 对于每个特定的正则表达式。它们也可以在代码中设置
// 就像正则表达式模式本身一样。
type RegexOptions int32

const (
	None                    RegexOptions = 0x0
	IgnoreCase                           = 0x0001 // "i"
	Multiline                            = 0x0002 // "m"
	ExplicitCapture                      = 0x0004 // "n"
	Compiled                             = 0x0008 // "c"
	SingleLine                           = 0x0010 // "s"
	IgnorePatternWhitespace              = 0x0020 // "x"
	RightToLeft                          = 0x0040 // "r"
	Debug                                = 0x0080 // "d"
	ECMAScript                           = 0x0100 // "e"
	RE2                                  = 0x0200 // RE2 (regexp package) compatibility mode
	Unicode                              = 0x0400 // "u"
)

func (re *Regexp) RightToLeft() bool {
	return re.options&RightToLeft != 0
}

func (re *Regexp) Debug() bool {
	return re.options&Debug != 0
}

// Replace 搜索输入字符串并用替换文本替换找到的每个匹配项。
// Count 将限制尝试匹配的次数，而 startAt 将允许
// 我们在输入开始时跳过可能的匹配（向左或向右取决于 RightToLeft 选项）。
// 将 startAt 和 count 设置为 -1 以遍历整个字符串
func (re *Regexp) Replace(input, replacement string, startAt, count int) (string, error) {
	data, err := syntax.NewReplacerData(replacement, re.caps, re.capsize, re.capNames, syntax.RegexOptions(re.options))
	if err != nil {
		return "", err
	}
	//TODO: cache ReplacerData

	return replace(re, data, nil, input, startAt, count)
}

// ReplaceFunc 搜索输入字符串并使用求值器中的字符串替换找到的每个匹配项
// Count 将限制尝试匹配的次数，而 startAt 将允许
// 我们在输入开始时跳过可能的匹配（向左或向右取决于 RightToLeft 选项）。
// 将 startAt 和 count 设置为 -1 以遍历整个字符串。
func (re *Regexp) ReplaceFunc(input string, evaluator MatchEvaluator, startAt, count int) (string, error) {
	return replace(re, nil, evaluator, input, startAt, count)
}

// FindStringMatch 在输入字符串中搜索正则表达式匹配项
func (re *Regexp) FindStringMatch(s string) (*Match, error) {
	// convert string to runes
	return re.run(false, -1, getRunes(s))
}

// FindRunesMatch 在输入符文切片中搜索正则表达式匹配项
func (re *Regexp) FindRunesMatch(r []rune) (*Match, error) {
	return re.run(false, -1, r)
}

// FindStringMatchStartingAt 在输入字符串中搜索从 startAt 索引开始的正则表达式匹配项
func (re *Regexp) FindStringMatchStartingAt(s string, startAt int) (*Match, error) {
	if startAt > len(s) {
		return nil, errors.New("startAt must be less than the length of the input string")
	}
	r, startAt := re.getRunesAndStart(s, startAt)
	if startAt == -1 {
		// we didn't find our start index in the string -- that's a problem
		return nil, errors.New("startAt must align to the start of a valid rune in the input string")
	}

	return re.run(false, startAt, r)
}

// FindRunesMatchStartingAt 在输入符文切片中搜索从 startAt 索引开始的正则表达式匹配
func (re *Regexp) FindRunesMatchStartingAt(r []rune, startAt int) (*Match, error) {
	return re.run(false, startAt, r)
}

// FindNextMatch 返回与 match 参数相同的输入字符串中的下一个匹配项。
// 如果没有下一个匹配或给定 nil 匹配，则返回 nil。
func (re *Regexp) FindNextMatch(m *Match) (*Match, error) {
	if m == nil {
		return nil, nil
	}

	// If previous match was empty, advance by one before matching to prevent
	// infinite loop
	startAt := m.textpos
	if m.Length == 0 {
		if m.textpos == len(m.text) {
			return nil, nil
		}

		if re.RightToLeft() {
			startAt--
		} else {
			startAt++
		}
	}
	return re.run(false, startAt, m.text)
}

// MatchString 如果字符串与正则表达式匹配，则 MatchString 返回 true
// 如果超时则设置错误
func (re *Regexp) MatchString(s string) (bool, error) {
	m, err := re.run(true, -1, getRunes(s))
	if err != nil {
		return false, err
	}
	return m != nil, nil
}

func (re *Regexp) getRunesAndStart(s string, startAt int) ([]rune, int) {
	if startAt < 0 {
		if re.RightToLeft() {
			r := getRunes(s)
			return r, len(r)
		}
		return getRunes(s), 0
	}
	ret := make([]rune, len(s))
	i := 0
	runeIdx := -1
	for strIdx, r := range s {
		if strIdx == startAt {
			runeIdx = i
		}
		ret[i] = r
		i++
	}
	if startAt == len(s) {
		runeIdx = i
	}
	return ret[:i], runeIdx
}

func getRunes(s string) []rune {
	return []rune(s)
}

// MatchRunes 如果符文与正则表达式匹配，则 MatchRunes 返回 true
// 如果超时则设置错误
func (re *Regexp) MatchRunes(r []rune) (bool, error) {
	m, err := re.run(true, -1, r)
	if err != nil {
		return false, err
	}
	return m != nil, nil
}

// GetGroupNames 返回用于命名表达式中的捕获组的字符串集。
func (re *Regexp) GetGroupNames() []string {
	var result []string

	if re.capsList == nil {
		result = make([]string, re.capsize)

		for i := 0; i < len(result); i++ {
			result[i] = strconv.Itoa(i)
		}
	} else {
		result = make([]string, len(re.capsList))
		copy(result, re.capsList)
	}

	return result
}

// GetGroupNumbers 返回与组名称对应的整数组编号。
func (re *Regexp) GetGroupNumbers() []int {
	var result []int

	if re.caps == nil {
		result = make([]int, re.capsize)

		for i := 0; i < len(result); i++ {
			result[i] = i
		}
	} else {
		result = make([]int, len(re.caps))

		for k, v := range re.caps {
			result[v] = k
		}
	}

	return result
}

// GroupNameFromNumber 检索与组编号对应的组名称。
// 对于未知的组号，它将返回“”。自动未命名组
// 接收一个名称，该名称是与其编号等效的十进制字符串。
func (re *Regexp) GroupNameFromNumber(i int) string {
	if re.capsList == nil {
		if i >= 0 && i < re.capsize {
			return strconv.Itoa(i)
		}

		return ""
	}

	if re.caps != nil {
		var ok bool
		if i, ok = re.caps[i]; !ok {
			return ""
		}
	}

	if i >= 0 && i < len(re.capsList) {
		return re.capsList[i]
	}

	return ""
}

// GroupNumberFromName 返回与组名称对应的组编号。
// 如果名称不是可识别的组名称，则返回 -1。编号组
// 自动获取一个组名称，该名称是与其编号等效的十进制字符串。
func (re *Regexp) GroupNumberFromName(name string) int {
	// look up name if we have a hashtable of names
	if re.capNames != nil {
		if k, ok := re.capNames[name]; ok {
			return k
		}

		return -1
	}

	// convert to an int if it looks like a number
	result := 0
	for i := 0; i < len(name); i++ {
		ch := name[i]

		if ch > '9' || ch < '0' {
			return -1
		}

		result *= 10
		result += int(ch - '0')
	}

	// return int if it's in range
	if result >= 0 && result < re.capsize {
		return result
	}

	return -1
}
