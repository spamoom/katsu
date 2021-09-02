package aws

import (
	cmdAssumeRole "github.com/netsells/katsu/cmd/aws/assume-role"
	cmdDocker "github.com/netsells/katsu/cmd/aws/docker"
	"github.com/spf13/cobra"
	// https://pkg.go.dev/github.com/MakeNowJust/heredoc
)

func NewCmdAws() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws <command>",
		Short: "Interact with AWS",
		Long:  "Work with AWS",
	}

	cmd.AddCommand(cmdAssumeRole.NewCmdAssumeRole())
	cmd.AddCommand(cmdDocker.NewCmdDocker())

	return cmd
}
