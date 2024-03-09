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
	Long:  `smart api market data websocket.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Invoking Angel One Smart API Websocket Login...")
		smartapiwslogin.SmartAPIWebSocketLogin(totp)
		log.Info("Finished executing Angel One Smart API Websocket Login!")
	},
}
