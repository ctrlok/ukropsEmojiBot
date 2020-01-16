package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

func main() {
	lambda.Start(HandleEventTest)
}

func HandleEventTest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}

	log.WithFields(log.Fields{"body": request.Body}).Debug("Start processing a body")
	lambdaEvent, err := parseBody(request.Body)
	if err != nil {
		log.Error(err)
		return response, err
	}
	log.WithFields(log.Fields{"recievedEvent": lambdaEvent}).Debug("Finish processing a body")

	log.Debug("Check is this request is a challenge?")
	if response, isChallenge := checkChallenge(lambdaEvent); isChallenge {
		log.Debug("Yep, this is a challenge, return 220 ok then")
		return response, nil
	}
	log.Debug("Nope, this is not a challenge. Processing then.")

	// If this request is not a challenge - proceed as reaction request
	log.WithFields(log.Fields{"reaction": lambdaEvent.Event.Reaction}).Debug("Start checking a reaction")
	if lambdaEvent.Event.Type == "reaction_added" && lambdaEvent.Event.Reaction == config.BestEmojiName {
		log.Debug("Yep, this is reaction we need to react. Will get all reactions for the message")
		reactions, err := client.GetReactions(slack.ItemRef{
			Channel:   lambdaEvent.Event.Item.Channel,
			Timestamp: lambdaEvent.Event.Item.Timestamp,
		}, slack.NewGetReactionsParameters())
		if err != nil {
			log.Error(err)
			return response, err
		}
		log.WithFields(log.Fields{"reactions": reactions}).Debugf("We got reactions for the message")

		log.Debug("Let's find a reaction we need...")
		for _, reaction := range reactions {
			if reaction.Name == config.BestEmojiName && reaction.Count == 1 {

				log.Debug("Lets create a permalink to the message we will quote")
				permalink, err := client.GetPermalink(&slack.PermalinkParameters{
					Channel: lambdaEvent.Event.Item.Channel,
					Ts:      lambdaEvent.Event.Item.Timestamp,
				})
				if err != nil {
					log.Error(err)
					return response, err
				}
				log.Debugf("Permalink created: %s", permalink)

				log.Debugf("Then we need to find a message in history, to copy a text from.")
				channelHistory, err := legacyClient.GetChannelHistory(lambdaEvent.Event.Item.Channel, slack.HistoryParameters{
					Latest:    lambdaEvent.Event.Item.Timestamp,
					Inclusive: true,
					Count:     1,
				})
				if err != nil {
					log.Error(err)
				}
				if len(channelHistory.Messages) != 1 {
					log.Errorf("Not a single message returned from a search: %v", channelHistory)
					return response, nil
				}

				log.Debug("So... We got a single message from a history. Let's process it.")
				headerString := fmt.Sprintf("<@%s> написал, _(а <@%s> добавил в лучшее)_: ", lambdaEvent.Event.ItemUser, lambdaEvent.Event.User)
				headerTextObject := slack.NewTextBlockObject("mrkdwn", headerString, false, false)
				headerSection := slack.NewSectionBlock(headerTextObject, nil, nil)

				bodyString := fmt.Sprintf("%s\n\n <%s|link to message>", channelHistory.Messages[0].Text, permalink)
				bodyText := slack.NewTextBlockObject("mrkdwn", bodyString, false, false)
				bodySection := slack.NewSectionBlock(bodyText, nil, nil)

				log.WithFields(log.Fields{"header": headerSection, "body": bodySection}).Debug("Message created and ready.")
				_, _, err = client.PostMessage(config.BestChannelId,
					slack.MsgOptionBlocks(headerSection),
					slack.MsgOptionAttachments(slack.Attachment{Blocks: []slack.Block{bodySection}}))
				if err != nil {
					log.Errorf("Problem with senging a message: %s", err)
				}
			}
		}
	}
	return response, err
}

type LambdaEvent struct {
	Token     string                   `json:"token"`
	Challenge string                   `json:"challenge"`
	Type      string                   `json:"type"`
	Event     slack.ReactionAddedEvent `json:"event"`
}

func parseBody(body string) (event LambdaEvent, err error) {
	err = json.Unmarshal([]byte(body), &event)
	return
}

// checkChallenge check request for a slack challenge.
// Slack uses same URL for the challenges and regular requests.
func checkChallenge(event LambdaEvent) (response events.APIGatewayProxyResponse, isChallenge bool) {
	if event.Challenge != "" {
		response.Body = event.Challenge
		response.StatusCode = 200
		isChallenge = true
	}
	return
}
