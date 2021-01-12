package slack

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

var slackToken = os.Getenv("SLACK_TOKEN")

func SendMessageToChannel(channel, message string) error {
	api := slack.New(slackToken, slack.OptionDebug(true))
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues

	channelID, timestamp, err := api.PostMessage(
		channel,
		slack.MsgOptionText(message, false),
		//slack.MsgOptionAttachments(attachment),
		//slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)

	fmt.Println(channelID, timestamp)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	return err
}
