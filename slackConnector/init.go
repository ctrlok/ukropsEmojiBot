package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/nlopes/slack"
	"os"
	"reflect"
	"strings"
	"unicode"
)

// Variables to share between Lambda runs
var (
	legacyClient *slack.Client
	client       *slack.Client
	sess         *session.Session
	config       *initConfig
)

// We use init for parse env variables
// and for creating slack sessions
// shared between Lambda runs
// Don't know how to test it without expensive integration tests :(
func init() {
	// Don't run init on testing
	if len(os.Args) > 1 && os.Args[1][:5] == "-test" {
		return
	}

	err := config.FillFromEnv()
	if err != nil {
		fmt.Printf("Initialization error: %s", err)
		os.Exit(2)
	}

	sess, _ = session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(config.AwsRegion)},
		SharedConfigState: session.SharedConfigEnable,
	})

	slackApiKey, err := GetSecrets(config.SsmSlackApiKeyPath)
	if err != nil {
		fmt.Printf("Error getting ssm key %s: %s", config.SsmSlackApiKeyPath, err)
		os.Exit(2)
	}
	client = slack.New(slackApiKey)

	slackApiKeyLegacy, err := GetSecrets(config.SsmSlackApiLegacyKeyPath)
	if err != nil {
		fmt.Printf("Error getting ssm key %s: %s", config.SsmSlackApiLegacyKeyPath, err)
		os.Exit(2)
	}
	legacyClient = slack.New(slackApiKeyLegacy)
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

	// Set initConfig fields which are mached to the env variables with EMOJIBOT_ prefix
	for _, envVarPair := range envVariablesRaw {
		if strings.HasPrefix(envVarPair, envVariablesPrefix) {
			envVarName := strings.SplitN(envVarPair, "=", 2)[0]
			trimedName := strings.TrimPrefix(envVarName, envVariablesPrefix)
			cammelCaseName := envToCamelCase(trimedName)
			field := reflectConfig.Elem().FieldByName(cammelCaseName)
			if field.IsValid() {
				field.SetString(os.Getenv(envVarName)) // Set only exist fields
			}
		}
	}

	// Check if all fields are setted up
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

// GetSecrets is a function for getting secrets from AWS SSM based on provided path
func GetSecrets(secretPath string) (string, error) {
	ps := ssm.New(sess)
	output, err := ps.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(secretPath),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *output.Parameter.Value, nil
}
