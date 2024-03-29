package ecr

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/netsells/katsu/helpers/aws"
	"github.com/netsells/katsu/helpers/aws/sts"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/netsells/katsu/helpers/config"
	"github.com/netsells/katsu/helpers/docker"
)

func AuthenticateDocker() error {
	token, err := getLoginToken()

	if err != nil {
		return err
	}

	username, password := loginTokenToUserPassword(token)

	if err != nil {
		return errors.New("unable to get docker password from AWS")
	}

	cliio.LogVerbosef("Got ECR password: %s", token)

	// Making use of GetAwsAccountIdDefault and GetCallerAccountId so we can auto detect the ID when not provided
	repoHostname := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", config.GetAwsAccountIdDefault(sts.GetCallerAccountId()), config.GetAwsRegion())
	currentUser := sts.GetCallerArn()

	cliio.Lines([]string{
		fmt.Sprintf("Targeting registry %s", cliio.CommentText(repoHostname)),
		fmt.Sprintf("Using user %s", cliio.CommentText(currentUser)),
		"",
	})

	_, err = docker.Login(repoHostname, username, password)

	if err != nil {
		return errors.New("unable to login to docker")
	}

	return nil
}

func getLoginToken() (string, error) {

	ctx := context.Background()

	client := ecr.NewFromConfig(aws.GetConfig())

	input := &ecr.GetAuthorizationTokenInput{}

	output, err := client.GetAuthorizationToken(ctx, input)

	if err != nil {
		return "", err
	}

	return *output.AuthorizationData[0].AuthorizationToken, nil
}

func loginTokenToUserPassword(token string) (string, string) {
	decodedBytes, _ := base64.StdEncoding.DecodeString(token)
	userPasswordSplit := strings.Split(string(decodedBytes), ":")

	return userPasswordSplit[0], userPasswordSplit[1]
}
