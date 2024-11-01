package styles

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	Black       = lipgloss.Color("#000000")
	White       = lipgloss.Color("#FFFFFF")
	Gray        = lipgloss.Color("240")
	ZambeziGray = lipgloss.Color("#585858")
	Red         = lipgloss.Color("1")
	Redder      = lipgloss.Color("#AA0000")
	Yellow      = lipgloss.Color("3")
	Blue        = lipgloss.Color("6")
	MayaBlue    = lipgloss.Color("#89B4FA")
	DarkOrange  = lipgloss.Color("#FF8700")
)

var (
	// text styles
	RegularStyle = lipgloss.NewStyle()
	FocusedStyle = lipgloss.NewStyle().Foreground(DarkOrange)
	HelpStyle    = lipgloss.NewStyle().Foreground(Gray)
	// InfoStyle         = lipgloss.NewStyle().Foreground(MayaBlue)
	InfoStyle         = lipgloss.NewStyle().Foreground(Blue)
	WarningStyle      = lipgloss.NewStyle().Foreground(Yellow).Bold(true).Underline(true).Margin(1, 0).Padding(0, 4)
	ErrorStyle        = lipgloss.NewStyle().Foreground(Red)
	InvalidInputStyle = lipgloss.NewStyle().Foreground(Red)

	// other
	CursorStyle     = lipgloss.NewStyle().Foreground(Yellow)
	FocusedButton   = lipgloss.NewStyle().Background(Yellow).Foreground(Black)
	UnfocusedButton = lipgloss.NewStyle().Background(Black).Foreground(White)
	HiddenBorder    = lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder()).Margin(0, 1)
)
