package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	katsuConfig "github.com/netsells/katsu/helpers/config"
)

func GetConfig() aws.Config {

	awsConfig := config.WithRegion(katsuConfig.GetAwsRegion())

	if katsuConfig.GetAwsProfile() != "" {
		awsConfig = config.WithSharedConfigProfile(katsuConfig.GetAwsProfile())
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), awsConfig)

	if err != nil {
		panic(err)
	}

	return cfg
}
