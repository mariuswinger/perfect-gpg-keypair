package confirm

import (
	"fmt"

	huh "github.com/charmbracelet/huh"

	styles "perfect-gpg-keypair/ui/styles"
)

func createTheme() *huh.Theme {
	theme := huh.ThemeCharm()
	theme.Focused.Title = styles.RegularStyle.Margin(0, 0, 0, 1).Bold(true)
	theme.Focused.FocusedButton = styles.FocusedButton.Margin(0, 1).Padding(0, 3)
	theme.Focused.BlurredButton = styles.UnfocusedButton.Margin(0, 1).Padding(0, 3)
	return theme
}

func createConfirm(prompt string, choice *bool) *huh.Confirm {
	return huh.NewConfirm().
		Affirmative("yes").
		Negative("no").
		Title(prompt).
		Value(choice)
}

func createThemedConfirmForm(prompt string, choice *bool) *huh.Form {
	confirm_form := huh.NewForm(huh.NewGroup(createConfirm(prompt, choice)))
	return confirm_form.WithTheme(createTheme())
	// WithShowHelp(o.ShowHelp).
}

func Confirm(prompt string) (bool, error) {
	choice := true
	confirm_form := createThemedConfirmForm(prompt, &choice)

	if err := confirm_form.Run(); err != nil {
		return false, fmt.Errorf("Unable to confirm user info: %w", err)
	}

	return choice, nil
}
