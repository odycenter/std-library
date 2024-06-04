package email_test

import (
	"log"
	"std-library/email"
	"testing"
)

func TestEMail(t *testing.T) {
	email.New(&email.Option{
		Address:    "smtp.xxx.xxx:587",
		AuthMethod: email.MethodPlainAuth,
		Auth:       email.Auth{Identity: "", Username: "xxx@xxx.com", Password: "xxxxx", Host: "smtp.xxx.com"},
	})
	msg := email.RFC822("xxx@xxx.com", "yyy@yyy.com", "验证码", "您的验证码为：\n123456")

	err := email.Cli().Send("no-reply", []string{"css123_123@hotmail.com"}, msg)
	if err != nil {
		log.Println(err)
		return
	}
}
