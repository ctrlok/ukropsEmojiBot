package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChallengeEvent_Process(t *testing.T) {
	c := ChallengeEvent{
		Challenge: "123",
	}
	response, err := c.Process()
	assert.NoError(t, err)
	assert.Equal(t, "123", response.Body)
	assert.Equal(t, 200, response.StatusCode)
}
