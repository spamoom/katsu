package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/manifoldco/promptui"
	"github.com/netsells/katsu/helpers"
	"github.com/netsells/katsu/helpers/aws/s3"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/netsells/katsu/helpers/config"
	"github.com/spf13/cobra"
)

type FileAction struct {
	File   *string
	Action string
}

func NewCmdManageEnv() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manage-env",
		Short: "Opens up editor for environment variables s3 files",
		Run:   runManageEnvCmd,
	}

	cmd.Flags().String("s3-bucket", "", "The S3 bucket name to look for files")

	return cmd
}

func runManageEnvCmd(cmd *cobra.Command, args []string) {
	helpers.SetCmd(cmd)

	bucketName := config.GetS3Bucket()
	if bucketName == "" {
		bucketName = os.Getenv("AWS_S3_ENV")
	}

	if bucketName == "" {
		cliio.FatalStep("An S3 bucket for the environment variable files must be specified. Calling aws:assume-role will do this automatically for you.")
	}

	cliio.Stepf("Looking in S3 bucket %s for environment files", bucketName)

	files := getFile(bucketName)

	tempFile, err := ioutil.TempFile("", "katsu-env")

	if err != nil {
		cliio.FatalStepf("Failed to write temp file")
	}

	defer os.Remove(tempFile.Name())

	if files.Action == "" {
		fileContents, err := s3.GetFile(bucketName, *files.File)

		if err != nil {
			cliio.FatalStepf("Failed to get file contents: %s", err)
		}

		bodyBytes, _ := ioutil.ReadAll(fileContents.Body)
		err = os.WriteFile(tempFile.Name(), bodyBytes, 0644)

		if err != nil {
			cliio.FatalStepf("Failed to write temp file")
		}

		cmd := exec.Command("nano", tempFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()

		newBodyBytes, _ := ioutil.ReadFile(tempFile.Name())

		edits := myers.ComputeEdits(span.URIFromPath("a.txt"), string(bodyBytes), string(newBodyBytes))

		if len(edits) == 0 {
			cliio.WarnStep("No changes")
			os.Exit(0)
		}

		diff := gotextdiff.ToUnified("old env", "new env", string(bodyBytes), edits)

		cliio.PrintDiff(diff)

		result, _ := cliio.ConfirmStep(cliio.DefaultQuestion{
			Question: "Are you sure you want to commit these changes?",
			Default:  "n",
		})

		if result == "y" {
			err = s3.PutFile(bucketName, *files.File, newBodyBytes)

			if err != nil {
				cliio.FatalStepf("Failed to commit changes: %s", err)
			}

			cliio.SuccessfulStep("Successfully committed changes")
			os.Exit(0)
		}

		cliio.WarnStep("Changes not committed")
		os.Exit(0)
	} else {
		newName, err := cliio.AskStep(cliio.DefaultQuestion{
			Question: "Enter a name for the new file",
			Default:  "",
		})

		if err != nil {
			cliio.FatalStepf("Failed to get new file name")
		}

		cmd := exec.Command("nano", tempFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()

		newBodyBytes, _ := ioutil.ReadFile(tempFile.Name())

		err = s3.PutFile(bucketName, newName, newBodyBytes)

		if err != nil {
			cliio.FatalStepf("Failed to create file: %s", err)
		}

		cliio.SuccessfulStep("Successfully created file")
		os.Exit(0)
	}

	os.Exit(0)
}

type Runner interface {
	Run() (int, string, error)
}

func getFile(bucketName string) FileAction {

	files, _ := s3.ListFiles(bucketName)

	files = append(files, "Create new file")

	filePrompt := getFilePrompt(files)
	fileIndex := askForFiles(&filePrompt)

	if files[fileIndex] == "Create new file" {
		return FileAction{
			File:   nil,
			Action: "create",
		}
	}

	return FileAction{
		File:   &files[fileIndex],
		Action: "",
	}
}

func askForFiles(runner Runner) int {
	index, _, err := runner.Run()

	if err != nil {
		os.Exit(0)
	}

	return index
}

func getFilePrompt(files []string) promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A1 {{ . | blue }}",
		Inactive: "  {{ . | white }}",
		Selected: fmt.Sprintf("%s File: {{ . | faint }}", promptui.IconGood),
	}

	return promptui.Select{
		Label:     "Choose a file to edit:",
		Items:     files,
		Templates: templates,
		Size:      10,
		Stdout:    &cliio.BellSkipper{},
	}
}
