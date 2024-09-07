package slack_test

import (
	"github.com/odycenter/std-library/sdk/slack"
	"testing"
)

func TestNew(t *testing.T) {
	slack.Register(&slack.Option{
		AliseName: "",
		Token:     "",
		ChannelID: "",
		Debug:     false,
	})
	err := slack.Bot().WithChannel("").Send(slack.Message{
		Pretext: "ABC",
		Text:    "DEF",
		Color:   "#4AF030",
	})
	if err != nil {
		return
	}
}
