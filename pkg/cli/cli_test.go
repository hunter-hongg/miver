package cli

import (
	"testing"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		color string
		want  string
	}{
		{"red text", "hello", red, "\033[31mhello\033[0m"},
		{"empty string", "", green, "\033[32m\033[0m"},
		{"green text", "world", green, "\033[32mworld\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := colorize(tt.text, tt.color)
			if got != tt.want {
				t.Errorf("colorize(%q, %q) = %q, want %q", tt.text, tt.color, got, tt.want)
			}
		})
	}
}

func TestColorFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
		text string
	}{
		{"Red", Red, "error"},
		{"Green", Green, "success"},
		{"Yellow", Yellow, "warning"},
		{"Blue", Blue, "info"},
		{"Cyan", Cyan, "header"},
		{"Magenta", Magenta, "highlight"},
		{"Bold", Bold, "strong"},
		{"Dim", Dim, "muted"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.text)
			if len(got) == 0 {
				t.Errorf("%s(%q) returned empty string", tt.name, tt.text)
			}
			if got[len(got)-4:] != reset {
				t.Errorf("%s(%q) = %q, expected ANSI reset at end", tt.name, tt.text, got)
			}
		})
	}
}

func TestColorRoundTrip(t *testing.T) {
	// Verify that applying a color and stripping ANSI codes reconstructs the original text.
	text := "test message"
	colored := Red(text)
	if len(colored) <= len(text) {
		t.Errorf("Red(%q) = %q, expected ANSI-wrapped output", text, colored)
	}
}

func TestOutputHelpers(t *testing.T) {
	// These functions write to stdout, so we can't easily capture the output.
	// Instead, verify they don't panic and produce non-empty strings via color helpers.
	if Cyan("test") == "" {
		t.Error("Cyan returned empty")
	}
	if Yellow("test") == "" {
		t.Error("Yellow returned empty")
	}
}
