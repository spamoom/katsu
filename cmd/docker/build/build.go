package cmd

import (
	"os"

	"github.com/netsells/katsu/helpers"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/netsells/katsu/helpers/config"
	"github.com/netsells/katsu/helpers/docker"
	"github.com/netsells/katsu/helpers/process"
	"github.com/spf13/cobra"
)

type CallBuildContext struct {
	Tag         string
	Service     string
	Services    []string
	Environment config.ConfigEnvironment
}

func NewCmdBuild() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Builds docker-compose ready for an environment",
		Run:   runDockerBuildCmd,
	}

	cmd.Flags().String("tag", helpers.GetCurrentSha(), "The tag that should be built with the images. Defaults to the current commit SHA")
	cmd.Flags().String("tag-prefix", "", "The tag prefix that should be built with the images. Defaults to null")
	cmd.Flags().String("environment", "", "The destination environment for the images")
	cmd.Flags().StringArray("services", []string{}, "The service that should be built. Not defining this will push all services")

	return cmd
}

func runDockerBuildCmd(cmd *cobra.Command, args []string) {
	helpers.SetCmd(cmd)

	if config.GetTag() == "" {
		cliio.ErrorStep("No tag set or available from git. Cannot proceed.")
		os.Exit(1)
	}

	helpers.CheckAndReportMissingBinaries([]string{"docker"})
	helpers.CheckAndReportMissingFiles([]string{"docker-compose.yml", "docker-compose.build.yml"})

	environment, err := config.GetCurrentEnvironment()

	if err != nil {
		cliio.FatalStepf("Unable to get environment, please ensure it is setup in .katsu.yml")
	}

	prefixedTag := docker.DockerPrefixedTag()
	services := config.GetDockerServices()

	if len(services) == 0 {
		cliio.Stepf("Building docker images for all services with tag %s", prefixedTag)

		success := callBuild(CallBuildContext{
			Tag:         prefixedTag,
			Environment: *environment,
			Services:    services,
		})

		if success {
			cliio.SuccessfulStep("Docker images built.")
			os.Exit(0)
		}

		os.Exit(1)
	}

	cliio.Stepf("Building docker images for services with tag %s: %v", prefixedTag, services)

	for _, service := range services {
		success := callBuild(CallBuildContext{
			Tag:         prefixedTag,
			Service:     service,
			Environment: *environment,
			Services:    services,
		})

		if !success {
			os.Exit(1)
		}
	}

	cliio.SuccessfulStep("Docker images built.")
	os.Exit(0)
}

func callBuild(context CallBuildContext) bool {
	parts := []string{
		"-f", "docker-compose.yml",
		"-f", "docker-compose.build.yml",
		"build", "--no-cache",
	}

	if context.Service != "" {
		parts = append(parts, context.Service)
	}

	process := process.NewProcess("docker-compose", parts...)
	process.SetEnv("TAG", context.Tag)

	docker.AppendEnvImageReposForServices(process, context.Services)

	process.EchoLineByLine = true
	_, err := process.Run()

	if err != nil {
		cliio.ErrorStep("Unable to build all images, check the above output for reasons why.")
		return false
	}

	return true
}
