package tg_test

import (
	"testing"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/odycenter/std-library/sdk/tg"
)

func TestTG(t *testing.T) {
	opts := []*tg.Opt{{
		AliasName: "",
		Token:     "87b0gnvj0761t8",
	}, {
		AliasName: "bot1",
		Token:     "87b0gnvj0761t8",
	},
	}
	_ = tg.Register(opts...)
	_ = tg.Bot().SendMsg("测试信息测试信息测试信息", 1234, 0, nil)
	_ = tg.Bot("bot1").SendMsg("测试信息测试信息测试信息", 1234, 0, func(msg tgBotApi.MessageConfig) {
		msg.ParseMode = "Markdown"
	})
}
