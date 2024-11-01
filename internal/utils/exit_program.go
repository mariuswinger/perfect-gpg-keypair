package utils

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

func ExitProgram(msg string) {
	if msg == "" {
		os.Exit(1)
	}
	logger.Fatalln(msg)
}
