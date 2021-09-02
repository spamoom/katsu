package aws

import (
	cmdLogin "github.com/netsells/katsu/cmd/aws/docker/login"
	"github.com/spf13/cobra"
	// https://pkg.go.dev/github.com/MakeNowJust/heredoc
)

func NewCmdDocker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker <command>",
		Short: "Interact with Docker & AWS",
		Long:  "Work with Docker & AWS",
	}

	cmd.AddCommand(cmdLogin.NewCmdLogin())

	return cmd
}
