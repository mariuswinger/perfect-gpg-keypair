package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	styles "perfect-gpg-keypair/ui/styles"
)

func PrintlnStyled(message string, style lipgloss.Style) {
	fmt.Println(style.Render(message))
}

func PrintfStyled(format string, message string, style lipgloss.Style) {
	fmt.Printf(format, style.Render(message))
}

func Print(message string) {
	PrintlnStyled(message, styles.RegularStyle)
}

func PrintHiddenBorder(message string) {
	PrintlnStyled(message, styles.HiddenBorder)
}

func ErrorPrint(message string) {
	PrintlnStyled(fmt.Sprintf("ERROR %s", message), styles.ErrorStyle)
}

func InfoPrint(message string) {
	PrintlnStyled(message, styles.InfoStyle)
}

func WarningPrint(message string) {
	PrintlnStyled(message, styles.WarningStyle)
}
