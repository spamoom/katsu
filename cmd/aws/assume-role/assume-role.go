package assumerole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/netsells/katsu/helpers"
	"github.com/netsells/katsu/helpers/aws"
	"github.com/netsells/katsu/helpers/aws/s3"
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

func NewCmdAssumeRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assume-role",
		Short: "Allows you to enter a cli as another role",
		Run:   runAssumeRoleCmd,
	}

	aws.RegisterCommonFlags(cmd)

	return cmd
}

func runAssumeRoleCmd(cmd *cobra.Command, args []string) {
	helpers.SetCmd(cmd)

	accountsFile, err := s3.GetFile("netsells-security-meta", "accounts.json")

	if err != nil {
		cliio.FatalStep("Unable to fetch Netsells Security meta data")
	}

	var m AccountsMeta

	jsonBytes, err := ioutil.ReadAll(accountsFile.Body)

	if err != nil {
		cliio.FatalStep("Unable to process Netsells Security meta data")
	}

	err = json.Unmarshal(jsonBytes, &m)

	if err != nil {
		cliio.FatalStep("Unable to process Netsells Security meta data")
	}

	if len(m.Accounts) == 0 {
		cliio.FatalStep("No accounts available")
	}

	accounts := getAccountPrompt(m.Accounts)
	accountIndex := askForAccount(&accounts)

	fmt.Println(accountIndex, err, m.Accounts[accountIndex].Name)

	os.Exit(0)
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

func getAccountPrompt(accounts []Account) promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .Name }}?",
		Active:   "\U000027A1 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | white }}",
		Selected: "\U000027A1 {{ .Name | red | cyan }}",
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
