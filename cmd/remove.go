package cmd

import (
	"fmt"
	"perfect-gpg-keypair/internal/utils"
	confirm "perfect-gpg-keypair/ui/confirm"

	"github.com/spf13/cobra"
)

func NewRemoveCmd() *cobra.Command {
	var force bool
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "remove existing GPG key by fingerprint",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fingerprint := args[0]

			err := validateFingerprint(fingerprint)
			if err != nil {
				utils.ExitProgram("invalid fingerprint: " + err.Error())
			}

			err = remove(fingerprint, force)
			if err != nil {
				utils.ExitProgram("could not delete gpg key: " + err.Error())
			}
			utils.InfoPrint(fmt.Sprintf("successfully removed key '%s'", fingerprint))
		},
	}

	// add flags
	deleteCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "forcefully remove GPG key")
	return deleteCmd
}

func remove(fingerprint string, force bool) error {
	if !force {
		confirmMsg := fmt.Sprintf("Are you really sure you want to delete the key with fingerprint '%s'", fingerprint)
		confirmDelete, err := confirm.Confirm(confirmMsg)
		if err != nil {
			return err
		}
		if !confirmDelete {
			return nil
		}
	}
	return utils.DeleteEntireKey(fingerprint)
}

func validateFingerprint(fingerprint string) error {
	if len(fingerprint) != 40 {
		return fmt.Errorf("must have length of 40")
	}
	return nil
}
