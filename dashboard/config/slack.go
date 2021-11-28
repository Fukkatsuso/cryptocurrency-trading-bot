package config

import "os"

var (
	SlackBotToken  string
	SlackChannelID string
)

func init() {
	SlackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	SlackChannelID = os.Getenv("SLACK_CHANNEL_ID")
}
