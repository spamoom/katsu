package aws

import (
	cmdManageEnv "github.com/netsells/katsu/cmd/aws/ecs/manage-env"
	"github.com/spf13/cobra"
	// https://pkg.go.dev/github.com/MakeNowJust/heredoc
)

func NewCmdEcs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecs <command>",
		Short: "Interact with AWS ECS",
		Long:  "Work with AWS ECS",
	}

	cmd.AddCommand(cmdManageEnv.NewCmdManageEnv())

	return cmd
}
