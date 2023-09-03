package websocket

import (
	"context"
	"github.com/olahol/melody"
	"net/http"
	"net/url"
	"sync"
)

var clients sync.Map //map[AliasName]*Client

// Client websocket管理结构
type Client struct {
	AliasName string         //别名
	m         *melody.Melody //ws对象
	opt       *Option        //配置
	started   chan struct{}
}

// Option websocket配置
type Option struct {
	AliasName       string                                               //别名
	IP              string                                               //监听IP
	Port            string                                               //监听端口
	Pattern         string                                               //链接路由
	OnConnect       func(session *melody.Session)                        `json:"-"`
	OnDisconnect    func(session *melody.Session)                        `json:"-"`
	OnClose         func(session *melody.Session, i int, s string) error `json:"-"`
	OnMessage       func(session *melody.Session, b []byte)              `json:"-"`
	OnMessageBinary func(session *melody.Session, b []byte)              `json:"-"`
}

func (opt *Option) getAliasName() string {
	if opt.AliasName == "" {
		return "default"
	}
	return opt.AliasName
}

func (opt *Option) getIP() string {
	if opt.IP == "" {
		return "localhost"
	}
	return opt.IP
}

func (opt *Option) getPort() string {
	if opt.AliasName == "" {
		return ":8765"
	}
	return ":" + opt.Port
}

func (opt *Option) getPattern() string {
	if opt.Pattern == "" {
		return "/"
	}
	return opt.Pattern
}

// WithOnConnect 响应连接回调方法
func (opt *Option) WithOnConnect(fn func(session *melody.Session)) {
	opt.OnConnect = fn
}

// WithOnDisconnect 响应断开链接回调方法
func (opt *Option) WithOnDisconnect(fn func(session *melody.Session)) {
	opt.OnDisconnect = fn
}

// WithOnClose 响应关闭连接回调方法
func (opt *Option) WithOnClose(fn func(session *melody.Session, i int, s string) error) {
	opt.OnClose = fn
}

// WithOnMessage 响应message处理回调方法
func (opt *Option) WithOnMessage(fn func(session *melody.Session, b []byte)) {
	opt.OnMessage = fn
}

// WithOnMessageBinary 响应binary message处理回调方法
func (opt *Option) WithOnMessageBinary(fn func(session *melody.Session, b []byte)) {
	opt.OnMessageBinary = fn
}

// New 创建websocket
func New(opt *Option) {
	clients.Store(opt.getAliasName(), &Client{
		AliasName: opt.getAliasName(),
		m:         melody.New(),
		opt:       opt,
		started:   make(chan struct{}, 1),
	})
}

// Cli 获取client
func Cli(aliasName ...string) *Client {
	name := "default"
	if len(aliasName) != 0 {
		name = aliasName[0]
	}
	cli, ok := clients.Load(name)
	if ok {
		return cli.(*Client)
	}
	return nil
}

func setKeys(r *http.Request) (m map[string]interface{}) {
	m = make(map[string]interface{})
	m["request"] = r
	return
}

// GetRequest 获取连接时的request
func GetRequest(session *melody.Session) *http.Request {
	v, ok := session.Get("request")
	if ok {
		return v.(*http.Request)
	}
	return nil
}

// GetURL 获取连接时的request的URL
func GetURL(session *melody.Session) *url.URL {
	v, ok := session.Get("request")
	if ok {
		return v.(*http.Request).URL
	}
	return nil
}

// ListenWithServer 开启监听并开启http server
// pattern 路由
func (cli *Client) ListenWithServer(ctx context.Context) error {
	http.HandleFunc(cli.opt.getPattern(), func(rw http.ResponseWriter, r *http.Request) {
		err := cli.m.HandleRequestWithKeys(rw, r, setKeys(r))
		if err != nil {
			context.WithValue(ctx, "ListenWithServer", err)
			return
		}
	})
	cli.setHandle()
	err := http.ListenAndServe(cli.opt.getIP()+cli.opt.getPort(), nil)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

// Listen 开启监听
// pattern 路由
func (cli *Client) Listen(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	err := cli.m.HandleRequestWithKeys(rw, r, setKeys(r))
	if err != nil {
		context.WithValue(ctx, "ListenWithServer", err)
		return
	}
	cli.setHandle()
}

func (cli *Client) setHandle() {
	if cli.opt.OnConnect != nil {
		cli.m.HandleConnect(cli.opt.OnConnect)
	}
	if cli.opt.OnDisconnect != nil {
		cli.m.HandleDisconnect(cli.opt.OnDisconnect)
	}
	if cli.opt.OnClose != nil {
		cli.m.HandleClose(cli.opt.OnClose)
	}
	if cli.opt.OnMessage != nil {
		cli.m.HandleMessage(cli.opt.OnMessage)
	}
	if cli.opt.OnMessageBinary != nil {
		cli.m.HandleMessageBinary(cli.opt.OnMessage)
	}
}

// Broadcast 向所有会话广播文本消息
func (cli *Client) Broadcast(msg []byte) error {
	return cli.m.Broadcast(msg)
}

// BroadcastBinary 向所有会话广播二进制消息。
func (cli *Client) BroadcastBinary(msg []byte) error {
	return cli.m.BroadcastBinary(msg)
}

// BroadcastOthers 向除 session 之外的所有会话广播文本消息。
func (cli *Client) BroadcastOthers(msg []byte, session *melody.Session) error {
	return cli.m.BroadcastOthers(msg, session)
}

// BroadcastFilter 向所有 fn 返回 true 的会话广播一条文本消息。
func (cli *Client) BroadcastFilter(msg []byte, fn func(session *melody.Session) bool) error {
	return cli.m.BroadcastFilter(msg, fn)
}

// BroadcastBinaryOthers 向除 session 之外的所有会话广播二进制消息。
func (cli *Client) BroadcastBinaryOthers(msg []byte, session *melody.Session) error {
	return cli.m.BroadcastBinaryOthers(msg, session)
}

// BroadcastBinaryFilter 向向所有 fn 返回 true 的会话广播二进制消息。
func (cli *Client) BroadcastBinaryFilter(msg []byte, session *melody.Session) error {
	return cli.m.BroadcastBinaryOthers(msg, session)
}

// BroadcastMultiple 将文本消息广播到会话片中给定的多个会话。
func (cli *Client) BroadcastMultiple(msg []byte, sessions []*melody.Session) error {
	return cli.m.BroadcastMultiple(msg, sessions)
}
