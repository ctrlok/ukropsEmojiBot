package main

import "github.com/aws/aws-lambda-go/events"

type ChallengeEvent struct {
	BasicEvent
	Challenge string                   `json:"challenge"`
}

func (c *ChallengeEvent) Process() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body:c.Challenge, StatusCode:200}, nil
}
