package slack

import (
	"strings"

	"github.com/slack-go/slack"
)

type SlackEmoji string

const (
	SlackEmojiCoin      SlackEmoji = ":coin:"
	SlackEmojiDizzyFace SlackEmoji = ":dizzy_face:"
)

func BuildTextMessage(lines ...string) string {
	msg := strings.Join(lines, "\n")
	return msg
}

func PostTextMessage(token, channelID, msg string) error {
	slackBot := slack.New(token)
	_, _, err := slackBot.PostMessage(channelID, slack.MsgOptionText(msg, true))
	return err
}
