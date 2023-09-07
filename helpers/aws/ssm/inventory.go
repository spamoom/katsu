package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws"
	katsuAws "github.com/netsells/katsu/helpers/aws"
)

type Instance struct {
	Id        string
	IpAddress string
}

func GetInstances() ([]Instance, error) {
	ctx := context.Background()
	client := ssm.NewFromConfig(katsuAws.GetConfig())

	output, err := client.GetInventory(ctx, &ssm.GetInventoryInput{
		Filters: []types.InventoryFilter{
			{
				Key:    aws.String("AWS:InstanceInformation.InstanceStatus"),
				Values: []string{"Active"},
				Type:   "Equal",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var instances []Instance

	for _, item := range output.Entities {
		instances = append(instances, Instance{
			Id:        *item.Id,
			IpAddress: item.Data["AWS:InstanceInformation"].Content[0]["IpAddress"],
		})
	}

	return instances, nil
}
