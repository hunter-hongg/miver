package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "miver",
	Short: "Miva version manager",
	Long:  `Miver is a CLI tool to manage and switch between different versions of the Miva programming language compiler.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
