package cmd

import (
	"github.com/spf13/cobra"
	"github/masudur-rahman/AppsCodeServer/api"
)

var port string
var bypass bool
var stopTime int16

var startApp = &cobra.Command{
	Use:   "start",
	Short: "Start the app",
	Long:  "This starts the AppsCodeServer API",
	Run: func(cmd *cobra.Command, args []string) {
		api.AssignValues(port, bypass, stopTime)
		api.StartTheApp()
	},
}

func init() {
	startApp.PersistentFlags().StringVar(&port, "port", "8080", "port number for the server")
	startApp.PersistentFlags().BoolVarP(&bypass, "bypass", "b", true, "Bypass authentication parameter")
	startApp.PersistentFlags().Int16VarP(&stopTime, "stopTime", "s", 2, "The time after which the server will stop")

	rootCmd.AddCommand(startApp)
}