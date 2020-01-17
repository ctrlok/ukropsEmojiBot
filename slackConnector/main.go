package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

func main() {
	lambda.Start(HandleEvent)
}

func HandleEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.WithFields(log.Fields{"body": request.Body}).Debug("Start processing a body")
	event, err := ParseEvent(request)
	if err != nil {
		log.WithFields(log.Fields{"request": request.Body}).Error("Incorect event type")
		return events.APIGatewayProxyResponse{StatusCode:500}, err
	}
	log.WithFields(log.Fields{"event": event}).Debug("Finish processing a body. Let's start event processing")
	// TODO: add token verification
	return event.Process()
}

func ParseEvent(request events.APIGatewayProxyRequest) (Event, error) {
	log.Debug("Start processing body to a basic event a basic event")
	var event Event
	basicEvent := BasicEvent{}
	err := json.Unmarshal([]byte(request.Body), &basicEvent)
	if err != nil {
		return event, err
	}
	log.WithFields(log.Fields{"basicEvent": basicEvent}).Debug("Processing finishing. Let's find a type of this event")
	switch basicEvent.Type {
	case "url_verification":
		log.Debug("Looks like this is a ChallengeEvent. Processing.")
		event = &ChallengeEvent{}
		err = json.Unmarshal([]byte(request.Body), event)
	case "event_callback":
		log.Debug("And this is a reaction item. Interesting...")
		rEvent := ReactionEvent{}
		err = json.Unmarshal([]byte(request.Body), &rEvent)
		log.Debug("Is this a reaction which are supposed be a best reaction?")
		switch rEvent.Event.Reaction {
		case config.BestEmojiName:
			log.Debug("Yep, this is a best reaction. Process then.")
			event = ToBestEvent{rEvent}
		default:
			log.Debug("Nope. It's a usual reaction.")
			event = rEvent
		}
	default:
		err = fmt.Errorf("unknown type of event: %s", basicEvent.Type)
	}
	return event, err
}

type BasicEvent struct {
	Token string `json:"token"`
	Type string `json:"type"`
}

func (e BasicEvent) VerifyToken(expectedToken string) bool {
	return expectedToken == e.Token
}

type Event interface {
	VerifyToken(string) bool
	Process() (events.APIGatewayProxyResponse, error)
}

// func HandleEventTest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//
// 	// If this request is not a challenge - proceed as reaction request
// 	log.WithFields(log.Fields{"reaction": lambdaEvent.Event.Reaction}).Debug("Start checking a reaction")
// 	if lambdaEvent.Event.Type == "reaction_added" && lambdaEvent.Event.Reaction == config.BestEmojiName {
//
// 		log.Debug("Let's find a reaction we need...")
// 		for _, reaction := range reactions {
// 			if reaction.Name == config.BestEmojiName && reaction.Count == 1 {
//
// 				log.Debug("Lets create a permalink to the message we will quote")
// 				permalink, err := client.GetPermalink(&slack.PermalinkParameters{
// 					Channel: lambdaEvent.Event.Item.Channel,
// 					Ts:      lambdaEvent.Event.Item.Timestamp,
// 				})
// 				if err != nil {
// 					log.Error(err)
// 					return response, err
// 				}
// 				log.Debugf("Permalink created: %s", permalink)
//
// 				log.Debugf("Then we need to find a message in history, to copy a text from.")
// 				channelHistory, err := legacyClient.GetChannelHistory(lambdaEvent.Event.Item.Channel, slack.HistoryParameters{
// 					Latest:    lambdaEvent.Event.Item.Timestamp,
// 					Inclusive: true,
// 					Count:     1,
// 				})
// 				if err != nil {
// 					log.Error(err)
// 				}
// 				if len(channelHistory.Messages) != 1 {
// 					log.Errorf("Not a single message returned from a search: %v", channelHistory)
// 					return response, nil
// 				}
//
// 				log.Debug("So... We got a single message from a history. Let's process it.")
// 				headerString := fmt.Sprintf("<@%s> написал, _(а <@%s> добавил в лучшее)_: ", lambdaEvent.Event.ItemUser, lambdaEvent.Event.User)
// 				headerTextObject := slack.NewTextBlockObject("mrkdwn", headerString, false, false)
// 				headerSection := slack.NewSectionBlock(headerTextObject, nil, nil)
//
// 				bodyString := fmt.Sprintf("%s\n\n <%s|link to message>", channelHistory.Messages[0].Text, permalink)
// 				bodyText := slack.NewTextBlockObject("mrkdwn", bodyString, false, false)
// 				bodySection := slack.NewSectionBlock(bodyText, nil, nil)
//
// 				log.WithFields(log.Fields{"header": headerSection, "body": bodySection}).Debug("Message created and ready.")
// 				_, _, err = client.PostMessage(config.BestChannelId,
// 					slack.MsgOptionBlocks(headerSection),
// 					slack.MsgOptionAttachments(slack.Attachment{Blocks: []slack.Block{bodySection}}))
// 				if err != nil {
// 					log.Errorf("Problem with senging a message: %s.", err)
// 				}
// 			}
// 		}
// 	}
// 	return response, err
// }
//
// type LambdaEvent struct {
// 	Token     string                   `json:"token"`
// 	Challenge string                   `json:"challenge"`
// 	Type      string                   `json:"type"`
// 	Event     slack.ReactionAddedEvent `json:"event"`
// }
