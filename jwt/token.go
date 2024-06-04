// Package jwt JWT令牌管理方法封装
// Deprecated: JWT Function not safe for payload
package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Gen 生成Token
// Deprecated: JWT Function not safe for payload
func Gen(claims map[string]any, salt string, exp time.Duration) string {
	if exp != 0 {
		claims["exp"] = time.Now().Add(exp).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	signed, err := token.SignedString([]byte(salt))
	if err != nil {
		fmt.Println("Token General Failed:", err)
	}
	return signed
}

// Parse 解析Token
// Deprecated: JWT Function not safe for payload
func Parse(otk, salt string) (claims map[string]any) {
	token, err := jwt.Parse(otk, func(token *jwt.Token) (any, error) {
		return []byte(salt), nil
	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("That's Not A Token:", token, err)
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				fmt.Println("Token Has Expired:", token, err)
				return
			} else {
				fmt.Println("Invalid Token:", token, err)
				return
			}
		} else {
			fmt.Println("Token Parse Failed:", token, err)
			return
		}
	}
	if !token.Valid {
		fmt.Println("Invalid Token:", err)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Parse Token Format Error:", err)
		return
	}
	return
}

// OriginParse 原生解析Token
// Deprecated: JWT Function not safe for payload
func OriginParse(otk string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return jwt.Parse(otk, keyFunc)
}
