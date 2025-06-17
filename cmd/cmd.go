package main

import (
	"log"

	env "github.com/pal-paul/go-libraries/pkg/env"
	slack "github.com/pal-paul/go-libraries/pkg/slack"
)

type Environment struct {
	Input struct {
		Title string `env:"INPUT_TITLE,required=true"`
		Text  string `env:"INPUT_TEXT,required=true"`
	}
	Slack struct {
		Token   string `env:"INPUT_SLACK_TOKEN,required=true"`
		Channel string `env:"INPUT_SLACK_CHANNEL,required=true"`
	}
}

var (
	envVar      Environment
	slackClient slack.ISlack
)

// Initialize environment variables and Slack client
func initializeApp() error {
	_, err := env.Unmarshal(&envVar)
	if err != nil {
		return err
	}

	slackClient = slack.New(
		slack.WithToken(envVar.Slack.Token),
	)
	return nil
}

func main() {
	if err := initializeApp(); err != nil {
		log.Fatal(err)
	}

	message := SlackMessageBuilder(envVar.Input.Title, envVar.Input.Text, envVar.Slack.Channel)
	_, err := slackClient.AddFormattedMessage(envVar.Slack.Channel, message)
	if err != nil {
		log.Fatalf("error while sending message to slack: %v", err)
	}
}

func SlackMessageBuilder(title string, text string, channel string) slack.Message {
	message := slack.Message{
		Channel: channel,
	}

	message.Blocks = append(message.Blocks, slack.Block{
		Type: slack.HeaderBlock,
		Text: &slack.Text{
			Type: slack.PlainText,
			Text: title,
		},
	})
	message.Blocks = append(message.Blocks, slack.Block{
		Type: slack.SectionBlock,
		Text: &slack.Text{
			Type: slack.Mrkdwn,
			Text: text,
		},
	})
	return message
}
