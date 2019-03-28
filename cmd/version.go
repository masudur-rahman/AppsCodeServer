package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of AppsCodeServer",
	Long:  "The version of the AppsCodeServer app is",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("AppsCodeServer - v2.0")
	},
}
