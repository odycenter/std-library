// Package tsGoogleAuth Google身份认证（google码）
package tsGoogleAuth

// GetQrcode 获取二维码
func GetQrcode(user string) (string, string, string) {
	// 秘钥
	secret := NewGoogleAuth().GetSecret()
	// 用户名
	qrCode := NewGoogleAuth().GetQrcode(user, secret)
	// 第三方二维码地址
	qrCodeUrl := NewGoogleAuth().GetQrcodeUrl(user, secret)
	return secret, qrCode, qrCodeUrl
}

// GetCode 动态码(每隔30s会动态生成一个6位数的数字)
func GetCode(secret string) (string, error) {
	return NewGoogleAuth().GetCode(secret)
}

// VerifyCode 验证动态吗
func VerifyCode(secret, code string) (bool, string, error) {
	return NewGoogleAuth().VerifyCode(secret, code)
}
