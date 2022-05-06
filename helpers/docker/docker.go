package docker

import (
	"fmt"
	"strings"

	"github.com/netsells/katsu/helpers/cliio"
	"github.com/netsells/katsu/helpers/config"
	"github.com/netsells/katsu/helpers/process"
)

func DockerPrefixedTag() string {
	tag := config.GetTag()

	tagPrefix := config.GetTaxPrefix()

	if tagPrefix != "" {
		return tagPrefix + "-" + tag
	}

	environmentTagPrefix := config.GetEnvironment()

	if environmentTagPrefix != "" {
		return environmentTagPrefix + "-" + tag
	}

	return tag
}

func DetermineTags(skipAdditionalTags bool) []string {
	tag := DockerPrefixedTag()

	tags := []string{tag}

	if !skipAdditionalTags {
		tags = append(tags, "latest")

		if config.GetEnvironment() != "" {
			tags = append(tags, config.GetEnvironment())
		}
	}

	return tags
}

func TagImages(newTags []string, service *string) bool {
	var services []string
	sourceTag := DockerPrefixedTag()

	if service == nil {
		cliio.FatalStep("Tag images with all services not yet implemented.")
	} else {
		services = []string{
			BuildRepoUrlForService(*service),
		}
	}

	for _, serviceUrl := range services {
		for _, tag := range newTags {
			if sourceTag == tag {
				// No point tagging the same thing
				continue
			}

			cliio.Stepf("Tagging %s:%s to %s:%s", serviceUrl, sourceTag, serviceUrl, tag)

			parts := []string{
				"tag",
				fmt.Sprintf("%s:%s", serviceUrl, sourceTag),
				fmt.Sprintf("%s:%s", serviceUrl, tag),
			}

			process := process.NewProcess("docker", parts...)
			process.SetEnv("TAG", tag)

			_, err := process.Run()

			if err != nil {
				cliio.ErrorStepf("Unable to tag %s:%s as %s:%s", serviceUrl, sourceTag, serviceUrl, tag)
				return false
			}
		}
	}

	return true
}

func BuildRepoUrlForService(service string) string {
	environment, err := config.GetCurrentEnvironment()

	if err != nil {
		cliio.FatalStepf("Unable to get environment, please ensure it is setup in .katsu.yml")
	}

	fmt.Println(service)

	return fmt.Sprintf("%d.dkr.ecr.%s.amazonaws.com/%s", environment.Aws.AccountId, environment.Aws.Region, environment.Aws.Ecs.Services[0].Ecr)
}

func AppendEnvImageReposForServices(process *process.Process, services []string) {
	for _, service := range services {

		serviceRepo := BuildRepoUrlForService(service)

		serviceEnvName := strings.ToUpper(fmt.Sprintf("%s_IMAGE_REPO", service))

		process.SetEnv(serviceEnvName, serviceRepo)

		cliio.LogVerbosef("Setting env %s to %s", serviceEnvName, serviceRepo)
	}
}
