package regexps_test

import (
	"log"
	"std-library/regexps"
	"testing"
)

func TestRegexps(t *testing.T) {
	reg, err := regexps.Compile(`^(?=(.*[a-zA-Z]){2})(?=.*\d)(?=.*[a-zA-Z0-9])(?=([a-zA-Z0-9])\1*(?!\1))\w{6,12}$`, regexps.None)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(reg.MatchString("abcdefg1234"))
	log.Println(reg.MatchString("abcdefghijk"))
	log.Println(reg.MatchString("abc"))
	log.Println(reg.MatchString("abcdefghijklmnopqr123"))
	log.Println(reg.MatchString("aaa1"))
	log.Println(reg.MatchString("aaabbb11 "))
	log.Println(reg.MatchString("aaabbbCC11 "))
	log.Println(reg.MatchString("aaabbbCC11"))
	log.Println(reg.MatchString("222aaabbbCC11"))
}
