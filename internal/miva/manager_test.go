package miva

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewPaths(t *testing.T) {
	home := "/home/testuser"
	paths := NewPaths(home)

	wantHome := "/home/testuser/.miver"
	wantRepo := "/home/testuser/.miver/repo"
	wantMiva := "/home/testuser/.miver/repo/miva"
	wantBin := "/home/testuser/.miver/bin"
	wantLib := "/home/testuser/.miver/lib"
	wantLibStd := "/home/testuser/.miver/lib/.std"

	if paths.Home != wantHome {
		t.Errorf("Home = %q, want %q", paths.Home, wantHome)
	}
	if paths.Repo != wantRepo {
		t.Errorf("Repo = %q, want %q", paths.Repo, wantRepo)
	}
	if paths.Miva != wantMiva {
		t.Errorf("Miva = %q, want %q", paths.Miva, wantMiva)
	}
	if paths.Bin != wantBin {
		t.Errorf("Bin = %q, want %q", paths.Bin, wantBin)
	}
	if paths.Lib != wantLib {
		t.Errorf("Lib = %q, want %q", paths.Lib, wantLib)
	}
	if paths.LibStd != wantLibStd {
		t.Errorf("LibStd = %q, want %q", paths.LibStd, wantLibStd)
	}
}

func TestNewPathsEmptyHome(t *testing.T) {
	paths := NewPaths("")
	// With empty home, filepath.Join("", ".miver") returns a relative path on Linux
	if paths.Home != ".miver" {
		t.Errorf("with empty home, Home = %q, want %q", paths.Home, ".miver")
	}
}

func TestDirExists(t *testing.T) {
	t.Run("existing directory", func(t *testing.T) {
		dir := t.TempDir()
		exists, err := DirExists(dir)
		if err != nil {
			t.Fatalf("DirExists(%q) returned error: %v", dir, err)
		}
		if !exists {
			t.Errorf("DirExists(%q) = false, want true", dir)
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		exists, err := DirExists("/tmp/does-not-exist-12345")
		if err != nil {
			t.Fatalf("DirExists returned error: %v", err)
		}
		if exists {
			t.Error("DirExists for non-existent path = true, want false")
		}
	})

	t.Run("path is a file, not directory", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "somefile")
		if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
			t.Fatal(err)
		}
		exists, err := DirExists(f)
		if err == nil {
			t.Error("DirExists on a file: expected error, got nil")
		}
		if exists {
			t.Error("DirExists on a file: expected false, got true")
		}
	})
}

func TestAvailableVersions(t *testing.T) {
	t.Run("reads CSV correctly", func(t *testing.T) {
		dir := t.TempDir()
		csvContent := "1.0.0,alpha,2024-01-01\n1.1.0,beta,2024-02-01\n"
		csvPath := filepath.Join(dir, "available.csv")
		if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
			t.Fatal(err)
		}

		records, err := AvailableVersions(dir)
		if err != nil {
			t.Fatalf("AvailableVersions returned error: %v", err)
		}

		if len(records) != 2 {
			t.Fatalf("got %d records, want 2", len(records))
		}

		want := [][]string{
			{"1.0.0", "alpha", "2024-01-01"},
			{"1.1.0", "beta", "2024-02-01"},
		}
		for i, rec := range records {
			for j, field := range rec {
				if field != want[i][j] {
					t.Errorf("records[%d][%d] = %q, want %q", i, j, field, want[i][j])
				}
			}
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := AvailableVersions(t.TempDir())
		if err == nil {
			t.Error("AvailableVersions on empty dir: expected error, got nil")
		}
	})

	t.Run("empty CSV", func(t *testing.T) {
		dir := t.TempDir()
		csvPath := filepath.Join(dir, "available.csv")
		if err := os.WriteFile(csvPath, []byte(""), 0644); err != nil {
			t.Fatal(err)
		}

		records, err := AvailableVersions(dir)
		if err != nil {
			t.Fatalf("AvailableVersions returned error: %v", err)
		}
		if len(records) != 0 {
			t.Errorf("got %d records, want 0", len(records))
		}
	})
}

func TestInitRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping InitRepo test in short mode (requires network)")
	}

	t.Run("creates directories", func(t *testing.T) {
		home := t.TempDir()
		paths := NewPaths(home)

		// DirExists should return false before InitRepo creates dirs
		exists, err := DirExists(paths.Repo)
		if err != nil {
			t.Fatal(err)
		}
		if exists {
			t.Fatal("Repo dir should not exist yet")
		}

		// InitRepo will attempt git clone, which requires network.
		// We only verify directory creation behavior here.
		err = InitRepo(paths)
		if err != nil {
			// If git clone fails (no network), InitRepo should still create the Repo dir
			// before attempting clone. Verify that.
			if !os.IsNotExist(err) {
				// Check that at least the repo directory was created
				exists, statErr := DirExists(paths.Repo)
				if statErr != nil {
					t.Fatalf("Repo dir check failed: %v", statErr)
				}
				if !exists {
					t.Error("InitRepo should create Repo dir before cloning")
				}
			}
		}
	})
}

func TestSwitchVersion(t *testing.T) {
	t.Run("creates dangling symlinks for non-existent version", func(t *testing.T) {
		home := t.TempDir()
		paths := NewPaths(home)

		// Create a minimal miva repo structure (without the version dir)
		mivaDir := filepath.Join(paths.Miva)
		if err := os.MkdirAll(mivaDir, 0755); err != nil {
			t.Fatal(err)
		}

		// SwitchVersion does not validate the version directory; it creates dangling symlinks
		// via ln -sf. The binary link will be dangling since the source doesn't exist.
		err := SwitchVersion(paths, "9.9.9")
		if err != nil {
			t.Fatalf("SwitchVersion returned unexpected error: %v", err)
		}

		// Verify dangling symlink was created
		binLink := filepath.Join(paths.Bin, "miva")
		info, err := os.Lstat(binLink)
		if err != nil {
			t.Fatalf("expected symlink %s to exist: %v", binLink, err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			t.Errorf("%s is not a symlink", binLink)
		}
	})

	t.Run("switches to existing version", func(t *testing.T) {
		home := t.TempDir()
		paths := NewPaths(home)

		// Create a fake miva repository structure with a version
		version := "1.0.0"
		mivaVersionDir := filepath.Join(paths.Miva, version, "bin")
		if err := os.MkdirAll(mivaVersionDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create fake miva linux binary
		binContent := []byte("fake binary")
		if err := os.WriteFile(filepath.Join(mivaVersionDir, "miva-linux"), binContent, 0755); err != nil {
			t.Fatal(err)
		}

		// Create fake lib directory for std
		libDir := filepath.Join(paths.Miva, version, "lib")
		if err := os.MkdirAll(libDir, 0755); err != nil {
			t.Fatal(err)
		}
		// Create fake header files
		for _, hdr := range []string{"mvp_builtin.h", "mvp_copyable.h"} {
			if err := os.WriteFile(filepath.Join(libDir, hdr), []byte("header"), 0644); err != nil {
				t.Fatal(err)
			}
		}

		err := SwitchVersion(paths, version)
		if err != nil {
			t.Fatalf("SwitchVersion returned error: %v", err)
		}

		// Verify symlinks were created
		tests := []struct {
			name string
			path string
		}{
			{"miva binary", filepath.Join(paths.Bin, "miva")},
			{"std lib", paths.LibStd},
			{"mvp_builtin.h header", filepath.Join(paths.Lib, "mvp_builtin.h")},
			{"mvp_copyable.h header", filepath.Join(paths.Lib, "mvp_copyable.h")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				info, err := os.Lstat(tt.path)
				if err != nil {
					t.Fatalf("expected %s to exist: %v", tt.path, err)
				}
				if info.Mode()&os.ModeSymlink == 0 {
					t.Errorf("%s is not a symlink", tt.path)
				}
			})
		}
	})

	t.Run("switches to version with frontend binary", func(t *testing.T) {
		home := t.TempDir()
		paths := NewPaths(home)

		version := "2.0.0"
		binDir := filepath.Join(paths.Miva, version, "bin")
		if err := os.MkdirAll(binDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create fake binaries including frontend
		for _, b := range []string{"miva-linux", "miva-frontend-linux"} {
			if err := os.WriteFile(filepath.Join(binDir, b), []byte("fake"), 0755); err != nil {
				t.Fatal(err)
			}
		}

		libDir := filepath.Join(paths.Miva, version, "lib")
		if err := os.MkdirAll(libDir, 0755); err != nil {
			t.Fatal(err)
		}
		for _, hdr := range []string{"mvp_builtin.h", "mvp_copyable.h"} {
			if err := os.WriteFile(filepath.Join(libDir, hdr), []byte("header"), 0644); err != nil {
				t.Fatal(err)
			}
		}

		err := SwitchVersion(paths, version)
		if err != nil {
			t.Fatalf("SwitchVersion returned error: %v", err)
		}

		// Verify frontend symlink exists
		frontendPath := filepath.Join(paths.Bin, "miva-frontend")
		info, err := os.Lstat(frontendPath)
		if err != nil {
			t.Fatalf("expected miva-frontend symlink to exist: %v", err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			t.Errorf("miva-frontend is not a symlink")
		}
	})

	t.Run("removes old lib before switching", func(t *testing.T) {
		home := t.TempDir()
		paths := NewPaths(home)

		// Create old lib with stale content
		oldLibFile := filepath.Join(paths.Lib, "old-file")
		if err := os.MkdirAll(paths.Lib, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(oldLibFile, []byte("stale"), 0644); err != nil {
			t.Fatal(err)
		}

		version := "3.0.0"
		binDir := filepath.Join(paths.Miva, version, "bin")
		if err := os.MkdirAll(binDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(binDir, "miva-linux"), []byte("fake"), 0755); err != nil {
			t.Fatal(err)
		}

		libDir := filepath.Join(paths.Miva, version, "lib")
		if err := os.MkdirAll(libDir, 0755); err != nil {
			t.Fatal(err)
		}
		for _, hdr := range []string{"mvp_builtin.h", "mvp_copyable.h"} {
			if err := os.WriteFile(filepath.Join(libDir, hdr), []byte("header"), 0644); err != nil {
				t.Fatal(err)
			}
		}

		if err := SwitchVersion(paths, version); err != nil {
			t.Fatalf("SwitchVersion returned error: %v", err)
		}

		// old-file should no longer exist since Lib was removed and recreated
		if _, err := os.Stat(oldLibFile); !os.IsNotExist(err) {
			t.Error("old lib file should have been removed after switch")
		}
	})
}
