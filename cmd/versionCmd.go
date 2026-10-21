/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/abubakar-sadiq001/projscan/helpers"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of the currently installed version of the application",
	Long:  `Go version returns the currently installed version of the application. This command is useful for verifying the version of the application that is currently installed on your system. It can also be used to check for updates or to troubleshoot issues related to version compatibility.`,
	Run: func(cmd *cobra.Command, args []string) {
		helpers.PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
