package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "AppsCodeServer",
	Short: "It's a server containing workers of appscode",
	Long: "All the worker profile of the AppsCode Ltd." +
		" is included in this server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome.........!!!")
	},
}

func Execute() {
	//rootCmd.AddCommand(versionCmd)
	//Init()
	if err := rootCmd.Execute(); err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
}
