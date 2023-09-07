package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	katsuAws "github.com/netsells/katsu/helpers/aws"
)

type Instance struct {
	Id               string
	PrivateIpAddress string
	Name             string
	Type             string
}

func GetInstancesById(ids []string) ([]Instance, error) {
	ctx := context.Background()
	client := ec2.NewFromConfig(katsuAws.GetConfig())

	output, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: ids,
	})

	if err != nil {
		return nil, err
	}

	var instances []Instance

	for _, item := range output.Reservations {
		nameValue := ""
		for _, tag := range item.Instances[0].Tags {
			if *tag.Key == "Name" {
				nameValue = *tag.Value
			}
		}

		instances = append(instances, Instance{
			Id:               *item.Instances[0].InstanceId,
			PrivateIpAddress: *item.Instances[0].PrivateIpAddress,
			Name:             nameValue,
			Type:             string(item.Instances[0].InstanceType),
		})
	}

	return instances, nil
}
