package cmd

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logLevel string
	verbose  bool
	version  = "0.0.1"
	rootCmd  = &cobra.Command{
		Use:          "perfect-gpg-keypair",
		Version:      version,
		Short:        "perfect-gpg-keypair is a simple CLI script for generating a super secure GPG keypair with a separate signing subkey",
		SilenceUsage: true,
	}
)

func init() {
	// global flags
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", logrus.InfoLevel.String(), "log level (debug, info, warn, error, fatal, panic")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	// Initialize logger
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogger(os.Stdout, logLevel); err != nil {
			return err
		}
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	// Hide auto generated 'Completion' subcommand:
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// add subcommands
	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewListCmd())
	rootCmd.AddCommand(NewRemoveCmd())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setUpLogger(out io.Writer, level string) error {
	logrus.SetOutput(out)
	// validate level string:
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	// set format
	Formatter := new(logrus.TextFormatter)
	Formatter.DisableTimestamp = true
	Formatter.DisableLevelTruncation = true
	logrus.SetFormatter(Formatter)
	return nil
}
