package fuzzy

import (
	"math"
	"sort"
	"strings"
	"unicode/utf8"
)

// Ratio 根据两个 unicode 字符串的 Levenshtein 编辑距离计算其接近程度的分数。
// 返回一个整数分数[0,100]，分数越高表示字符串越接近。
func Ratio(s1, s2 string) int {
	return int(round(100 * floatRatio([]rune(s1), []rune(s2))))
}

// PartialRatio 计算一个字符串与另一个字符串中最相似的子字符串的接近程度的分数。
// 参数的顺序并不重要。
// 返回一个整数分数[0,100]，分数越高表示字符串和子串越接近。
func PartialRatio(s1, s2 string) int {
	shorter, longer := []rune(s1), []rune(s2)
	if len(shorter) > len(longer) {
		longer, shorter = shorter, longer
	}
	matchingBlocks := getMatchingBlocks(shorter, longer)

	bestScore := 0.0
	for _, block := range matchingBlocks {
		longStart := block.dpos - block.spos
		if longStart < 0 {
			longStart = 0
		}
		longEnd := longStart + len(shorter)
		if longEnd > len(longer) {
			longEnd = len(longer)
		}
		longSubStr := longer[longStart:longEnd]

		r := floatRatio(shorter, longSubStr)
		if r > .995 {
			return 100
		} else if r > bestScore {
			bestScore = r
		}
	}

	return int(round(100 * bestScore))
}

func floatRatio(chrs1, chrs2 []rune) float64 {
	lenSum := len(chrs1) + len(chrs2)
	if lenSum == 0 {
		return 0.0
	}
	editDistance := optimizedEditDistance(chrs1, chrs2, 1)
	return float64(lenSum-editDistance) / float64(lenSum)
}

// QRatio 计算的分数与 Ratio 类似，不同之处在于两个字符串都被修剪、清除了非 ASCII 字符并进行了大小写标准化。
func QRatio(s1, s2 string) int {
	return quickRatioHelper(s1, s2, true)
}

// UQRatio UQRatio 计算的分数与 Ratio 类似，只不过两个字符串都被修剪并进行了大小写标准化。
func UQRatio(s1, s2 string) int {
	return quickRatioHelper(s1, s2, false)
}

func quickRatioHelper(s1, s2 string, asciiOnly bool) int {
	c1 := Cleanse(s1, asciiOnly)
	c2 := Cleanse(s2, asciiOnly)

	if len(c1) == 0 || len(c2) == 0 {
		return 0
	}
	return Ratio(c1, c2)
}

// WRatio 通过以下步骤计算分数：
//  1. 清理两个字符串，删除非 ASCII 字符。
//  2. 以比率作为基线分数。
//  3. 运行一些启发式方法来确定是否应采用部分比率。
//  4. 如果确定需要部分比率，
//     则计算 PartialRatio、PartialTokenSetRatio 和 PartialTokenSortRatio。
//     否则，计算 TokenSortRatio 和 TokenSetRatio。
//  5. 返回所有计算比率的最大值。
func WRatio(s1, s2 string) int {
	return weightedRatioHelper(s1, s2, true)
}

// UWRatio UWRatio 计算的分数与 WRatio 类似，但允许使用非 ASCII 字符。
func UWRatio(s1, s2 string) int {
	return weightedRatioHelper(s1, s2, false)
}

func weightedRatioHelper(s1, s2 string, asciiOnly bool) int {
	c1 := Cleanse(s1, asciiOnly)
	c2 := Cleanse(s2, asciiOnly)

	if len(c1) == 0 || len(c2) == 0 {
		return 0
	}

	unbaseScale := .95
	partialScale := .9
	baseScore := float64(Ratio(c1, c2))
	lengthRatio := float64(utf8.RuneCountInString(c1)) / float64(utf8.RuneCountInString(c2))
	if lengthRatio < 1 {
		lengthRatio = 1 / lengthRatio
	}

	tryPartial := true
	if lengthRatio < 1.5 {
		tryPartial = false
	}

	if lengthRatio > 8 {
		partialScale = .6
	}

	if tryPartial {
		partialScore := float64(PartialRatio(c1, c2)) * partialScale
		tokenSortScore := float64(PartialTokenSortRatio(c1, c2, asciiOnly, false)) *
			unbaseScale * partialScale
		tokenSetScore := float64(PartialTokenSetRatio(c1, c2, asciiOnly, false)) *
			unbaseScale * partialScale
		return int(round(max(baseScore, partialScore, tokenSortScore, tokenSetScore)))
	}
	tokenSortScore := float64(TokenSortRatio(c1, c2, asciiOnly, false)) * unbaseScale
	tokenSetScore := float64(TokenSetRatio(c1, c2, asciiOnly, false)) * unbaseScale
	return int(round(max(baseScore, tokenSortScore, tokenSetScore)))
}

func max(args ...float64) float64 {
	maxVal := args[0]
	for _, arg := range args {
		if arg > maxVal {
			maxVal = arg
		}
	}
	return maxVal
}

// TokenSortRatio 计算与 Ratio 类似的分数，除了在比较之前对标记进行排序和（可选）清理
func TokenSortRatio(s1, s2 string, opts ...bool) int {
	return tokenSortRatioHelper(s1, s2, false, opts...)
}

// PartialTokenSortRatio 计算类似于 PartialRatio 的分数，只不过在比较之前对标记进行排序和（可选）清理。
func PartialTokenSortRatio(s1, s2 string, opts ...bool) int {
	return tokenSortRatioHelper(s1, s2, true, opts...)
}

func tokenSortRatioHelper(s1, s2 string, partial bool, opts ...bool) int {
	asciiOnly, cleanse := false, false
	for i, val := range opts {
		switch i {
		case 0:
			asciiOnly = val
		case 1:
			cleanse = val
		}
	}

	sorted1 := tokenSort(s1, asciiOnly, cleanse)
	sorted2 := tokenSort(s2, asciiOnly, cleanse)

	if partial {
		return PartialRatio(sorted1, sorted2)
	}
	return Ratio(sorted1, sorted2)
}

func tokenSort(s string, asciiOnly, cleanse bool) string {
	if cleanse {
		s = Cleanse(s, asciiOnly)
	} else if asciiOnly {
		s = ASCIIOnly(s)
	}

	tokens := strings.Fields(s)
	sort.Strings(tokens)
	return strings.Join(tokens, " ")
}

// TokenSetRatio 从每个输入字符串中提取标记，将它们添加到一个集合中，
// 构造<已排序的交集><已排序的余数>形式的字符串，
// 获取这两个字符串的比率，
// 并返回最大值。
func TokenSetRatio(s1, s2 string, opts ...bool) int {
	return tokenSetRatioHelper(s1, s2, false, opts...)
}

// PartialTokenSetRatio 从每个输入字符串中提取标记，将它们添加到一个集合中，
// 构造两个 <排序交集> <排序余数> 形式的字符串，
// 获取这两个字符串的部分比率，
// 并返回最大值。
func PartialTokenSetRatio(s1, s2 string, opts ...bool) int {
	return tokenSetRatioHelper(s1, s2, true, opts...)
}

func tokenSetRatioHelper(s1, s2 string, partial bool, opts ...bool) int {
	asciiOnly, cleanse := false, false
	for i, val := range opts {
		switch i {
		case 0:
			asciiOnly = val
		case 1:
			cleanse = val
		}
	}

	if cleanse {
		s1 = Cleanse(s1, asciiOnly)
		s2 = Cleanse(s2, asciiOnly)
	} else if asciiOnly {
		s1 = ASCIIOnly(s1)
		s2 = ASCIIOnly(s2)
	}

	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	set1 := NewStringSet(strings.Fields(s1))
	set2 := NewStringSet(strings.Fields(s2))
	intersection := set1.Intersect(set2).ToSlice()
	diff1to2 := set1.Difference(set2).ToSlice()
	diff2to1 := set2.Difference(set1).ToSlice()

	sort.Strings(intersection)
	sort.Strings(diff1to2)
	sort.Strings(diff2to1)

	sortedIntersect := strings.TrimSpace(strings.Join(intersection, " "))
	combined1to2 := strings.TrimSpace(sortedIntersect + " " + strings.Join(diff1to2, " "))
	combined2to1 := strings.TrimSpace(sortedIntersect + " " + strings.Join(diff2to1, " "))

	var ratioFunction func(string, string) int
	if partial {
		ratioFunction = PartialRatio
	} else {
		ratioFunction = Ratio
	}

	score := ratioFunction(sortedIntersect, combined1to2)
	if alt1 := ratioFunction(sortedIntersect, combined2to1); alt1 > score {
		score = alt1
	}
	if alt2 := ratioFunction(combined1to2, combined2to1); alt2 > score {
		score = alt2
	}

	return score
}

func round(x float64) float64 {
	if x < 0 {
		return math.Ceil(x - 0.5)
	}
	return math.Floor(x + 0.5)
}
