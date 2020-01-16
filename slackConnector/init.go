package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nlopes/slack"
	"os"
	"reflect"
	"strings"
	"unicode"
)

// Variables to share between sessions
var (
	legacyClient *slack.Client
	client       *slack.Client
	sess         *session.Session
	config       *initConfig
)

func init() {
	// Don't run init on testing
	if len(os.Args) > 1 && os.Args[1][:5] == "-test" {
		return
	}
	sess, _ = session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("us-east-1")},
		SharedConfigState: session.SharedConfigEnable,
	})

	// config := config{
	// 	slackApiKeyName:      "/ukrops/emojiBot/slackAPI",
	// 	slackApiKeyLeacyName: "/ukrops/emojiBot/slackAPILegacy",
	// }
	// err := config.GetSecrets()
	// if err != nil {
	// 	panic(err)
	// }
	// client = slack.New(config.slackApiKey)
	// legacyClient = slack.New(config.slackApiKeyLeacy)
}

type initConfig struct {
	AwsRegion                string
	SsmSlackApiKeyPath       string
	SsmSlackApiLegacyKeyPath string
	BestChannelId            string
}

// FillFromEnv will check env variables with "EMOJIBOT_" prefix
// and fill config based on results
// Probably can use some package for this
// But while it's a my own package, can do it myself for fun
func (c *initConfig) FillFromEnv() error {
	envVariablesPrefix := "EMOJIBOT_"
	envVariablesRaw := os.Environ()

	reflectConfig := reflect.ValueOf(c)

	for _, envVarPair := range envVariablesRaw {
		if strings.HasPrefix(envVarPair, envVariablesPrefix) {
			envVarName := strings.SplitN(envVarPair, "=", 2)[0]
			trimedName := strings.TrimPrefix(envVarName, envVariablesPrefix)
			cammelCaseName := envToCamelCase(trimedName)
			field := reflectConfig.Elem().FieldByName(cammelCaseName)
			if field.IsValid() {
				field.SetString(os.Getenv(envVarName))
			}
		}
	}

	for i := 0; i < reflectConfig.Elem().NumField(); i++ {
		if reflectConfig.Elem().Field(i).String() == "" {
			return fmt.Errorf("error parsing %s argument", reflectConfig.Elem().Type().Field(i).Name)
		}
	}
	return nil
}

// envToCamelCase() is a helper function to change THIS_ENV_VAR to thisEnvVar
// It's mean every symbol after '_' should be an UPPERCASE
func envToCamelCase(incomingString string) string {
	isThisAFirstLetter := true
	mapFunc := func(r rune) rune {
		switch {
		case r == '_':
			isThisAFirstLetter = true
			return -1
		case isThisAFirstLetter:
			isThisAFirstLetter = false
			return unicode.ToUpper(r)
		default:
			return unicode.ToLower(r)
		}
	}
	return strings.Map(mapFunc, incomingString)
}

// func (c *config) GetSecrets() error {
// 	ps := ssm.New(sess)
// 	output, err := ps.GetParameter(&ssm.GetParameterInput{
// 		Name:           aws.String(c.slackApiKeyName),
// 		WithDecryption: aws.Bool(true),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	c.slackApiKey = *output.Parameter.Value
// 	output, err = ps.GetParameter(&ssm.GetParameterInput{
// 		Name:           aws.String(c.slackApiKeyLeacyName),
// 		WithDecryption: aws.Bool(true),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	c.slackApiKeyLeacy = *output.Parameter.Value
// 	return nil
// }
