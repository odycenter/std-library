package slack

import (
	"context"
	"encoding/json"
	"github.com/slack-go/slack"
	"strconv"
	"sync"
	"time"
)

var bots sync.Map //map[name]*Client

type Client struct {
	cli       *slack.Client
	channelID string
}

// Option Slcak机器人配置
type Option struct {
	AliseName string
	Token     string
	ChannelID string
	Debug     bool
}

func (o *Option) getAliseName() string {
	if o.AliseName == "" {
		return "default"
	}
	return o.AliseName
}

// Register 按配置创建Slcak机器人
func Register(opts ...*Option) {
	for _, opt := range opts {
		create(opt)
	}
}

func create(opt *Option) {
	bots.Store(opt.getAliseName(), &Client{slack.New(opt.Token, slack.OptionDebug(opt.Debug)), opt.ChannelID})
}

// Bot 按aliseName获取机器人
func Bot(aliseName ...string) *Client {
	name := "default"
	if len(aliseName) != 0 {
		name = aliseName[0]
	}
	client, ok := bots.Load(name)
	if !ok {
		return nil
	}
	return client.(*Client)
}

type Message struct {
	Pretext string
	Text    string
	Color   string
}

// WithChannel 指定非配置channelID
func (c *Client) WithChannel(channelID string) *Client {
	if channelID == "" {
		return c
	}
	c.channelID = channelID
	return c
}

// Send 发送消息
func (c *Client) Send(msg Message) error {
	_, _, err := c.cli.PostMessage(c.channelID, slack.MsgOptionAttachments(slack.Attachment{
		Color:   msg.Color,
		Pretext: msg.Pretext,
		Text:    msg.Text,
		Ts:      json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}))
	return err
}

// SendWithCtx 发送消息with context
func (c *Client) SendWithCtx(ctx context.Context, msg Message) error {
	_, _, err := c.cli.PostMessageContext(ctx, c.channelID, slack.MsgOptionAttachments(slack.Attachment{
		Color:   msg.Color,
		Pretext: msg.Pretext,
		Text:    msg.Text,
		Ts:      json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}))
	return err
}
