// Package tg 基于http的telegram信息发送
package tg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// SimpleSendMsg 发送tg信息
// url 这是目前telegram服务部署的地址
// msgType 消息类型，这个就是消息机器人的名字
// title，body 这是消息内容
func SimpleSendMsg(u, msgType, title, body string) error {
	if u == "" {
		return errors.New("telegram需要配置服务器地址")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("type", msgType)
	q.Add("message", fmt.Sprintf("%s\r\n%s", title, body))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	return err
}

// SendMsg 调用发送接口
// chantId 渠道id，类似发送者ID
// botToken 机器人唯一编号
func SendMsg(chatID, botToken, message string) error {
	if chatID == "" || botToken == "" {
		return errors.New("args error")
	}
	PostUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	resp, err := http.Post(PostUrl, "application/x-www-form-urlencoded", strings.NewReader("chat_id="+chatID+"&text="+message))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	return err
}
