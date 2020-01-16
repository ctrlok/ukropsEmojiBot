package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
)

func main() {
	lambda.Start(HandleEventTest)
}

func HandleEventTest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}

	lambdaEvent, err := parseBody(request.Body)
	if err != nil {
		fmt.Println(err)
		return response, err
	}

	// is this request a slack challenge?
	if response, isChallenge := checkChallenge(lambdaEvent); isChallenge {
		return response, nil
	}

	// If this request is not a challenge - proceed as reaction request
	if lambdaEvent.Event.Type == "reaction_added" && lambdaEvent.Event.Reaction == "to_best" {
		reactions, err := client.GetReactions(slack.ItemRef{
			Channel:   lambdaEvent.Event.Item.Channel,
			Timestamp: lambdaEvent.Event.Item.Timestamp,
		}, slack.NewGetReactionsParameters())
		if err != nil {
			fmt.Println(err)
			return response, err
		}

		for _, reaction := range reactions {
			if reaction.Name == config.BestEmojiName && reaction.Count == 1 {
				permalink, _ := client.GetPermalink(&slack.PermalinkParameters{
					Channel: lambdaEvent.Event.Item.Channel,
					Ts:      lambdaEvent.Event.Item.Timestamp,
				})

				channelHistory, err := legacyClient.GetChannelHistory(lambdaEvent.Event.Item.Channel, slack.HistoryParameters{
					Latest:    lambdaEvent.Event.Item.Timestamp,
					Inclusive: true,
					Count:     1,
				})
				if err != nil {
					fmt.Println(err)
				}
				if len(channelHistory.Messages) != 1 {
					fmt.Printf("Not a single message returned from a search: %v", channelHistory)
					return response, nil
				}
				headerString := fmt.Sprintf("<@%s> написал, _(а <@%s> добавил в лучшее):_", lambdaEvent.Event.ItemUser, lambdaEvent.Event.User)
				headerTextObject := slack.NewTextBlockObject("mrkdwn", headerString, false, false)
				headerSection := slack.NewSectionBlock(headerTextObject, nil, nil)

				bodyString := fmt.Sprintf("%s\n\n <%s|link to message>", channelHistory.Messages[0].Text, permalink)
				bodyText := slack.NewTextBlockObject("mrkdwn", bodyString, false, false)
				bodySection := slack.NewSectionBlock(bodyText, nil, nil)

				_, _, err = client.PostMessage(config.BestChannelId,
					slack.MsgOptionBlocks(headerSection),
					slack.MsgOptionAttachments(slack.Attachment{Blocks: []slack.Block{bodySection}}))
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
