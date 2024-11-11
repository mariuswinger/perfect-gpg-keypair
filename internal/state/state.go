package state

import (
	"errors"
	"fmt"
	"perfect-gpg-keypair/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	logger "github.com/sirupsen/logrus"

	userinfo "perfect-gpg-keypair/internal/state/user_info"
	tmpdir "perfect-gpg-keypair/internal/tmp_dir"
	confirm "perfect-gpg-keypair/ui/confirm"
	spinner "perfect-gpg-keypair/ui/spinner"
	userinput "perfect-gpg-keypair/ui/user_input"
)

type State struct {
	TmpDir   tmpdir.TmpDir
	UserInfo userinfo.UserInfo
}

func NewState(debug bool) State {
	return State{
		TmpDir: tmpdir.NewTmpDir(debug),
	}
}

func (state *State) SetUserInfoFromInput() error {
	for {
		userInfoInputModel := userinput.NewUserInfoInputModel()
		if err := userInfoInputModel.GetInput(); err != nil {
			return err
		}
		userInfo := userinfo.UserInfo{
			FullName: userInfoInputModel.Name.Value(),
			Email:    userInfoInputModel.Email.Value(),
			Expiry:   userInfoInputModel.Expiry.Value(),
		}

		utils.InfoPrint("You have entered:")
		utils.PrintHiddenBorder(userInfo.String())

		confirmed, err := confirm.Confirm("Is the entered information correct?")
		// TODO: Handle ctrl+c inside here? Causes panic
		if err != nil {
			return err
		}
		if confirmed {
			state.UserInfo = userInfo
			break
		}
	}
	return nil
}

func (state State) GenerateKeys() error {
	utils.InfoPrint("The master keypair will be protected by a passphrase")
	utils.WarningPrint("Ensure that you keep this passphrase in a safe space (e.g. a key vault)!")
	passphrase, err := GetPassphrase()
	if err != nil {
		return err
	}

	// Generate master keypair
	logger.Debugf("generating master keypair\n")
	generateMasterKeypairSpinner := spinner.NewSpinnerModel(
		"Generating an RSA keypair with 4096 bit size ...",
		generateMasterKeypair(state, passphrase),
	)
	err = spinner.Spinner(&generateMasterKeypairSpinner)
	if err != nil {
		return err
	}
	masterFingerprint := generateMasterKeypairSpinner.ActionOutput()
	logger.Debugf("successfully generated master keypair with fingerprint: %s\n", masterFingerprint)

	// Add signing subkey:
	logger.Debugf("adding additional signing subkey for use with this computer\n")
	addSigningSubkeySpinner := spinner.NewSpinnerModel(
		"Adding additional signing subkey for use with this computer ...",
		addSigningSubkey(state, passphrase, masterFingerprint),
	)
	err = spinner.Spinner(&addSigningSubkeySpinner)
	if err != nil {
		return err
	}
	logger.Debugf("successfully added a signing subkey\n")

	// Create revocation certificate
	logger.Debugf("creating revocation certificate at: '%s'\n", state.TmpDir.RevocationCertFilePath())
	createRevCertSpinner := spinner.NewSpinnerModel(
		"Creating revocation certificate ...",
		createRevocationCertificate(state, passphrase, masterFingerprint),
	)
	err = spinner.Spinner(&createRevCertSpinner)
	if err != nil {
		return err
	}
	logger.Debugf("exported revocation certificate to: %s\n", state.TmpDir.RevocationCertFilePath())

	// export gpg keys to temporary files
	logger.Debugf("exporting gpg keys to: %s\n", state.TmpDir.ExportedKeysDirPath())
	exportKeysSpinner := spinner.NewSpinnerModel(
		"Exporting gpg keys ...",
		exportGpgKeys(state, passphrase, masterFingerprint),
	)
	err = spinner.Spinner(&exportKeysSpinner)
	if err != nil {
		return err
	}
	logger.Debugf(fmt.Sprintf("files exported to: %s\n", state.TmpDir.ExportedKeysDirPath()))

	// confirm keys are backed up
	utils.InfoPrint(fmt.Sprintf("Files exported to: %s", state.TmpDir.ExportedKeysDirPath()))
	utils.WarningPrint("Ensure that these files are backed up (e.g. in a key vault)!\nThey will automatically be deleted after confirming they are backed up.")
	for {
		backedUp, err := confirm.Confirm("Have you backed up the files?")
		if err != nil {
			return err
		}
		if backedUp {
			break
		}
	}

	// reimport and use only secret key on this laptop:
	logger.Debugf("removing master keys and reimporting signing subkey\n")
	reimportKeySpinner := spinner.NewSpinnerModel(
		"Removing master keypair and reimport signing subkey ...",
		removeMasterAndImportSubkey(state, passphrase, masterFingerprint),
	)
	err = spinner.Spinner(&reimportKeySpinner)
	if err != nil {
		return err
	}
	logger.Debugf("successfully added a signing subkey\n")

	utils.InfoPrint("\nYour generated GPG keypair is:")
	utils.ListKeys(true, "long", masterFingerprint)
	utils.InfoPrint(fmt.Sprintf("Ensure that the key with SC attributes and the fingerprint '%s' is prepended by 'sec#'\n", masterFingerprint))
	return nil
}

func generateMasterKeypair(state State, passphrase string) tea.Cmd {
	return func() tea.Msg {
		err := utils.GenerateMasterKeypair(passphrase, state.TmpDir.StatusFilePath(), state.TmpDir.ParametersFilePath())
		if err != nil {
			return spinner.SpinnerErrMsg(fmt.Errorf("could not generate master keypair: %w", err))
		}

		masterFingerprint, err := state.TmpDir.ReadStatusFileKeyId()
		if err != nil {
			return spinner.SpinnerErrMsg(fmt.Errorf("could not read id from status file: %w", err))
		}
		return spinner.ActionCompleteSpinnerMsg(masterFingerprint)
	}
}

func addSigningSubkey(state State, passphrase string, masterFingerprint string) tea.Cmd {
	return func() tea.Msg {
		err := utils.AddSigningSubKey(passphrase, masterFingerprint, state.UserInfo.Expiry)
		if err != nil {
			return spinner.SpinnerErrMsg(fmt.Errorf("could not add signing subkey: %w", err))
		}
		return spinner.ActionCompleteSpinnerMsg("")
	}
}

func createRevocationCertificate(state State, passphrase string, masterFingerprint string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CreateRevocationCertificate(state.TmpDir.Path(), passphrase, state.TmpDir.RevocationCertFilePath(), masterFingerprint)
		if err != nil {
			return spinner.SpinnerErrMsg(fmt.Errorf("could not generate revocation certificate: %s\n", err.Error()))
		}
		return spinner.ActionCompleteSpinnerMsg("")
	}
}

func exportGpgKeys(state State, passphrase string, masterFingerprint string) tea.Cmd {
	privateMasterKeyFilePath := state.TmpDir.PrivateMasterKeyFilePath()
	publicMasterKeyFilePath := state.TmpDir.PublicMasterKeyFilePath()
	signingSubkeyFilePath := state.TmpDir.SigningSubkeyFilePath()
	return func() tea.Msg {
		privateKeyExportError := utils.ExportPrivateMasterKey(passphrase, masterFingerprint, privateMasterKeyFilePath)
		if privateKeyExportError != nil {
			logger.Debugln(fmt.Sprintf("could not export private master key: %s\n", privateKeyExportError.Error()))
		}
		publicKeyExportError := utils.ExportPublicMasterKey(masterFingerprint, publicMasterKeyFilePath)
		if publicKeyExportError != nil {
			logger.Debugln(fmt.Sprintf("could not export public master key: %s\n", publicKeyExportError.Error()))
		}
		subkeyExportError := utils.ExportSigningSubkey(passphrase, masterFingerprint, signingSubkeyFilePath)
		if subkeyExportError != nil {
			logger.Debugln(fmt.Sprintf("could not export signing subkey: %s\n", privateKeyExportError.Error()))
		}
		if privateKeyExportError != nil || publicKeyExportError != nil || subkeyExportError != nil {
			return spinner.SpinnerErrMsg(errors.New("could not export all GPG keys. This must be done manually"))
		}

		return spinner.ActionCompleteSpinnerMsg("")
	}
}

func removeMasterAndImportSubkey(state State, passphrase string, masterFingerprint string) tea.Cmd {
	return func() tea.Msg {
		err := utils.DeleteSecretKeys(passphrase, masterFingerprint)
		if err != nil {
			return spinner.SpinnerErrMsg(fmt.Errorf("could not delete master key: %w", err))
		}

		err = utils.ImportKey(passphrase, state.TmpDir.SigningSubkeyFilePath())
		if err != nil {
			return spinner.SpinnerErrMsg(fmt.Errorf("could not import signing subkey: %w", err))
		}
		return spinner.ActionCompleteSpinnerMsg("")
	}
}
