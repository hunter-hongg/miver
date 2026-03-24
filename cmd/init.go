/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miver/internal/system"
	"miver/pkg/errs"
	"miver/pkg/out"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		curos := runtime.GOOS
		out.Step("Checking system...")
		out.Info(fmt.Sprintf("Current system is %s", curos))
		err := system.CheckSys(curos)
		errs.DealError(err)

		if curos == "linux" {
			home := os.Getenv("HOME")
			miver_home := fmt.Sprintf("%s/.miver", home)
			miver_repo := fmt.Sprintf("%s/repo", miver_home)
			miva_repo := fmt.Sprintf("%s/miva", miver_repo)

			out.Step("Creating Miver directories...")
			err = os.MkdirAll(miver_repo, 0755)
			errs.DealError(err)
			out.Info(fmt.Sprintf("Miver directories created at %s", miver_home))

			out.Step("Cloning Miva repository...")
			exists, err := system.DirExists(miva_repo)
			errs.DealError(err)
			if exists {
				out.Info(fmt.Sprintf("Miva repository already exists at %s", miva_repo))
			} else {
				cmd := exec.Command(
					"git", "clone", "--depth=1",
					"https://github.com/hunter-hongg/miva-publish",
					miva_repo,
				)
				err := cmd.Run()
				errs.DealError(err)
				out.Info(fmt.Sprintf("Miva repository cloned to %s", miva_repo))
				out.Info("Miver initialization completed.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
