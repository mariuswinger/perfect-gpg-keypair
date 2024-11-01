package tmpdir

import (
	"fmt"
	"os"

	userinfo "perfect-gpg-keypair/internal/state/user_info"
)

type ParametersFile struct {
	Path string
}

func (parameters_file ParametersFile) Create(user_info userinfo.UserInfo) error {
	f, err := os.Create(parameters_file.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(parameters_file.contents(user_info))
	if err != nil {
		return err
	}
	return nil
}

func (parameters_file ParametersFile) contents(user_info userinfo.UserInfo) string {
	return fmt.Sprintf(
		"Key-Type: RSA\n"+
			"Key-Length: 4096\n"+
			"Key-Usage: sign\n"+
			"Subkey-Type: RSA\n"+
			"Subkey-Length: 406\n"+
			"Subkey-Usage: encrypt\n"+
			"Name-Real: %s\n"+
			"Name-Email: %s\n"+
			"Expire-Date: %s\n"+
			"Preferences: SHA512 SHA384 SHA256 SHA224 AES256 AES192 AES CAST5 ZLIB BZIP2 ZIP Uncompressed\n"+
			"%%commit",
		user_info.FullName,
		user_info.Email,
		user_info.Expiry,
	)
}
