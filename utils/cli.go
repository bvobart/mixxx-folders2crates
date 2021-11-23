package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/mattn/go-isatty"
)

var (
	ErrAbort          = errors.New("user said no")
	ErrNotInteractive = errors.New("not an interactive terminal: cannot ask user confirmation")
)

// PromptConfirm will ask the user for confirmation
func PromptConfirm(format string, args ...interface{}) error {
	if !IsInteractive() {
		return ErrNotInteractive
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf(format, args...),
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err == promptui.ErrAbort {
		return ErrAbort
	}
	return err
}

// IsInteractive returns true if `deb` is currently being executed in a terminal,
// thus probably with a human interacting with it.
func IsInteractive() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) // linux and macos, might fail on windows, but fuck windows anyways
}
