package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
)

func CheckGpgIsInstalled() error {
	logger.Debugln("running: 'which gpg'")
	cmd := exec.Command("which", "gpg")
	_, err := cmd.Output()
	return err
}

func ListKeys(secret bool, format string, name string) error {
	var c GpgCommandArgs
	if secret {
		c = NewGpgCommand("--list-secret-keys").addOption("--keyid-format", format).addArg(name)
	} else {
		c = NewGpgCommand("--list-keys").addOption("--keyid-format", format).addArg(name)
	}
	cmd := c.toCommand()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	return cmd.Run()
}

func DeleteEntireKey(fingerprint string) error {
	c := NewGpgCommand("--delete-secret-and-public-keys").addFlag("--batch").addFlag("--yes").addArg(fingerprint)
	cmd := c.toCommand()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	return cmd.Run()
}

func DeleteSecretKeys(passphrase string, fingerprint string) error {
	c := NewGpgCommand("--delete-secret-keys").addFlag("--batch").addFlag("--yes").addPassphrase(passphrase).addArg(fingerprint)
	cmd := c.toCommand()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	return cmd.Run()
}

func GenerateMasterKeypair(passphrase string, statusFilepath string, parametersFilepath string) error {
	c := NewGpgCommand("--generate-key").addArg("--no-tty").addArg("--batch").addPassphrase(passphrase).addOption("--status-file", statusFilepath).addArg(parametersFilepath)
	cmd := c.toCommand()
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	_, err := cmd.Output()
	return err
}

func AddSigningSubKey(passphrase string, masterKeyId string, expiry string) error {
	fingerprint, err := getKeyFingerprint(masterKeyId)
	if err != nil {
		return fmt.Errorf("could not get fingerprint for key: %w", err)
	}
	c := NewGpgCommand("--quick-add-key").addArg("--no-tty").addArg("--batch").addPassphrase(passphrase).addArg(fingerprint).addArg("rsa4096").addArg("sign").addArg(expiry)
	cmd := c.toCommand()
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	_, err = cmd.Output()
	return err
}

func CreateRevocationCertificate(tmpDir string, passphrase string, outputFilepath string, masterKeyId string) error {
	commandFilePath := filepath.Join(tmpDir, ".rev-cert-input")
	err := createRevocationCertificateCommandFile(commandFilePath)
	if err != nil {
		return fmt.Errorf("could not create input file: %w", err)
	}
	defer os.Remove(commandFilePath)
	c := NewGpgCommand("--gen-revoke").addArg("--no-tty").addPassphrase(passphrase).addOption("--command-file", commandFilePath).addOutput(outputFilepath).addArg(masterKeyId)
	cmd := c.toCommand()
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	_, err = cmd.Output()
	return err
}

func getKeyFingerprint(keyId string) (string, error) {
	c := NewGpgCommand("--fingerprint").addArg(keyId)
	cmd := c.toCommand()

	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(out) < 2 {
		return "", fmt.Errorf("failed to parse gpg output")
	}
	fingerprintLine := strings.Split(string(out), "\n")[1]
	return strings.Join(strings.Fields(fingerprintLine), " "), nil
}

func createRevocationCertificateCommandFile(fp string) error {
	commandFileContents := "y\n1\n\ny\n"
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(commandFileContents)
	if err != nil {
		return err
	}
	return nil
}

func ExportPublicMasterKey(masterKeyId string, outputFilepath string) error {
	c := NewGpgCommand("--export").addArg("--armor").addOutput(outputFilepath).addArg(masterKeyId)
	cmd := c.toCommand()
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func ExportPrivateMasterKey(passphrase string, masterKeyId string, outputFilepath string) error {
	c := NewGpgCommand("--export-secret-keys").addArg("--armor").addPassphrase(passphrase).addOutput(outputFilepath).addArg(masterKeyId)
	cmd := c.toCommand()
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func ExportSigningSubkey(passphrase string, masterKeyId string, outputFilepath string) error {
	c := NewGpgCommand("--export-secret-subkeys").addArg("--armor").addPassphrase(passphrase).addOutput(outputFilepath).addArg(masterKeyId)
	cmd := c.toCommand()
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func ImportKey(passphrase string, filePath string) error {
	c := NewGpgCommand("--import").addPassphrase(passphrase).addArg(filePath)
	cmd := c.toCommand()
	logger.Debugln(fmt.Sprintf("running: '%s'\n", c.getCommandString()))
	_, err := cmd.Output()
	return err
}
