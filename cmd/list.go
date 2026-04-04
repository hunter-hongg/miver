/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miver/internal/csvr"
	"miver/internal/system"
	"miver/pkg/color"
	_ "miver/pkg/color"
	"miver/pkg/out"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "linux" {
			home := os.Getenv("HOME")
			miver_home := fmt.Sprintf("%s/.miver", home)
			miver_repo := fmt.Sprintf("%s/repo", miver_home)
			miva_repo := fmt.Sprintf("%s/miva", miver_repo)
			if res, err := system.DirExists(miva_repo); err == nil && res {
				lists := csvr.ReadAvailable(miva_repo)

				out.Info("Listing available versions...\n")

				var (
					versionMax int
					animalMax  int
					timeMax    int
				)
				dealMax := func(max int, cur string) int {
					l := len(cur)
					if l > max {
						return l
					} else {
						return max
					}
				}
				toMax := func(max int, cur string) string {
					return fmt.Sprintf("%-*s", max, cur)
				}
				cpStr := func(times int, s string) string {
					ss := ""
					for i := 0; i < times; i++ {
						ss += s
					}
					return ss
				}
				for _, list := range lists {
					version := list[0]
					animal := list[1]
					time := list[2]
					versionMax = dealMax(versionMax, version)
					animalMax = dealMax(animalMax, animal)
					timeMax = dealMax(timeMax, time)
				}
				versionMax = dealMax(versionMax, "Version")
				animalMax = dealMax(animalMax, "Version Codename")
				timeMax = dealMax(timeMax, "Release Time")
				outFg := fmt.Sprintf(
					"+-%s-+-%s-+-%s-+",
					cpStr(versionMax, "-"),
					cpStr(animalMax, "-"),
					cpStr(timeMax, "-"),
				)
				outFgN := fmt.Sprintf(
					"%s\n",
					outFg,
				)
				fmt.Print(color.Cyan(
					outFgN +
						fmt.Sprintf(
							"| %s | %s | %s |\n",
							toMax(versionMax, "Version"),
							toMax(animalMax, "Version Codename"),
							toMax(timeMax, "Release Time"),
						) +
						outFgN,
				))
				for _, list := range lists {
					version := list[0]
					animal := list[1]
					time := list[2]
					fmt.Println(color.Cyan(
						fmt.Sprintf(
							"| %s | %s | %s |\n",
							toMax(versionMax, version),
							toMax(animalMax, animal),
							toMax(timeMax, time),
						) + outFg,
					))
				}
			} else {
				out.Error(fmt.Sprintf(
					"Miva repository does not exist at %s. Please run 'miver init' first.",
					miva_repo,
				))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
