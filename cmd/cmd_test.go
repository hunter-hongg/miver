package cmd

import (
	"testing"
)

func TestCheckPlatform(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		wantErr bool
		errMsg  string
	}{
		{"linux is supported", "linux", false, ""},
		{"windows not supported", "windows", true, "Windows is not supported. Please use WSL instead"},
		{"macOS not supported", "darwin", true, "macOS is not supported. Please use Linux instead"},
		{"unknown os not supported", "freebsd", true, "unsupported system: freebsd"},
		{"empty string", "", true, "unsupported system: "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPlatform(tt.goos)
			if tt.wantErr {
				if err == nil {
					t.Errorf("checkPlatform(%q) expected error, got nil", tt.goos)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("checkPlatform(%q) error = %q, want %q", tt.goos, err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("checkPlatform(%q) returned unexpected error: %v", tt.goos, err)
				}
			}
		})
	}
}
