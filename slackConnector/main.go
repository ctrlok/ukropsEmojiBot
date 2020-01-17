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
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
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
	Type  string `json:"type"`
}

func (e BasicEvent) VerifyToken(expectedToken string) bool {
	return expectedToken == e.Token
}

type Event interface {
	VerifyToken(string) bool
	Process() (events.APIGatewayProxyResponse, error)
}
