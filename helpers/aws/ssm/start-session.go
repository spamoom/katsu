package ssm

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mmmorris1975/ssm-session-client/ssmclient"
)

type StartSessionInput struct {
	InstanceId string
}

func StartSesson(input StartSessionInput) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	in := ssmclient.PortForwardingInput{
		Target:     input.InstanceId,
		RemotePort: 22,
	}

	log.Fatal(ssmclient.SSHSession(cfg, &in))

	return nil
}
