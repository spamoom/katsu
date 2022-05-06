package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/netsells/katsu/helpers"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/spf13/viper"
)

type ConfigEnvironment struct {
	Name string
	Aws  ConfigEnvironmentAws
}

type ConfigEnvironmentAws struct {
	Name      string
	Region    string
	AccountId int `mapstructure:"account-id"`
	Ecs       ConfigEnvironmentAwsEcs
}

type ConfigEnvironmentAwsEcs struct {
	TaskDefinition string `mapstructure:"task-definition"`
	Service        string
	Services       []ConfigEnvironmentAwsEcsService
}

type ConfigEnvironmentAwsEcsService struct {
	Name string
	Ecr  string
}

func GetTag() string {
	return getString("tag", "", "")
}

func GetTaxPrefix() string {
	return getString("tag-prefix", "", "")
}

func GetEnvironment() string {
	return getString("environment", "", "")
}

func GetDefaultEnvironment() string {
	return getString("default-environment", "default-environment", "")
}

func GetAwsRegion() string {
	return getString("aws-region", "docker.aws.region", "eu-west-2")
}

func GetAwsProfile() string {
	return getString("aws-profile", "", "")
}

func GetAwsAccountId() string {
	return getString("aws-account-id", "docker.aws.account-id", "")
}

func GetAwsAccountIdDefault(defaultId string) string {
	return getString("aws-account-id", "docker.aws.account-id", defaultId)
}

func GetDockerServices() []string {
	return getStringArray("services", "docker.services")
}

func GetS3Bucket() string {
	return getString("s3-bucket", "", "")
}

func GetCurrentEnvironment() (*ConfigEnvironment, error) {
	environmentName := GetEnvironment()

	if environmentName == "" {
		cliio.LogVerbose("No environment name, using default environment")
		environmentName = GetDefaultEnvironment()
	}

	return GetNamedEnvironment(environmentName)
}

func GetNamedEnvironment(environmentName string) (*ConfigEnvironment, error) {
	v := viper.GetViper()
	environments := v.Get("environments")

	// Loop through environments
	for _, environment := range environments.([]interface{}) {
		var env ConfigEnvironment
		err := mapstructure.Decode(environment, &env)
		if err != nil {
			cliio.LogVerbosef("Failed to decode environment file. %s", err.Error())
			os.Exit(1)
		}

		if env.Name == environmentName {
			return &env, nil
		}
	}

	return nil, fmt.Errorf("environment [%s] not found", environmentName)
}

func getString(flag string, filePath string, defaultValue string) string {
	v := viper.GetViper()

	cliio.LogVerbosef("Fetching config for flag %s", flag)

	// Try from cli argument
	value, _ := helpers.GetCmd().Flags().GetString(flag)

	if value != "" {
		cliio.LogVerbosef("Got value for flag %s - %s", flag, value)
		return value
	}

	if filePath != "" {
		cliio.LogVerbosef("Now trying the katsu file in path %s", filePath)

		pathValue := v.GetString(filePath)

		if pathValue != "" {
			cliio.LogVerbosef("Got %s from file path %s", pathValue, filePath)
			return pathValue
		}
	}

	cliio.LogVerbosef("Unable to get value from flag or katsu file, falling back to default: %s", defaultValue)

	return defaultValue
}

func getStringArray(flag string, filePath string) []string {
	v := viper.GetViper()

	// Try from cli argument
	values, _ := helpers.GetCmd().Flags().GetStringArray(flag)

	if len(values) > 0 {
		return values
	}

	if filePath == "" {
		return []string{}
	}

	return v.GetStringSlice(filePath)
}
