package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/breeze"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var websocketCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of alpha-hft",
	Long:  `All software has versions. This is alpha-hft's`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Running the Websocket Command")
		breeze.Login(sessionToken, useProxy)
	},
}
