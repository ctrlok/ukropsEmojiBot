package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

type ToBestEvent struct {
	ReactionEvent
}

func (r ToBestEvent) Process() (events.APIGatewayProxyResponse, error) {
	log.Debug("Processing slack reaction event at first")
	_, err := r.ReactionEvent.Process()
	if err != nil {
		log.Error(err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	err = r.process()
	if err != nil {
		log.Error(err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	log.Debug("So... We got a single message from a history. Let's process it.")
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func (r ToBestEvent) process() error {
	log.Debug("Let's find a reaction we need...")
	reaction, err := r.getReaction()
	if err != nil {
		return err
	}
	// TODO: rewrite to storing this reaction and message into DB
	if reaction.Count != 1 {
		log.Debug("More than single reaction on the event")
		return nil
	}
	log.Debugf("We need to find a message in history, to copy a text from.")
	textMessage, err := r.getRelatedMessageText()
	if err != nil {
		return err
	}
	log.Debug("Lets create a permalink to the message we will quote")
	permalink, err := r.getPermalinkToMessage()
	if err != nil {
		return err
	}
	log.Debugf("Permalink created: %s", permalink)

	slackPost := r.createPost(textMessage, permalink)

	log.Debug("Looks like we did all we can. Let's post this message to a slack.")
	_, _, err = client.PostMessage(config.BestChannelId, slackPost...)
	if err != nil {
		log.Errorf("Problem with senging a message: %s.", err)
	}
	return nil
}

func (r ToBestEvent) getReaction() (reaction slack.ItemReaction, err error) {
	reactions, err := client.GetReactions(slack.ItemRef{
		Channel:   r.ReactionEvent.Event.Item.Channel,
		Timestamp: r.ReactionEvent.Event.Item.Timestamp,
	}, slack.NewGetReactionsParameters())
	if err != nil {
		return
	}
	log.WithFields(log.Fields{"reactions": reactions}).Debugf("We got reactions for the message")
	for _, itemReaction := range reactions {
		if itemReaction.Name == config.BestEmojiName {
			reaction = itemReaction
		}
	}
	return
}

func (r ToBestEvent) getPermalinkToMessage() (string, error) {
	return client.GetPermalink(&slack.PermalinkParameters{
		Channel: r.ReactionEvent.Event.Item.Channel,
		Ts:      r.ReactionEvent.Event.Item.Timestamp,
	})
}

func (r ToBestEvent) getRelatedMessageText() (string, error) {
	channelHistory, err := legacyClient.GetChannelHistory(r.ReactionEvent.Event.Item.Channel, slack.HistoryParameters{
		Latest:    r.ReactionEvent.Event.Item.Timestamp,
		Inclusive: true,
		Count:     1,
	})
	if err != nil {
		return "", err
	}
	if len(channelHistory.Messages) != 1 {
		err = fmt.Errorf("not a single message returned from a search: %v", channelHistory)
		return "", err
	}
	log.Debug("So... We got a single message from a history. Let's process it.")
	return channelHistory.Messages[0].Text, nil
}

func (r ToBestEvent) createPost(text string, link string) []slack.MsgOption {
	log.Debug("Let's create a nice looking post message")
	headerString := fmt.Sprintf("<@%s> написал, _(а <@%s> добавил в лучшее)_: ", r.ReactionEvent.Event.ItemUser, r.ReactionEvent.Event.User)
	headerTextObject := slack.NewTextBlockObject("mrkdwn", headerString, false, false)
	headerSection := slack.NewSectionBlock(headerTextObject, nil, nil)

	bodyString := fmt.Sprintf("%s\n\n <%s|link to message>", text, link)
	bodyText := slack.NewTextBlockObject("mrkdwn", bodyString, false, false)
	bodySection := slack.NewSectionBlock(bodyText, nil, nil)
	log.WithFields(log.Fields{"header": headerSection, "body": bodySection}).Debug("Message created and ready.")
	return []slack.MsgOption{slack.MsgOptionBlocks(headerSection), slack.MsgOptionAttachments(slack.Attachment{Blocks: []slack.Block{bodySection}})}
}
