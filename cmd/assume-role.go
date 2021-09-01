package cmd

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

var AssumeRoleCmd = &cobra.Command{
	Use:   "aws:assume-role",
	Short: "Allows you to enter a cli as another role",
	Run:   runAssumeRoleCmd,
}

func init() {
	rootCmd.AddCommand(AssumeRoleCmd)

	aws.RegisterCommonFlags(AssumeRoleCmd)
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

	accountPrompt := promptui.Select{
		Label:     "Choose an Account:",
		Items:     m.Accounts,
		Templates: templates,
		Size:      10,
		Stdout:    &cliio.BellSkipper{},
	}

	_, account, err := accountPrompt.Run()

	fmt.Println(account, err)

	os.Exit(0)
}
