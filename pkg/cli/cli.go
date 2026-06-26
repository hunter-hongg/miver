// Package cli provides unified CLI output formatting, colors, and error handling.
package cli

import (
	"fmt"
	"os"
)

// ANSI escape codes
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	faint   = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

// Symbols used in output (ASCII-only for broad terminal compatibility)
const (
	symSuccess = "*"
	symInfo    = "i"
	symWarning = "!"
	symError   = "x"
	symStep    = ">"
)

func colorize(text, color string) string {
	return color + text + reset
}

func Bold(text string) string    { return bold + text + reset }
func Red(text string) string     { return colorize(text, red) }
func Green(text string) string   { return colorize(text, green) }
func Yellow(text string) string  { return colorize(text, yellow) }
func Blue(text string) string    { return colorize(text, blue) }
func Cyan(text string) string    { return colorize(text, cyan) }
func Magenta(text string) string { return colorize(text, magenta) }

func Dim(text string) string {
	return faint + text + reset
}

func Warning(text string) {
	fmt.Printf(" %s  %s\n", Yellow(symWarning), text)
}

func Error(text string) {
	fmt.Printf(" %s  %s\n", Red(symError), text)
}

func Info(text string) {
	fmt.Printf(" %s  %s\n", Blue(symInfo), text)
}

func Success(text string) {
	fmt.Printf(" %s  %s\n", Green(symSuccess), text)
}

func Step(text string) {
	fmt.Printf("\n %s  %s\n", Cyan(symStep), Bold(text))
}

// Fatal prints an error message and exits with code 1.
func Fatal(err error) {
	if err != nil {
		Error(err.Error())
		os.Exit(1)
	}
}
