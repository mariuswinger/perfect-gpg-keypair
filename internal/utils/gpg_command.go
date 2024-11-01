package utils

import (
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

type GpgCommandArgs struct {
	args []string
}

func NewGpgCommand(subcommand string) GpgCommandArgs {
	return GpgCommandArgs{[]string{subcommand}}
}

func (c GpgCommandArgs) hasPassphrase() bool {
	return slices.Contains(c.args, "--passphrase")
}

func (c GpgCommandArgs) append(a ...string) GpgCommandArgs {
	return GpgCommandArgs{append(c.args, a...)}
}

func (c GpgCommandArgs) addFlag(flag string) GpgCommandArgs {
	return c.append(flag)
}

func (c GpgCommandArgs) addArg(arg string) GpgCommandArgs {
	if arg != "" {
		return c.append(arg)
	}
	return c
}

func (c GpgCommandArgs) addOption(option string, arg string) GpgCommandArgs {
	return c.append([]string{option, arg}...)
}

func (c GpgCommandArgs) addPassphrase(passphrase string) GpgCommandArgs {
	return c.addOption("--pinentry", "loopback").addOption("--passphrase", passphrase)
}

func (c GpgCommandArgs) addOutput(outputFilepath string) GpgCommandArgs {
	return c.addOption("--output", outputFilepath)
}

func (c GpgCommandArgs) toCommand() *exec.Cmd {
	return exec.Command("gpg", c.args...)
}

func (c GpgCommandArgs) getCommandString() string {
	out := "gpg " + strings.Join(c.args, " ")
	if c.hasPassphrase() {
		re := regexp.MustCompile(`(--passphrase) ([^\s]+)`)
		return re.ReplaceAllString(out, `$1 XXXXX`)
	}
	return out
}
