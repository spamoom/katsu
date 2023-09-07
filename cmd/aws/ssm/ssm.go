package aws

import (
	cmdConnect "github.com/netsells/katsu/cmd/aws/ssm/connect"
	"github.com/spf13/cobra"
	// https://pkg.go.dev/github.com/MakeNowJust/heredoc
)

func NewCmdSsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssm <command>",
		Short: "Interact with AWS SSM",
		Long:  "Work with AWS SSM",
	}

	cmd.AddCommand(cmdConnect.NewCmdConnect())

	return cmd
}
