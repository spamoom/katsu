package aws

import (
	cmdAssumeRole "github.com/netsells/katsu/cmd/aws/assume-role"
	cmdDocker "github.com/netsells/katsu/cmd/aws/docker"
	cmdEcs "github.com/netsells/katsu/cmd/aws/ecs"
	cmdSsm "github.com/netsells/katsu/cmd/aws/ssm"
	"github.com/spf13/cobra"
	// https://pkg.go.dev/github.com/MakeNowJust/heredoc
)

func NewCmdAws() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws <command>",
		Short: "Interact with AWS",
		Long:  "Work with AWS",
	}

	// These will persist through all the aws sub-commands
	cmd.PersistentFlags().String("aws-region", "", "Override the default AWS region")
	cmd.PersistentFlags().String("aws-account-id", "", "Override the default AWS account ID")
	cmd.PersistentFlags().String("aws-profile", "", "Override the AWS profile to use")

	cmd.AddCommand(cmdAssumeRole.NewCmdAssumeRole())
	cmd.AddCommand(cmdDocker.NewCmdDocker())
	cmd.AddCommand(cmdEcs.NewCmdEcs())
	cmd.AddCommand(cmdSsm.NewCmdSsm())

	return cmd
}
