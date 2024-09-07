package regexps

import (
	"bytes"
	"errors"
	"github.com/odycenter/std-library/regexps/syntax"
)

const (
	replaceSpecials     = 4
	replaceLeftPortion  = -1
	replaceRightPortion = -2
	replaceLastGroup    = -3
	replaceWholeString  = -4
)

// MatchEvaluator 是一个接受匹配并返回要使用的替换字符串函数
type MatchEvaluator func(Match) string

// 下面出现了三个非常相似的算法：replace (pattern),
// 替换（评估器），并拆分。

// Replace 将字符串中所有出现的正则表达式替换为
// 替换模式。
//
// 请注意，没有匹配的特殊情况是单独处理的：
// 如果没有匹配项，则输入字符串原样返回。
// 从右到左的情况被分开，因为 StringBuilder
// 不能很好地直接处理从右到左的字符串构建。
func replace(regex *Regexp, data *syntax.ReplacerData, evaluator MatchEvaluator, input string, startAt, count int) (string, error) {
	if count < -1 {
		return "", errors.New("Count too small")
	}
	if count == 0 {
		return "", nil
	}

	m, err := regex.FindStringMatchStartingAt(input, startAt)

	if err != nil {
		return "", err
	}
	if m == nil {
		return input, nil
	}

	buf := &bytes.Buffer{}
	text := m.text

	if !regex.RightToLeft() {
		prevat := 0
		for m != nil {
			if m.Index != prevat {
				buf.WriteString(string(text[prevat:m.Index]))
			}
			prevat = m.Index + m.Length
			if evaluator == nil {
				replacementImpl(data, buf, m)
			} else {
				buf.WriteString(evaluator(*m))
			}

			count--
			if count == 0 {
				break
			}
			m, err = regex.FindNextMatch(m)
			if err != nil {
				return "", nil
			}
		}

		if prevat < len(text) {
			buf.WriteString(string(text[prevat:]))
		}
	} else {
		prevat := len(text)
		var al []string

		for m != nil {
			if m.Index+m.Length != prevat {
				al = append(al, string(text[m.Index+m.Length:prevat]))
			}
			prevat = m.Index
			if evaluator == nil {
				replacementImplRTL(data, &al, m)
			} else {
				al = append(al, evaluator(*m))
			}

			count--
			if count == 0 {
				break
			}
			m, err = regex.FindNextMatch(m)
			if err != nil {
				return "", nil
			}
		}

		if prevat > 0 {
			buf.WriteString(string(text[:prevat]))
		}

		for i := len(al) - 1; i >= 0; i-- {
			buf.WriteString(al[i])
		}
	}

	return buf.String(), nil
}

// Given a Match, emits into the StringBuilder the evaluated
// substitution pattern.
func replacementImpl(data *syntax.ReplacerData, buf *bytes.Buffer, m *Match) {
	for _, r := range data.Rules {

		if r >= 0 { // string lookup
			buf.WriteString(data.Strings[r])
		} else if r < -replaceSpecials { // group lookup
			m.groupValueAppendToBuf(-replaceSpecials-1-r, buf)
		} else {
			switch -replaceSpecials - 1 - r { // special insertion patterns
			case replaceLeftPortion:
				for i := 0; i < m.Index; i++ {
					buf.WriteRune(m.text[i])
				}
			case replaceRightPortion:
				for i := m.Index + m.Length; i < len(m.text); i++ {
					buf.WriteRune(m.text[i])
				}
			case replaceLastGroup:
				m.groupValueAppendToBuf(m.GroupCount()-1, buf)
			case replaceWholeString:
				for i := 0; i < len(m.text); i++ {
					buf.WriteRune(m.text[i])
				}
			}
		}
	}
}

func replacementImplRTL(data *syntax.ReplacerData, al *[]string, m *Match) {
	l := *al
	buf := &bytes.Buffer{}

	for _, r := range data.Rules {
		buf.Reset()
		if r >= 0 { // string lookup
			l = append(l, data.Strings[r])
		} else if r < -replaceSpecials { // group lookup
			m.groupValueAppendToBuf(-replaceSpecials-1-r, buf)
			l = append(l, buf.String())
		} else {
			switch -replaceSpecials - 1 - r { // special insertion patterns
			case replaceLeftPortion:
				for i := 0; i < m.Index; i++ {
					buf.WriteRune(m.text[i])
				}
			case replaceRightPortion:
				for i := m.Index + m.Length; i < len(m.text); i++ {
					buf.WriteRune(m.text[i])
				}
			case replaceLastGroup:
				m.groupValueAppendToBuf(m.GroupCount()-1, buf)
			case replaceWholeString:
				for i := 0; i < len(m.text); i++ {
					buf.WriteRune(m.text[i])
				}
			}
			l = append(l, buf.String())
		}
	}

	*al = l
}
