package email

type AuthMethod int

// 身份验证方式
const (
	MethodPlainAuth = AuthMethod(iota)
	MethodLoginAuth
)

const defaultRFC822 = `From: %s
To: %s
Subject: %s

%s`
