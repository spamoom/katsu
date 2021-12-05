package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	katsuAws "github.com/netsells/katsu/helpers/aws"
)

func GetUserMfaDevices() ([]types.MFADevice, error) {
	ctx := context.Background()
	client := iam.NewFromConfig(katsuAws.GetConfig())

	input := &iam.ListMFADevicesInput{}

	devices, err := client.ListMFADevices(ctx, input)

	if err != nil {
		return nil, err
	}

	return devices.MFADevices, nil
}
