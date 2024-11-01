package spinner

import (
	"fmt"
	"perfect-gpg-keypair/internal/utils"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	checkMarkIcon = '\U00002705'
	bigXIcon      = '\U0000274C'
)

type ActionCompleteSpinnerMsg string

type SpinnerErrMsg error

type SpinnerModel struct {
	title         string
	action        tea.Cmd
	spinner       spinner.Model
	userInterrupt bool
	isComplete    bool
	actionOutput  string
	error         error
}

func (m *SpinnerModel) ActionOutput() string {
	return m.actionOutput
}

func NewSpinnerModel(title string, action tea.Cmd) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return SpinnerModel{spinner: s, title: title, action: action}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Got key input:
	case tea.KeyMsg:
		switch msg.Type {

		// Exit program:
		case tea.KeyEsc, tea.KeyCtrlC:
			m.userInterrupt = true
			return m, tea.Quit

		// Do nothing for any other keys:
		default:
			return m, nil
		}

	// Got an error, return error and quit:
	case SpinnerErrMsg:
		m.error = msg
		return m, tea.Quit

	// Action is done, set output and return:
	case ActionCompleteSpinnerMsg:
		m.isComplete = true
		m.actionOutput = string(msg)
		return m, tea.Quit

	// "increment" the spinner:
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	// Starts action at first "iteration"
	default:
		return m, m.action
	}
}

func (m SpinnerModel) View() string {
	if m.error != nil {
		return fmt.Sprintf("%c %s\n", bigXIcon, m.title)
	}
	str := fmt.Sprintf("%s %s", m.spinner.View(), m.title)
	if m.userInterrupt {
		return str + "\n"
	}
	if m.isComplete {
		return fmt.Sprintf("%c %s\n", checkMarkIcon, m.title)
	}
	return str
}

func Spinner(m *SpinnerModel) error {
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return err
	}
	if m.userInterrupt {
		return &utils.UserInterrupt{}
	}
	if m.error != nil {
		return m.error
	}
	return nil
}
