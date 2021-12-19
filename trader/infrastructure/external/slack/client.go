package slack

import (
	"github.com/slack-go/slack"
)

type Client struct {
	token     string
	channelId string
	client    *slack.Client
}

func NewClient(token, channelId string) *Client {
	client := slack.New(token)
	return &Client{
		token:     token,
		channelId: channelId,
		client:    client,
	}
}
