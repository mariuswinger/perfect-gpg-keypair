package userinput

import (
	"strings"

	textinput "github.com/charmbracelet/bubbles/textinput"

	utils "perfect-gpg-keypair/internal/utils"
	styles "perfect-gpg-keypair/ui/styles"
)

type userInput struct {
	input         textinput.Model
	acceptedInput bool
	prompt        string
	helpMsg       string
	// TODO: valueObject
	defaultValue    string
	userInterrupt   bool
	validator       utils.ValidatorFunction
	validationError error
	err             error
}

func (u userInput) Value() string {
	return u.input.Value()
}

func (u userInput) ValidateInput() error {
	if u.validator == nil {
		return nil
	}
	return u.validator(u.input.Value())
}

func (u userInput) StyledPromptMsg() string {
	return styles.InfoStyle.Render(strings.TrimSuffix(u.prompt, "\n"))
}

func (u userInput) StyledValidationErrorMsg() string {
	if u.validationError != nil {
		return styles.ErrorStyle.Render(u.validationError.Error())
	}
	return ""
}

func (u userInput) StyledHelpMsg() string {
	helpMsg := "\n"
	if u.helpMsg != "" {
		helpMsg += strings.TrimSuffix(u.helpMsg, "\n") + "\n\n"
	}
	helpMsg += "Ctrl-C or Esc to exit\n"
	return styles.HelpStyle.Render(helpMsg)
}

func (u *userInput) HideCursor() {
	u.input.Cursor.SetMode(2)
}

func initialTextInputModel(placeholder string, char_limit int) textinput.Model {
	text := textinput.New()
	text.Cursor.Style = styles.CursorStyle
	text.PromptStyle = styles.CursorStyle
	text.Placeholder = placeholder
	text.CharLimit = char_limit
	text.Focus()
	return text
}

func NewNameInputModel() userInput {
	return userInput{
		input:           initialTextInputModel("Full Name", 50),
		prompt:          "Please enter your (real) full name:",
		userInterrupt:   false,
		validator:       utils.ValidateName,
		validationError: nil,
	}
}

func NewEmailInputModel() userInput {
	return userInput{
		input:           initialTextInputModel("Email", 50),
		prompt:          "Please enter your email address:",
		userInterrupt:   false,
		validator:       utils.ValidateEmail,
		validationError: nil,
	}
}

func NewExpiryInputModel() userInput {
	description := "Input is '<n>w|m|y', where n is an integer\nInput 0 for a keypair that never expires (NOT RECOMMENDED)\nThe default, recommended, value is '1y'"
	return userInput{
		input:           initialTextInputModel("<n>w|m|y", 4),
		prompt:          "Please specify how long the key should be valid:",
		helpMsg:         description,
		defaultValue:    "1y",
		userInterrupt:   false,
		validator:       utils.ValidateExpiry,
		validationError: nil,
	}
}

type UserInfoInputModel struct {
	Name   *userInput
	Email  *userInput
	Expiry *userInput
}

func NewUserInfoInputModel() UserInfoInputModel {
	nameModel := NewNameInputModel()
	emailModel := NewEmailInputModel()
	expiryModel := NewExpiryInputModel()
	return UserInfoInputModel{
		Name:   &nameModel,
		Email:  &emailModel,
		Expiry: &expiryModel,
	}
}

func NewPassphraseInputModel(prompt string) userInput {
	text := initialTextInputModel("Passphrase", 30)
	text.EchoMode = textinput.EchoPassword
	text.EchoCharacter = '•'
	return userInput{
		input:           text,
		prompt:          prompt,
		userInterrupt:   false,
		validator:       utils.ValidatePassphrase,
		validationError: nil,
	}
}

func (m *UserInfoInputModel) GetInput() error {
	if err := GetUserInput(m.Name); err != nil {
		return err
	}
	if err := GetUserInput(m.Email); err != nil {
		return err
	}
	if err := GetUserInput(m.Expiry); err != nil {
		return err
	}
	return nil
}

func (m *userInput) GetPassphrase() error {
	if err := GetUserInput(m); err != nil {
		return err
	}
	return nil
}
