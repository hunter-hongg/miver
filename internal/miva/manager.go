// Package miva manages Miva version repositories.
package miva

import (
	"encoding/csv"
	"fmt"
	"miver/pkg/cli"
	"os"
	"os/exec"
	"path/filepath"
)

// Paths holds all relevant filesystem paths for a Miva installation.
type Paths struct {
	Home   string // ~/.miver
	Repo   string // ~/.miver/repo
	Miva   string // ~/.miver/repo/miva
	Bin    string // ~/.miver/bin
	Lib    string // ~/.miver/lib
	LibStd string // ~/.miver/lib/.std
}

// NewPaths constructs all Miva paths from the user's home directory.
func NewPaths(home string) Paths {
	miverHome := filepath.Join(home, ".miver")
	miverRepo := filepath.Join(miverHome, "repo")
	mivaRepo := filepath.Join(miverRepo, "miva")
	return Paths{
		Home:   miverHome,
		Repo:   miverRepo,
		Miva:   mivaRepo,
		Bin:    filepath.Join(miverHome, "bin"),
		Lib:    filepath.Join(miverHome, "lib"),
		LibStd: filepath.Join(miverHome, "lib", ".std"),
	}
}

// InitRepo clones the miva-publish repository if it doesn't exist.
func InitRepo(p Paths) error {
	cli.Step("Creating Miver directories")
	if err := os.MkdirAll(p.Repo, 0755); err != nil {
		return fmt.Errorf("create miver repo dir: %w", err)
	}
	cli.Success(fmt.Sprintf("Miver directories created at %s", p.Home))

	cli.Step("Cloning Miva repository")
	exists, err := DirExists(p.Miva)
	if err != nil {
		return err
	}
	if exists {
		cli.Info(fmt.Sprintf("Miva repository already exists at %s", p.Miva))
		return nil
	}

	cli.Info("Cloning from github.com/hunter-hongg/miva-publish...")
	cmd := exec.Command("git", "clone", "--depth=1",
		"https://github.com/hunter-hongg/miva-publish", p.Miva)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clone miva repo: %w", err)
	}
	cli.Success(fmt.Sprintf("Miva repository cloned to %s", p.Miva))
	cli.Success("Miver initialization completed.")
	return nil
}

// AvailableVersions reads and returns the version records from available.csv.
func AvailableVersions(mivaRepo string) ([][]string, error) {
	csvPath := filepath.Join(mivaRepo, "available.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("open available.csv: %w", err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read available.csv: %w", err)
	}
	return records, nil
}

// SwitchVersion switches the active Miva version to the given version tag.
func SwitchVersion(p Paths, version string) error {
	// Remove old lib
	if err := os.RemoveAll(p.Lib); err != nil {
		return fmt.Errorf("remove old lib: %w", err)
	}

	if err := os.MkdirAll(p.Bin, 0755); err != nil {
		return fmt.Errorf("create bin dir: %w", err)
	}
	if err := os.MkdirAll(p.Lib, 0755); err != nil {
		return fmt.Errorf("create lib dir: %w", err)
	}

	// Link miva binary
	binSrc := filepath.Join(p.Miva, version, "bin", "miva-linux")
	binDst := filepath.Join(p.Bin, "miva")
	if err := symlink(binSrc, binDst); err != nil {
		return fmt.Errorf("link miva binary: %w", err)
	}
	cli.Success(fmt.Sprintf("miva binary linked to %s", p.Bin))

	// Link frontend binary if present
	frontSrc := filepath.Join(p.Miva, version, "bin", "miva-frontend-linux")
	if _, err := os.Stat(frontSrc); err == nil {
		frontDst := filepath.Join(p.Bin, "miva-frontend")
		if err := symlink(frontSrc, frontDst); err != nil {
			return fmt.Errorf("link miva-frontend binary: %w", err)
		}
		cli.Success(fmt.Sprintf("miva frontend binary linked to %s", p.Bin))
	}

	// Link std library
	stdSrc := filepath.Join(p.Miva, version, "lib")
	stdDst := p.LibStd
	if err := symlink(stdSrc, stdDst); err != nil {
		return fmt.Errorf("link std lib: %w", err)
	}

	// Link headers into lib/
	for _, hdr := range []string{"mvp_builtin.h", "mvp_copyable.h"} {
		hdrSrc := filepath.Join(p.LibStd, hdr)
		hdrDst := filepath.Join(p.Lib, hdr)
		if err := symlink(hdrSrc, hdrDst); err != nil {
			return fmt.Errorf("link header %s: %w", hdr, err)
		}
	}

	// Link std ->
	stdLinkSrc := filepath.Join(p.Lib, ".std", "std")
	stdLinkDst := filepath.Join(p.Lib, "std")
	if err := symlink(stdLinkSrc, stdLinkDst); err != nil {
		return fmt.Errorf("link std symlink: %w", err)
	}

	cli.Success(fmt.Sprintf("miva std linked to %s", p.LibStd))
	return nil
}

// DirExists reports whether a directory exists.
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return true, nil
		}
		return false, fmt.Errorf("%s is not a directory", path)
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func symlink(src, dst string) error {
	return exec.Command("ln", "-sf", src, dst).Run()
}
