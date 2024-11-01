package userinput

import (
	"fmt"
	"perfect-gpg-keypair/internal/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type UserInputErrMsg error

func (model userInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m *userInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Switch on message type
	switch msg := msg.(type) {

	// Got key input:
	case tea.KeyMsg:
		switch msg.Type {

		// Exit program
		case tea.KeyEsc, tea.KeyCtrlC:
			m.userInterrupt = true
			m.HideCursor()
			return m, tea.Quit

		// Enter input and check if valid
		case tea.KeyEnter:
			// Set default if default is present and value is empty
			if m.input.Value() == "" && m.defaultValue != "" {
				m.input.SetValue(m.defaultValue)
			}
			if err := m.ValidateInput(); err != nil {
				m.acceptedInput = false
				m.validationError = err
				return m, nil
			} else {
				m.acceptedInput = true
				m.validationError = nil
				m.HideCursor()
				return m, tea.Quit
			}

		// User entered text, update textinput:
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

	// Got an error, return and quit:
	case UserInputErrMsg:
		m.err = msg
		return m, tea.Quit

	// Keep blinking cursor:
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (model userInput) View() string {
	out := fmt.Sprintf("%s\n%s", model.StyledPromptMsg(), model.input.View())
	if model.userInterrupt || model.acceptedInput {
		return out + "\n"
	}
	if model.validationError != nil {
		out += "\n" + model.StyledValidationErrorMsg()
	}
	out += "\n" + model.StyledHelpMsg()
	return out
}

func GetUserInput(model *userInput) error {
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		return err
	}
	if model.userInterrupt {
		return &utils.UserInterrupt{}
	}
	if model.err != nil {
		return model.err
	}
	return nil
}
