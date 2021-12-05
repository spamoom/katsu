package sts

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/netsells/katsu/helpers/aws"
)

type AssumeRoleEnvVars struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

type AssumeRoleInput struct {
	AccountId   string
	Role        string
	SessionUser string
	MfaDevice   *string
	MfaCode     *string
}

func AssumeRole(input AssumeRoleInput) (*AssumeRoleEnvVars, error) {
	ctx := context.Background()

	client := sts.NewFromConfig(aws.GetConfig())

	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", input.AccountId, input.Role)
	sessionName := fmt.Sprintf("%s-on-%s", input.SessionUser, input.AccountId)

	assumeRoleInput := &sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: &sessionName,
	}

	if input.MfaDevice != nil && input.MfaCode != nil {
		assumeRoleInput.SerialNumber = input.MfaDevice
		assumeRoleInput.TokenCode = input.MfaCode
	}

	output, err := client.AssumeRole(ctx, assumeRoleInput)

	if err != nil {
		return nil, err
	}

	return &AssumeRoleEnvVars{
		AccessKeyID:     *output.Credentials.AccessKeyId,
		SecretAccessKey: *output.Credentials.SecretAccessKey,
		SessionToken:    *output.Credentials.SessionToken,
	}, nil
}
