package jwt_test

import (
	"fmt"
	"std-library/jwt"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	tk := jwt.Gen(map[string]any{
		"Browser":          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
		"CreateTime":       1685501436,
		"HId":              0,
		"Id":               1113,
		"Ip":               "54.255.114.140",
		"UseNewPermission": true,
		"exp":              1717037436,
	}, "TokenGenSalt", time.Minute*10)
	fmt.Println(tk)
	payload := jwt.Parse(tk, "TokenGenSalt")
	fmt.Printf("%+v\n", payload)
}
