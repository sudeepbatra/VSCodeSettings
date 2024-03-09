package cmd

import (
	"github.com/spf13/cobra"
)

var totp string

func init() {
	smartApiWSCmd.Flags().StringVarP(&totp, "totp", "s", "", "Specify the totp")
	rootCmd.AddCommand(smartApiWSCmd)
}

var smartApiWSCmd = &cobra.Command{
	Use:   "smartAPIWSLogin",
	Short: "smart api websocket login command",
	Long:  `smart api websocket login command. It outputs the feedtoken, accesstoken and renewtoken needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Invoking Angel One Smart API Websocket Login...")
		log.Info("Finished executing Angel One Smart API Websocket Login!")
	},
}
