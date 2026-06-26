package cmd

import (
	"fmt"
	"miver/internal/miva"
	"miver/pkg/cli"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Miva versions",
	Long:  "Display all available Miva versions from the local repository in a formatted table.",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "linux" {
			return
		}

		paths := miva.NewPaths(os.Getenv("HOME"))
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

		// Calculate column widths
		widths := [3]int{len("Version"), len("Version Codename"), len("Release Time")}
		for _, v := range versions {
			if l := len(v[0]); l > widths[0] {
				widths[0] = l
			}
			if l := len(v[1]); l > widths[1] {
				widths[1] = l
			}
			if l := len(v[2]); l > widths[2] {
				widths[2] = l
			}
		}

		// ASCII box-drawing characters (no Unicode, broad terminal compatibility)
		const (
			horz = "-"
			vert = "|"
			plus = "+"
		)

		// Build separator lines
		col := func(n int) string {
			return strings.Repeat(horz, n+2)
		}
		topBorder := plus + col(widths[0]) + plus + col(widths[1]) + plus + col(widths[2]) + plus
		midBorder := plus + col(widths[0]) + plus + col(widths[1]) + plus + col(widths[2]) + plus
		botBorder := plus + col(widths[0]) + plus + col(widths[1]) + plus + col(widths[2]) + plus

		rowFmt := fmt.Sprintf("%s %%-%ds %s %%-%ds %s %%-%ds %s", vert, widths[0], vert, widths[1], vert, widths[2], vert)

		// Render table
		fmt.Println()
		fmt.Println(cli.Cyan(topBorder))
		fmt.Println(cli.Magenta(fmt.Sprintf(rowFmt, "Version", "Version Codename", "Release Time")))
		fmt.Println(cli.Cyan(midBorder))
		for _, v := range versions {
			fmt.Println(cli.Cyan(fmt.Sprintf(rowFmt, v[0], v[1], v[2])))
		}
		fmt.Println(cli.Cyan(botBorder))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
