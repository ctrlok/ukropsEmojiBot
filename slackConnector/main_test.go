package main

import (
	"testing"
)

var reactionBody = `{
	"token":"TOKEN",
	"team_id":"T0VFKM9MK",
	"api_app_id":"APP_ID",
	"event":{
		"type":"reaction_removed",
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
	"authed_users":["U0VFW0KPC"]}`

func Test_parseBody(t *testing.T) {

}
