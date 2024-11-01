package state

import (
	utils "perfect-gpg-keypair/internal/utils"
	userinput "perfect-gpg-keypair/ui/user_input"
)

func GetPassphrase() (string, error) {
	for {
		firstPassphraseInputModel := userinput.NewPassphraseInputModel("Please enter a passphrase:")
		if err := userinput.GetUserInput(&firstPassphraseInputModel); err != nil {
			return "", err
		}
		passphrase := firstPassphraseInputModel.Value()

		confirmPassphraseInputModel := userinput.NewPassphraseInputModel("Please re-enter the passphrase to confirm:")
		if err := userinput.GetUserInput(&confirmPassphraseInputModel); err != nil {
			return "", err
		}

		if !(passphrase == confirmPassphraseInputModel.Value()) {
			utils.ErrorPrint("Passphrases do not match!")
		} else {
			return passphrase, nil
		}
	}
}
