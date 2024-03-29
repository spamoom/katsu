package assumerole

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
	"github.com/manifoldco/promptui"
	"github.com/netsells/katsu/helpers"
	"github.com/netsells/katsu/helpers/aws/iam"
	"github.com/netsells/katsu/helpers/aws/s3"
	"github.com/netsells/katsu/helpers/aws/sts"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/spf13/cobra"
)

type AccountsMeta struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Customer bool   `json:"customer"`
	Name     string `json:"name"`
	Id       string `json:"id"`
	S3Env    string `json:"s3env"`
}

type RolesMeta struct {
	Roles []Role `json:"roles"`
}

type Role struct {
	Name string `json:"name"`
}

func NewCmdAssumeRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assume-role",
		Short: "Allows you to enter a cli as another role",
		Run:   runAssumeRoleCmd,
	}

	return cmd
}

func runAssumeRoleCmd(cmd *cobra.Command, args []string) {
	helpers.SetCmd(cmd)

	accountsMeta := getAccounts()

	if len(accountsMeta.Accounts) == 0 {
		cliio.FatalStep("No accounts available")
	}

	accounts := getAccountPrompt(accountsMeta.Accounts)
	accountIndex := askForAccount(&accounts)

	rolesMeta := getRoles()

	if len(rolesMeta.Roles) == 0 {
		cliio.FatalStep("No Roles available")
	}

	roles := getRolePrompt(rolesMeta.Roles)
	roleIndex := askForRole(&roles)

	callerArn := sts.GetCallerArn()
	sessionUser := "unknown.user"

	if strings.Contains(callerArn, "user/") {
		arnParts := strings.Split(callerArn, "user/")
		sessionUser = arnParts[1]
	}

	_, err := sts.AssumeRole(sts.AssumeRoleInput{
		AccountId:   accountsMeta.Accounts[accountIndex].Id,
		Role:        rolesMeta.Roles[roleIndex].Name,
		SessionUser: sessionUser,
	})

	if err != nil {
		var oe *smithy.GenericAPIError
		if errors.As(err, &oe) {
			if oe.ErrorCode() == "AccessDenied" {
				// There's a high chance that MFA is required for this, let's try that.
				mfaDevice, err := getMfaDevice()

				if err != nil {
					cliio.FatalStepf("Access was denied assuming role %s. We tried to initiate an MFA session but you have no devices available for user %s.", rolesMeta.Roles[roleIndex].Name, sessionUser)
				}

				mfaCode, err := getMFACode()

				envVars, err := sts.AssumeRole(sts.AssumeRoleInput{
					AccountId:   accountsMeta.Accounts[accountIndex].Id,
					Role:        rolesMeta.Roles[roleIndex].Name,
					SessionUser: sessionUser,
					MfaDevice:   mfaDevice.SerialNumber,
					MfaCode:     &mfaCode,
				})

				cliio.Stepf("Now opening a session following you (%s) assuming the role %s on %s (%s) . Type `exit` to leave this shell.", sessionUser, rolesMeta.Roles[roleIndex].Name, accountsMeta.Accounts[accountIndex].Name, accountsMeta.Accounts[accountIndex].Id)

				// Decide which shell the user uses

				assumePrompt := fmt.Sprintf("%s:%s", sessionUser, accountsMeta.Accounts[accountIndex].Name)

				cmd := exec.Command("bash")
				cmd.Env = append(os.Environ(),
					"BASH_SILENCE_DEPRECATION_WARNING=1",
					fmt.Sprintf("PS1=\\e[32mnscli\\e[34m(%s)$\\e[39m ", assumePrompt),
					fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", envVars.AccessKeyID),
					fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", envVars.SecretAccessKey),
					fmt.Sprintf("AWS_SESSION_TOKEN=%s", envVars.SessionToken),
					fmt.Sprintf("AWS_S3_ENV=%s", accountsMeta.Accounts[accountIndex].S3Env),
				)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()

				os.Exit(0)
			} else {
				cliio.FatalStep("Unable to assume role")
			}
		} else {
			cliio.FatalStep("Unable to assume role")
		}
	}

	cliio.FatalStep("Assume role without MFA not yet implemented")
}

func getMFACode() (string, error) {
	prompt := promptui.Prompt{
		Label: "Please enter the code generated by your MFA device...",
	}

	return prompt.Run()
}

func getMfaDevice() (*types.MFADevice, error) {
	userDevices, err := iam.GetUserMfaDevices()

	if err != nil {
		return nil, err
	}

	if len(userDevices) == 0 {
		cliio.FatalStep("No MFA devices for current user")
	}

	if len(userDevices) > 1 {
		// Ask which device to use
		devicesPrompt := getMFAPrompt(userDevices)
		deviceIndex := askForMFADevice(&devicesPrompt)

		return &userDevices[deviceIndex], nil
	}

	return &userDevices[0], nil
}

func getAccounts() AccountsMeta {
	accountsFile, err := s3.GetFile("netsells-security-meta", "accounts.json")

	if err != nil {
		cliio.FatalStep("Unable to fetch Netsells Security meta data")
	}

	var accountsMeta AccountsMeta

	jsonBytes, err := ioutil.ReadAll(accountsFile.Body)

	if err != nil {
		cliio.FatalStep("Unable to process Netsells Security meta data")
	}

	err = json.Unmarshal(jsonBytes, &accountsMeta)

	if err != nil {
		cliio.FatalStep("Unable to process Netsells Security meta data")
	}

	return accountsMeta
}

func getRoles() RolesMeta {
	rolesFile, err := s3.GetFile("netsells-security-meta", "roles.json")

	if err != nil {
		cliio.FatalStep("Unable to fetch Netsells Security meta data")
	}

	var rolesMeta RolesMeta

	jsonBytes, err := ioutil.ReadAll(rolesFile.Body)

	if err != nil {
		cliio.FatalStep("Unable to process Netsells Security meta data")
	}

	err = json.Unmarshal(jsonBytes, &rolesMeta)

	if err != nil {
		cliio.FatalStep("Unable to process Netsells Security meta data")
	}

	return rolesMeta
}

type Runner interface {
	Run() (int, string, error)
}

func askForAccount(runner Runner) int {
	accountIndex, _, err := runner.Run()

	if err != nil {
		cliio.FatalStep("Unable to select account")
	}

	return accountIndex
}

func askForRole(runner Runner) int {
	index, _, err := runner.Run()

	if err != nil {
		cliio.FatalStep("Unable to select role")
	}

	return index
}

func askForMFADevice(runner Runner) int {
	index, _, err := runner.Run()

	if err != nil {
		cliio.FatalStep("Unable to select MFA device")
	}

	return index
}

func getAccountPrompt(accounts []Account) promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .Name }}?",
		Active:   "\U000027A1 {{ .Name | blue }}",
		Inactive: "  {{ .Name | white }}",
		Selected: fmt.Sprintf("%s Account: {{ .Name | faint }}", promptui.IconGood),
		// Searcher: https://github.com/manifoldco/promptui/blob/master/_examples/custom_select/main.go#L42
		Details: `
--------- Account ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Account ID:" | faint }}	{{ .Id }}
{{ "Env S3:" | faint }}	{{ .S3Env }}
`,
	}

	return promptui.Select{
		Label:     "Choose an Account:",
		Items:     accounts,
		Templates: templates,
		Size:      10,
		Stdout:    &cliio.BellSkipper{},
	}
}

func getRolePrompt(roles []Role) promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .Name }}?",
		Active:   "\U000027A1 {{ .Name | blue }}",
		Inactive: "  {{ .Name | white }}",
		Selected: fmt.Sprintf("%s Role: {{ .Name | faint }}", promptui.IconGood),
	}

	return promptui.Select{
		Label:     "Choose a role to assume:",
		Items:     roles,
		Templates: templates,
		Size:      10,
		Stdout:    &cliio.BellSkipper{},
	}
}

func getMFAPrompt(devices []types.MFADevice) promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .SerialNumber }}?",
		Active:   "\U000027A1 {{ .SerialNumber | blue }}",
		Inactive: "  {{ .SerialNumber | white }}",
		Selected: fmt.Sprintf("%s MFA Device: {{ .SerialNumber | faint }}", promptui.IconGood),
	}

	return promptui.Select{
		Label:     "Choose a MFA device to authenticate with:",
		Items:     devices,
		Templates: templates,
		Size:      10,
		Stdout:    &cliio.BellSkipper{},
	}
}
