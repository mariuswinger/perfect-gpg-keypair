package cmd

import (
	"perfect-gpg-keypair/internal/utils"

	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	var longFormat bool
	var secret bool
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all existing public GPG keys",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			format := getFormat(longFormat)
			if err := utils.ListKeys(secret, format, ""); err != nil {
				utils.ExitProgram("Failed to list gpg keys: " + err.Error())
			}
		},
	}

	// add flags
	listCmd.PersistentFlags().BoolVar(&longFormat, "long", false, "use long key format")
	listCmd.PersistentFlags().BoolVar(&secret, "secret", false, "list secret keys")
	return listCmd
}

func getFormat(longFormat bool) string {
	if longFormat {
		return "long"
	}
	return "short"
}
