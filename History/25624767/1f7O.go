package cmd

import (
	"github.com/spf13/cobra"
)

var totp string

func init() {
	smartAPIWSCmd.Flags().StringVarP(&totp, "totp", "s", "", "Specify the totp")
	rootCmd.AddCommand(smartAPIWSCmd)
}

var smartAPIWSCmd = &cobra.Command{
	Use:   "smartAPIWSLogin",
	Short: "smart api websocket login command",
	Long:  `smart api websocket login command. It outputs the feedtoken, accesstoken and renewtoken needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if totp == "" {
			log.Error("Error: No totp passed")
			return
		}
		log.Info("Invoking Angel One Smart API Websocket Login...")
		log.Info("Finished executing Angel One Smart API Websocket Login!")
	},
}
