package cmd

import (
	"os"
	"strings"

	"github.com/netsells/katsu/helpers"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/netsells/katsu/helpers/config"
	"github.com/netsells/katsu/helpers/docker"
	"github.com/netsells/katsu/helpers/process"
	"github.com/spf13/cobra"
)

type CallPushContext struct {
	Tags        []string
	Service     string
	Services    []string
	Environment config.ConfigEnvironment
}

func NewCmdPush() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Pushes built images to an environment",
		Run:   runDockerPushCmd,
	}

	cmd.Flags().Bool("skip-additional-tags", false, "Skips the latest and environment tags")

	cmd.Flags().String("tag", helpers.GetCurrentSha(), "The tag that should be built with the images. Defaults to the current commit SHA")
	cmd.Flags().String("tag-prefix", "", "The tag prefix that should be built with the images. Defaults to null")
	cmd.Flags().String("environment", "", "The destination environment for the images")
	cmd.Flags().StringArray("services", []string{}, "The service that should be built. Not defining this will push all services")

	return cmd
}

func runDockerPushCmd(cmd *cobra.Command, args []string) {
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

	services := config.GetDockerServices()

	skipAdditonalTags, _ := cmd.Flags().GetBool("skip-additional-tags")
	tags := docker.DetermineTags(skipAdditonalTags)

	if len(services) == 0 {
		cliio.Stepf("Pushing docker images for all services with tags %s", strings.Join(tags, ", "))

		success := callPush(CallPushContext{
			Tags:        tags,
			Environment: *environment,
			Services:    services,
		})

		if success {
			cliio.SuccessfulStep("Docker images built.")
			os.Exit(0)
		}

		os.Exit(1)
	}

	cliio.Stepf("Pushing docker images for services with tags %s: %s", strings.Join(tags, ", "), strings.Join(services, ", "))

	for _, service := range services {
		success := callPush(CallPushContext{
			Tags:        tags,
			Service:     service,
			Environment: *environment,
			Services:    services,
		})

		if !success {
			os.Exit(1)
		}
	}

	cliio.SuccessfulStep("All images pushed.")
	os.Exit(0)
}

func callPush(context CallPushContext) bool {

	// Need to make the new tags first
	docker.TagImages(context.Tags, &context.Service)

	for _, tag := range context.Tags {
		parts := []string{
			"-f", "docker-compose.yml",
			"-f", "docker-compose.build.yml",
			"push",
		}

		if context.Service != "" {
			parts = append(parts, context.Service)
		}

		process := process.NewProcess("docker-compose", parts...)
		process.SetEnv("TAG", tag)

		docker.AppendEnvImageReposForServices(process, context.Services)

		process.EchoLineByLine = false
		_, err := process.Run()

		if err != nil {
			cliio.ErrorStepf("Unable to push tag %s for service %s, check the above output for reasons why.", tag, context.Service)
			return false
		} else {
			cliio.SuccessfulStepf("Pushed %s image with %s tag", context.Service, tag)
		}
	}

	return true
}
