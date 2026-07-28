/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/abubakar-sadiq001/projscan/helpers"
	_ "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var ExtensionArgsSlice []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "my-cli [argument]",
	Short: "A powerful CLI tool built with Cobra",
	Long: `projscan is lightweight CLI tool for scanning My CLI is a demonstration application built with Cobra.
	 
	  This application shows how to create professional command-line
   	  tools with proper flag handling, subcommands, and configuration.`,

	Args: func(cmd *cobra.Command, args []string) error { // 1 argument is required to run this command.
		if len(args) < 1 {
			return fmt.Errorf("requires at least one argument")
		}

		if len(args) > 1 {
			return fmt.Errorf("too many arguments (max 1)")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
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
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Opening your current working directory: %s\n", wd)
		} else {

			fmt.Printf("Opening directory %s for scanning\n", args[0])
		}

		fmt.Printf("\n\n")
		helpers.ScanDirectory(args[0], ExtensionArgsSlice)
		helpers.PrintStats()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringSliceVarP(&ExtensionArgsSlice, "extension", "e", ExtensionArgsSlice, "Specify any extension type to scan")
}
