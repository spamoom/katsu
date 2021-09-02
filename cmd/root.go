package cmd

import (
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/netsells/katsu/helpers/cliio"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	awsCmd "github.com/netsells/katsu/cmd/aws"
	dockerCmd "github.com/netsells/katsu/cmd/docker"
)

var envPrefix = "KATSU"

func Logo() string {
	myFigure := figure.NewFigure("Katsu", "doom", true)

	return strings.Join(myFigure.Slicify(), "\n")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := &cobra.Command{
		Use:   "katsu",
		Short: fmt.Sprintf("%s\n\nEasily manage apps and infrastructure", Logo()),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// Run init config on every command so we can do ENV fallbacks
			initConfig(cmd)

			cliio.ConfigureLogLevel(cmd)

			return nil
		},
	}

	cmd.PersistentFlags().IntP("verbosity", "v", 0, "Print verbose logs (0-3)")

	cmd.AddCommand(awsCmd.NewCmdAws())
	cmd.AddCommand(dockerCmd.NewCmdDocker())

	cobra.CheckErr(cmd.Execute())
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) {
	viper.SetEnvPrefix(envPrefix)

	// Pull in the config file for the particular project
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".katsu.yml")
	viper.ReadInConfig()

	viper.AutomaticEnv()

	bindFlags(cmd, viper.GetViper())
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
