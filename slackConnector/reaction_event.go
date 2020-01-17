package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
)

// TODO: implement storing reactions at DynamoDB
type ReactionEvent struct {
	BasicEvent
	Event     slack.ReactionAddedEvent `json:"event"`
}

func (r ReactionEvent) Process() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{StatusCode:200}, nil
}

