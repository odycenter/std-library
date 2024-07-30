// Package jwt JWT令牌管理方法封装
// Deprecated: JWT Function not safe for payload
package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			fmt.Println("That's Not A Token:", token, err)
			return
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			fmt.Println("Token Has Expired:", token, err)
			return
		default:
			fmt.Println("Invalid Token:", err)
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
