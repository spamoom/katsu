package cmd

import (
	"os"

	"github.com/netsells/katsu/helpers"
	"github.com/netsells/katsu/helpers/aws"
	"github.com/netsells/katsu/helpers/aws/ecr"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/spf13/cobra"
)

var dockerLoginCmd = &cobra.Command{
	Use:   "docker:aws:login",
	Short: "Logs into docker via the AWS account",
	Run:   runDockerLoginCmd,
}

func init() {
	rootCmd.AddCommand(dockerLoginCmd)

	aws.RegisterCommonFlags(dockerLoginCmd)
}

func runDockerLoginCmd(cmd *cobra.Command, args []string) {
	helpers.SetCmd(cmd)

	helpers.CheckAndReportMissingBinaries([]string{"docker"})

	cliio.Step("Logging into docker")

	err := ecr.AuthenticateDocker()
	if err != nil {
		cliio.FatalStep(err.Error())
	}

	cliio.SuccessfulStep("Successfully logged into docker")

	os.Exit(0)
}
