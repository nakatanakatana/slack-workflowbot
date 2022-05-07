package slackworkflowbot

import (
	slackdev "github.com/nakatanakatana/slack"
	"github.com/slack-go/slack"
)

func CreateSlackClient(b BotToken) *slack.Client {
	return slack.New(string(b))
}

func CreateSlackDevClient(b BotToken) *slackdev.Client {
	return slackdev.New(string(b))
}
