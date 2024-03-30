package main

import (
	"errors"
	"os"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/szkiba/k6pack"
	"golang.org/x/term"
)

func formatError(err error) string {
	width, color := formatOptions(int(os.Stderr.Fd())) //nolint:forbidigo

	var perr k6pack.PackError
	if errors.As(err, &perr) {
		return perr.Format(width, color)
	}

	return strings.Join(
		api.FormatMessages(
			[]api.Message{{Text: err.Error()}},
			api.FormatMessagesOptions{TerminalWidth: width, Color: color},
		),
		"\n",
	)
}

func formatOptions(fd int) (int, bool) {
	color := false
	width := 0

	if term.IsTerminal(fd) {
		if os.Getenv("NO_COLOR") != "true" { //nolint:forbidigo
			color = true
		}

		if w, _, err := term.GetSize(fd); err == nil {
			width = w
		}
	}

	return width, color
}
