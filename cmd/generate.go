package cmd

import (
	"errors"
	"fmt"
	"os"
	"perfect-gpg-keypair/internal/utils"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	state "perfect-gpg-keypair/internal/state"
)

func NewGenerateCmd() *cobra.Command {
	var debug bool
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate a GPG keypair along with a separate signing subkey",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			mainState := state.NewState(debug)
			err := generate(&mainState)
			cleanup(&mainState, debug)
			if err != nil {
				handleError(err)
			}
		},
	}

	// add flags
	generateCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode (affects path of tmp dir)")
	return generateCmd
}

func cleanup(mainState *state.State, debug bool) {
	logger.Debugln("removing temporary exported keys directory")
	err := os.RemoveAll(mainState.TmpDir.ExportedKeysDirPath())
	if err != nil {
		logger.Errorf(
			"Failed to remove temporary files at '%s'."+
				"These files contain information about your keys and should be deleted "+
				"if you intend to use the generated keys!\n", mainState.TmpDir.Path(),
		)
	}
	// Keep tmp dir after run for debug runs:
	if !debug {
		logger.Debugln("removing temporary directory")
		err := os.RemoveAll(mainState.TmpDir.Path())
		if err != nil {
			logger.Errorf("could not remove temporary directory at '%s'", mainState.TmpDir.Path())
		}
	}
}

func handleError(err error) {
	interrupt := &utils.UserInterrupt{}
	if errors.As(err, &interrupt) {
		fmt.Printf("SIGINT: %s\n", err.Error())
		os.Exit(130)
	} else {
		utils.ExitProgram(err.Error())
	}
}

func generate(mainState *state.State) error {
	// TODO: Add output path argument?
	if err := utils.CheckGpgIsInstalled(); err != nil {
		return fmt.Errorf("gpg command could not be found: %w", err)
	}

	utils.InfoPrint(
		"Welcome! This program will follow you through the process of generating a secure GPG keypair.\n" +
			"The program will generate a master keypair (public and private keys) that should be stored in a safe place.\n" +
			"In addition, the program will generate a signing subkey to use for this computer.\n" +
			"Finally, the program will remove the master keypair (after ensuring they are backed up!) and import the " +
			"signing subkey so that the master key can not be obtained from this computer." +
			"At any time you can press C-c or Esc to quit the program.",
	)

	// Create temp dir
	if err := mainState.TmpDir.Create(); err != nil {
		return fmt.Errorf("could not create temporary directory: %w", err)
	}

	// Set info from user input
	utils.InfoPrint("In order to generate a GPG keypair, we need some information about you")
	if err := mainState.SetUserInfoFromInput(); err != nil {
		if err == err.(*utils.UserInterrupt) {
			return err
		} else {
			return fmt.Errorf("could not get user input: %w", err)
		}
	}

	// Write parameter file
	if err := mainState.TmpDir.CreateParametersFile(mainState.UserInfo); err != nil {
		return fmt.Errorf("could not create parameters file: %w", err)
	}

	// Generate key pair
	if err := mainState.GenerateKeys(); err != nil {
		return fmt.Errorf("could not generate GPG keys: %w", err)
	}

	utils.InfoPrint(
		"You may now want to add the signing key to your git config\n" +
			"This can be done with the following command in a git repository:\n" +
			"  'git config user.signingkey [key_id]'\n" +
			"Add '--global' to use the signing subkey globally " +
			"the key_id to use is the 16 hex digits after 'sec#  rsa4096/'",
	)
	return nil
}
