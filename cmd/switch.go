package cmd

import (
	"fmt"
	"miver/internal/miva"
	"miver/pkg/cli"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch <version>",
	Short: "Switch to a specific Miva version",
	Long: `Switch the active Miva version by linking the selected version's binaries and libraries.

The version must match one of the versions listed by 'miver list'.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "linux" {
			return
		}

		paths := miva.NewPaths(os.Getenv("HOME"))
		version := args[0]

		exists, err := miva.DirExists(paths.Miva)
		if err != nil {
			cli.Fatal(err)
		}
		if !exists {
			cli.Fatal(fmt.Errorf("Miva repository not found at %s. Please run 'miver init' first", paths.Miva))
		}

		versions, err := miva.AvailableVersions(paths.Miva)
		if err != nil {
			cli.Fatal(err)
		}

		found := false
		for _, v := range versions {
			if v[0] == version {
				found = true
				break
			}
		}
		if !found {
			cli.Fatal(fmt.Errorf("version %s not found in available versions", version))
		}

		cli.Info(fmt.Sprintf("Switching to Miva version %s", version))
		if err := miva.SwitchVersion(paths, version); err != nil {
			cli.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
