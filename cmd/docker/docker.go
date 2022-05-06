package docker

import (
	cmdBuild "github.com/netsells/katsu/cmd/docker/build"
	cmdPush "github.com/netsells/katsu/cmd/docker/push"
	"github.com/spf13/cobra"
	// https://pkg.go.dev/github.com/MakeNowJust/heredoc
)

func NewCmdDocker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker <command>",
		Short: "Interact with Docker",
		Long:  "Work with Docker",
	}

	cmd.AddCommand(cmdBuild.NewCmdBuild())
	cmd.AddCommand(cmdPush.NewCmdPush())

	return cmd
}
