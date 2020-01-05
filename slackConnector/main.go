package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaEvent struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

func HandleEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event LambdaEvent
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}
	response := events.APIGatewayProxyResponse{StatusCode: 200, Body: event.Challenge}
	return response, nil
}

func main() {
	lambda.Start(HandleEvent)
}
