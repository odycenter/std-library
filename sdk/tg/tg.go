// Package tg Telegram 机器人SDK封装
package tg

import (
	"fmt"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

var tgBots sync.Map //map[name]*tgBotApi.BotAPI

type BotApi struct {
	bot *tgBotApi.BotAPI
}

// Opt 机器人配置
type Opt struct {
	AliasName string
	Token     string
}

func (o *Opt) getAliasName() string {
	if o.AliasName == "" {
		return "default"
	}
	return o.AliasName
}

// Register 注册新的tg机器人
func Register(opts ...*Opt) error {
	for _, opt := range opts {
		bot, err := tgBotApi.NewBotAPI(opt.Token)
		if err != nil {
			return fmt.Errorf("tg bot register failed:%s", err.Error())
		}
		tgBots.Store(opt.getAliasName(), &BotApi{bot})
	}
	return nil
}

// Bot 获取一个指定的tgBot
// 不指定name则默认default
func Bot(name ...string) *BotApi {
	n := "default"
	if len(name) != 0 {
		n = name[0]
	}
	v, ok := tgBots.Load(n)
	if !ok {
		return nil
	}
	botApi, ok := v.(*BotApi)
	if !ok {
		return nil
	}
	return botApi
}

// SendMsg 发送文本消息
// text 文本内容
// chatID 频道ID
// replyMsgID 返回信息ID，如果不需要自己维护可传入0
func (b *BotApi) SendMsg(text string, chatID int64, replyMsgID int, fn func(msg tgBotApi.MessageConfig)) error {
	msg := tgBotApi.NewMessage(chatID, text)
	if replyMsgID != 0 {
		msg.ReplyToMessageID = replyMsgID
	}
	if fn != nil {
		fn(msg)
	}
	_, err := b.bot.Send(msg)
	return err
}

// SendBytes 发送二进制文本
// text 文本内容
// chatID 频道ID
// replyMsgID 返回信息ID，如果不需要自己维护可传入0
func (b *BotApi) SendBytes(docName string, doc []byte, chatID int64, replyMsgID int, fn func(msg tgBotApi.DocumentConfig)) error {
	msg := tgBotApi.NewDocument(chatID, tgBotApi.FileBytes{
		Name:  docName,
		Bytes: doc,
	})
	if replyMsgID != 0 {
		msg.ReplyToMessageID = replyMsgID
	}
	if fn != nil {
		fn(msg)
	}
	_, err := b.bot.Send(msg)
	return err
}

// SendFile 发送文件
// filePath 文件路径
// chatID 频道ID
// replyMsgID 返回信息ID，如果不需要自己维护可传入0
// caption 标题内容
func (b *BotApi) SendFile(filePath string, chatID int64, replyMsgID int, fn func(msg tgBotApi.DocumentConfig), caption ...string) error {
	msg := tgBotApi.NewDocument(chatID, tgBotApi.FilePath(filePath))
	if replyMsgID != 0 {
		msg.ReplyToMessageID = replyMsgID
	}
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	if fn != nil {
		fn(msg)
	}
	_, err := b.bot.Send(msg)
	return err
}

// SendImg 发送图片
// filePath 文件路径
// chatID 频道ID
// replyMsgID 返回信息ID，如果不需要自己维护可传入0
// caption 标题内容
func (b *BotApi) SendImg(filePath string, chatID int64, replyMsgID int, fn func(msg tgBotApi.PhotoConfig), caption ...string) error {
	msg := tgBotApi.NewPhoto(chatID, tgBotApi.FilePath(filePath))
	if replyMsgID != 0 {
		msg.ReplyToMessageID = replyMsgID
	}
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	if fn != nil {
		fn(msg)
	}
	_, err := b.bot.Send(msg)
	return err
}
