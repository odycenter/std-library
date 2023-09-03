package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

type Conn struct {
	dial *websocket.Conn
}

// GetConn 获取原始conn
func (c *Conn) GetConn() *websocket.Conn {
	return c.dial
}

// Send 发送消息
func (c *Conn) Send(msg []byte) error {
	return c.dial.WriteMessage(websocket.BinaryMessage, msg)
}

// SendText 发送文本消息UTF-8编码
func (c *Conn) SendText(msg []byte) error {
	return c.dial.WriteMessage(websocket.TextMessage, msg)
}

// SendJson 发送json
func (c *Conn) SendJson(i any) error {
	return c.dial.WriteJSON(i)
}

// Close 关闭
func (c *Conn) Close() {
	_ = c.dial.Close()
}

// DialWithCtx 使用context创建websocket客户端
func DialWithCtx(ctx context.Context, scheme, host, path string, reqHeader *http.Header) (conn *Conn, resp *http.Response, err error) {
	var dial *websocket.Conn
	u := url.URL{Scheme: scheme, Host: host, Path: path}
	if reqHeader == nil {
		dial, resp, err = websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	} else {
		dial, resp, err = websocket.DefaultDialer.DialContext(ctx, u.String(), *reqHeader)
	}
	if err != nil {
		return
	}
	return &Conn{dial}, resp, err
}

// Dial 创建websocket客户端
func Dial(scheme, host, path string, reqHeader *http.Header) (conn *Conn, resp *http.Response, err error) {
	return DialWithCtx(context.Background(), scheme, host, path, reqHeader)
}
