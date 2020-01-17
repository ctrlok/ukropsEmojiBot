package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

var reactionAddedBody = `
{
	"token":"TOKEN",
	"team_id":"T0VFKM9MK",
	"api_app_id":"APP_ID",
	"event":{
		"type":"reaction_added",
		"user":"U0VFW0KPC",
		"item":{
			"type":"message",
			"channel":"CS9R5CUQG",
			"ts":"1578237353.000400"
		},
		"reaction":"sweat_smile",
		"item_user":"U0VFW0KPC",
		"event_ts":"1578245461.001600"
	},
	"type":"event_callback",
	"event_id":"EvRZKTUM5F",
	"event_time":1578245461,
	"authed_users":["U0VFW0KPC"]
}`

var challengeBody = `
{
  "token": "token",
  "challenge": "challenge",
  "type": "url_verification"
}`

func TestParseEvent(t *testing.T) {
	config = &initConfig{}
	request := events.APIGatewayProxyRequest{Body: challengeBody}
	event, err := ParseEvent(request)
	assert.NoError(t, err)
	assert.IsType(t, &ChallengeEvent{}, event)
	assert.Equal(t, "challenge", event.(*ChallengeEvent).Challenge)

	_, err = ParseEvent(events.APIGatewayProxyRequest{Body: ""})
	assert.Error(t, err)

	_, err = ParseEvent(events.APIGatewayProxyRequest{Body: `{"token": "token","type": "not_existed_type"}`})
	assert.Error(t, err)

	request = events.APIGatewayProxyRequest{Body: reactionAddedBody}
	config.BestEmojiName = "to_best"
	event, err = ParseEvent(request)
	assert.NoError(t, err)
	assert.IsType(t, ReactionEvent{}, event)
	assert.Equal(t, "sweat_smile", event.(ReactionEvent).Event.Reaction)

	request = events.APIGatewayProxyRequest{Body: reactionAddedBody}
	config.BestEmojiName = "sweat_smile"
	event, err = ParseEvent(request)
	assert.NoError(t, err)
	assert.IsType(t, ToBestEvent{}, event)
	assert.Equal(t, "sweat_smile", event.(ToBestEvent).Event.Reaction)
}
