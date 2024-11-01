package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

type ValidatorFunction func(string) error

type ValidationError struct {
	propertyName string
	msg          string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", ve.propertyName, ve.msg)
}

func InvalidNameError(msg string) error {
	return &ValidationError{"name", msg}
}

func InvalidEmailError(msg string) error {
	return &ValidationError{"email", msg}
}

func InvalidExpiryError(msg string) error {
	return &ValidationError{"expiry", msg}
}

func InvalidPassphraseError(msg string) error {
	return &ValidationError{"passphrase", msg}
}

func ValidateName(name string) error {
	if name == "" {
		return InvalidNameError("can not be empty")
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return InvalidEmailError("can not be empty")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return InvalidEmailError(strings.TrimPrefix(err.Error(), "mail: "))
	}
	return nil
}

func ValidateExpiry(expiry string) error {
	if expiry == "0" {
		return nil
	}
	re := regexp.MustCompile(`^\d{1,4}[wmy]$`)
	match := re.MatchString(expiry)
	if match {
		return nil
	}
	return InvalidExpiryError("ensure expiry is of format '<n>w|m|y'")
}

func ValidatePassphrase(passphrase string) error {
	// TODO: Allow empty?
	if passphrase == "" {
		return InvalidPassphraseError("can not be empty")
	}
	if regexp.MustCompile(`\s`).MatchString(passphrase) {
		return InvalidPassphraseError("can not contain whitespace")
	}
	return nil
}
