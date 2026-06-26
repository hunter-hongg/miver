package cmd

import (
	"fmt"
	"miver/internal/miva"
	"miver/pkg/cli"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Miver environment",
	Long: `Set up the Miver directory structure and clone the Miva version repository.

This command:
1. Creates ~/.miver/repo/ for storing version data
2. Clones the miva-publish repository with available version definitions`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Step("Checking system...")
		cli.Info(fmt.Sprintf("Current system is %s", runtime.GOOS))

		if err := checkPlatform(runtime.GOOS); err != nil {
			cli.Fatal(err)
		}

		if runtime.GOOS != "linux" {
			return
		}

		paths := miva.NewPaths(os.Getenv("HOME"))
		if err := miva.InitRepo(paths); err != nil {
			cli.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func checkPlatform(goos string) error {
	switch goos {
	case "windows":
		return fmt.Errorf("Windows is not supported. Please use WSL instead")
	case "darwin":
		return fmt.Errorf("macOS is not supported. Please use Linux instead")
	case "linux":
		return nil
	default:
		return fmt.Errorf("unsupported system: %s", goos)
	}
}
