/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/abubakar-sadiq001/projscan/helpers"
	_ "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var ExtensionArgsSlice []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "projscan",
	Short: "A developer's pocket companion for quickly understanding what's inside any project directory.",
	Long:  `A developer's pocket companion for quickly understanding what's inside any project directory. ProjScan scans folders, breaks down file types by extension, counts lines of code, and flags React components that have grown a little too large. No configuration, no fuss, just actionable insight straight from the command line.`,

	Args: func(cmd *cobra.Command, args []string) error { // 1 argument is required to run this command.
		// if len(args) < 1 {
		// 	return fmt.Errorf("requires at least one argument")
		// }

		// if len(args) > 1 {
		// 	return fmt.Errorf("too many arguments (max 1)")
		// }

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		fs := helpers.FileStats{}

		version, _ := cmd.Flags().GetBool("version")
		if version {
			helpers.PrintVersion()
			return
		}

		if len(args) >= 1 {

			info, err := os.Stat(args[0])
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("Error: Directory '%s' does not exist\n", args[0])
					return
				}
				fmt.Printf("Error: %v\n", err)
				return
			}
			if !info.IsDir() {
				fmt.Printf("Error: '%s' is not a directory\n", args[0])
				return
			}

			time.Sleep(1 * time.Second)
			if args[0] == "." {
				fmt.Println("Opening current dir for scanning")
			} else {

				fmt.Printf("Opening dir %s for scanning\n", args[0])
			}

			fmt.Printf("\n\n")
			now := time.Now()

			stopSpinner := make(chan struct{})
			spinnerDone := make(chan struct{})
			go func() {
				fs.ShowSpinner(stopSpinner)
				close(spinnerDone)
			}()

			wg.Add(1)
			go helpers.ScanDirectory(args[0], ExtensionArgsSlice, &wg)
			wg.Wait()
			close(stopSpinner)
			<-spinnerDone

			fs.PrintFilesWithLines()
			fs.PrintTotalFilesStat()

			elapsed := time.Since(now)
			fmt.Printf("\n\nScan completed in %s\n", elapsed)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringSliceVarP(&ExtensionArgsSlice, "extension", "e", ExtensionArgsSlice, "Specify any extension type to scan")
	rootCmd.Flags().BoolP("version", "v", false, "Print the version information")
}
