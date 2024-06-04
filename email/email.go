package email

import (
	"fmt"
	"net/smtp"
	"sync"
)

var clients sync.Map //map[string]*Client

// Auth 身份验证参数，按需选填
// PlainAuth： Identity Username Password Host
// LoginAuth： Username Password
type Auth struct {
	Identity string
	Username string
	Password string
	Host     string
	Secret   string
}

// Option 邮件发送配置
type Option struct {
	AliasName  string
	Address    string     //smtp服务器地址
	AuthMethod AuthMethod //身份验证方式
	Auth       Auth
	auth       smtp.Auth
}

func (o *Option) getAliseName() string {
	if o.AliasName == "" {
		o.AliasName = "default"
	}
	return o.AliasName
}

func (o *Option) getAuth() smtp.Auth {
	switch o.AuthMethod {
	case MethodPlainAuth:
		return smtp.PlainAuth(o.Auth.Identity, o.Auth.Username, o.Auth.Password, o.Auth.Host)
	case MethodLoginAuth:
		return LoginAuth(o.Auth.Username, o.Auth.Password)
	default:
		return nil
	}
}

// Client 结构体
type Client struct {
	opt *Option
}

// New 创建Email客户端
func New(opt ...*Option) {
	for _, o := range opt {
		clients.Store(o.getAliseName(), &Client{opt: o})
	}
}

// Cli 使用 AliseName 获取发送Client
func Cli(aliseName ...string) *Client {
	if len(aliseName) == 0 {
		aliseName = append(aliseName, "default")
	}
	v, ok := clients.Load(aliseName[0])
	if ok {
		return v.(*Client)
	}
	return nil
}

// Address 返回client的地址
func (c *Client) Address() string {
	return c.opt.Auth.Username
}

// Send 发送
// msg格式需要遵循 RFC822-style
// For example:
// msg := []byte(
//
//	"To: recipient@example.net\r\n" +
//	"Subject: discount Gophers!\r\n" +
//	"\r\n" +
//	"This is the email body.\r\n"
//	)
//
// 可以使用 RFC822 来组织默认格式化的邮件
func (c *Client) Send(from string, to []string, msg []byte) error {
	return smtp.SendMail(c.opt.Address, c.opt.getAuth(), from, to, msg)
}

// RFC822 组装标准 RFC822 格式邮件内容
// 规定：
// to 必须是收件人的地址
func RFC822(from, to, subject, body string) []byte {
	return []byte(fmt.Sprintf(defaultRFC822, from, to, subject, body))
}
