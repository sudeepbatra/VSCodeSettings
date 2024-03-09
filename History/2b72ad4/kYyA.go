package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/fivepaisa"
)

func init() {
	rootCmd.AddCommand(nseCorporateActionsCmd)
}

var nseCorporateActionsCmd = &cobra.Command{
	Use:   "nseCorporateActions",
	Short: "Five Paisa Websocket",
	Long:  "Five Paisa Websocket",
	Run: func(cmd *cobra.Command, args []string) {
		fivepaisa.WebsocketConnect()
	},
}
