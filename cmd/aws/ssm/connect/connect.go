package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/netsells/katsu/helpers"
	"github.com/netsells/katsu/helpers/aws/ec2"
	"github.com/netsells/katsu/helpers/aws/ssm"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/spf13/cobra"
)

type Instance struct {
	Id               string
	IpAddress        string
	PrivateIpAddress string
	Name             string
	Type             string
}

func NewCmdConnect() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Starts an SSM SSH session",
		Run:   runConnectCmd,
	}
	return cmd
}

func runConnectCmd(cmd *cobra.Command, args []string) {
	helpers.SetCmd(cmd)

	instances := getInstances()

	if len(instances) == 0 {
		cliio.FatalStep("No instances available")
	}

	instancePrompt := getInstancePrompt(instances)
	instanceIndex := askForInstance(&instancePrompt)

	usernameOptions := []string{"ubuntu", "ec2-user", "admin", "root", "custom username"}
	usernamePrompt := getUsernamePrompt(usernameOptions)
	usernameIndex := askForUsername(&usernamePrompt)
	username := usernameOptions[usernameIndex]

	if username == "custom username" {
		prompt := promptui.Prompt{
			Label: "Enter Username",
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		username = result
	}

	cliio.SuccessfulStep("Starting SSM SSH session")

	cliio.Line(string(username))

	ssm.StartSesson(ssm.StartSessionInput{
		InstanceId: instances[instanceIndex].Id,
	})

	os.Exit(0)
}

func getInstances() []Instance {
	instances, err := ssm.GetInstances()

	if err != nil {
		cliio.FatalStep(err.Error())
	}

	// We need to enrich the fleet data with info from ec2
	instanceIds := make([]string, len(instances))
	for i, instance := range instances {
		instanceIds[i] = instance.Id
	}

	instancesMeta, err := ec2.GetInstancesById(instanceIds)
	if err != nil {
		cliio.FatalStep(err.Error())
	}

	formattedInstances := make([]Instance, len(instances))

	for i, instance := range instances {

		instanceStruct := Instance{
			Id:        instance.Id,
			IpAddress: instance.IpAddress,
		}

		// Find the meta
		for _, meta := range instancesMeta {
			if instance.Id == meta.Id {
				instanceStruct.Name = meta.Name
				instanceStruct.Type = meta.Type
				instanceStruct.PrivateIpAddress = meta.PrivateIpAddress
			}
		}

		formattedInstances[i] = instanceStruct
	}

	return formattedInstances
}

func getInstancePrompt(instances []Instance) promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "[{{ .Id }}] .Name?",
		Active:   "\U000027A1 [{{ .Id | blue }}] {{ .Name | blue }}",
		Inactive: "  [{{ .Id | white }}] {{ .Name}}",
		Selected: fmt.Sprintf("%s Instance: {{ .Id | faint }} {{ .Name | faint }}", promptui.IconGood),
		Details: `
--------- Instance ----------
{{ "ID:" | faint }}	{{ .Id }}
{{ "Name:" | faint }}	{{ .Name }}
{{ "Type:" | faint }}	{{ .Type }}
{{ "Private IP Address:" | faint }}	{{ .PrivateIpAddress }}
`,
	}

	return promptui.Select{
		Label:     "Choose an instance:",
		Items:     instances,
		Templates: templates,
		Size:      20,
		Stdout:    &cliio.BellSkipper{},
	}
}

type Runner interface {
	Run() (int, string, error)
}

func askForInstance(runner Runner) int {
	accountIndex, _, err := runner.Run()

	if err != nil {
		cliio.FatalStep("Unable to select instance")
	}

	return accountIndex
}

func getUsernamePrompt(usernameOptions []string) promptui.Select {
	return promptui.Select{
		Label:  "Choose an instance:",
		Items:  usernameOptions,
		Size:   10,
		Stdout: &cliio.BellSkipper{},
	}
}

func askForUsername(runner Runner) int {
	usernameIndex, _, err := runner.Run()

	if err != nil {
		cliio.FatalStep("No username specified")
	}

	return usernameIndex
}
