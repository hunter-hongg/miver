/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miver/internal/csvr"
	"miver/internal/system"
	"miver/pkg/errs"
	"miver/pkg/out"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
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
			if res, err := system.DirExists(miva_repo); err != nil || !res {
				errs.DealError(err)
				errs.DealError(fmt.Errorf("miver directory not found"))
			}
			if len(args) < 1 {
				errs.DealError(fmt.Errorf("switch version required"))
			} else if len(args) > 1 {
				out.Warning("there are too many args")
			}
			version := args[0]
			out.Info(fmt.Sprintf("switch version %s", version))

			lists := csvr.ReadAvailable(miva_repo)
			hasVer := false

			for _, list := range lists {
				if list[0] == version {
					hasVer = true
					break
				}
			}

			if !hasVer {
				errs.DealError(fmt.Errorf("version %s not found", version))
			}

			miver_bin := fmt.Sprintf("%s/bin", miver_home)
			miver_lib := fmt.Sprintf("%s/lib", miver_home)
			miver_lib_std := fmt.Sprintf("%s/lib/.std", miver_home)

			os.MkdirAll(miver_bin, 0755)
			os.MkdirAll(miver_lib, 0755)

			lnToBin := exec.Command(
				"ln", "-sf", 
				fmt.Sprintf("%s/%s/bin/miva-linux", miva_repo, version),
				fmt.Sprintf("%s/miva", miver_bin),
			)
			err := lnToBin.Run()
			errs.DealError(err)
			out.Info(fmt.Sprintf("miva binary linked to %s", miver_bin))

			lnToStd := exec.Command(
				"ln", "-sf", 
				fmt.Sprintf("%s/%s/lib", miva_repo, version),
				miver_lib_std,
			)
			err = lnToStd.Run()
			errs.DealError(err)

			lnHeader := exec.Command(
				"ln", "-sf", 
				fmt.Sprintf("%s/mvp_builtin.h", miver_lib_std),
				fmt.Sprintf("%s/mvp_builtin.h", miver_lib),
			)
			err = lnHeader.Run()
			errs.DealError(err)

			lnHeader = exec.Command(
				"ln", "-sf", 
				fmt.Sprintf("%s/mvp_copyable.h", miver_lib_std),
				fmt.Sprintf("%s/mvp_copyable.h", miver_lib),
			)
			err = lnHeader.Run()
			errs.DealError(err)

			lnStd := exec.Command(
				"ln", "-sf", 
				fmt.Sprintf("%s/.std/std", miver_lib),
				fmt.Sprintf("%s/std", miver_lib),
			)
			err = lnStd.Run()
			errs.DealError(err)

			out.Info(fmt.Sprintf("miva std linked to %s", miver_lib_std))
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
