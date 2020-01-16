package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_envToCamelCase(t *testing.T) {
	assert.Equal(t, "BigBraveString", envToCamelCase("BIG_BRAVE_STRING"))
}

func Test_initConfig_FillFromEnv(t *testing.T) {
	// Fail because of empty fields
	c := &initConfig{}
	err := c.FillFromEnv()
	assert.Error(t, err)

	// Succ test
	os.Setenv("EMOJIBOT_AWS_REGION", "us-west-1")
	os.Setenv("EMOJIBOT_BEST_CHANNEL_ID", "CHANID")
	os.Setenv("EMOJIBOT_SSM_SLACK_API_KEY_PATH", "APIPATH1")
	os.Setenv("EMOJIBOT_SSM_SLACK_API_LEGACY_KEY_PATH", "L_APIPATH")
	c = &initConfig{}
	err = c.FillFromEnv()
	assert.NoError(t, err)
	assert.Equal(t, "us-west-1", c.AwsRegion)
	assert.Equal(t, "CHANID", c.BestChannelId)
	assert.Equal(t, "APIPATH1", c.SsmSlackApiKeyPath)
	assert.Equal(t, "L_APIPATH", c.SsmSlackApiLegacyKeyPath)

	// Test panic
	c = &initConfig{}
	_ = os.Setenv("EMOJIBOT_NOTEXIST_FIELD", "non-field-value")
	assert.NotPanics(t, func() { _ = c.FillFromEnv() })
}
