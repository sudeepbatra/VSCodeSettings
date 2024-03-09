package cmd

import (
	"github.com/spf13/cobra"
	smartapiwslogin "github.com/sudeepbatra/alpha-hft/broker/smartapi"
)

func init() {
	rootCmd.AddCommand(smartAPIWSCmd)
}

var smartAPIMarketDataWSCmd = &cobra.Command{
	Use:   "smartAPIMarketDataWS",
	Short: "smart api market data websocket",
	Long:  `smart api websocket login command. It outputs the feedtoken, accesstoken and renewtoken needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if totp == "" {
			log.Error("Error: No totp passed")
			return
		}
		log.Info("Invoking Angel One Smart API Websocket Login...")
		smartapiwslogin.SmartAPIWebSocketLogin(totp)
		log.Info("Finished executing Angel One Smart API Websocket Login!")
	},
}
