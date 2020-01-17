package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChallenge_VerifyToken(t *testing.T) {
	c := ChallengeEvent{
		BasicEvent{Token:"123"},
		"challenge_string",
	}
	assert.True(t, c.VerifyToken("123"))
}

func TestChallengeEvent_Process(t *testing.T) {
	c := ChallengeEvent{
		Challenge:  "123",
	}
	response, err := c.Process()
	assert.NoError(t, err)
	assert.Equal(t, "123", response.Body)
	assert.Equal(t, 200, response.StatusCode)
}

