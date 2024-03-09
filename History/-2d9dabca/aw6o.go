package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fivePaisaWSCmd)
}

var fivePaisaWSCmd = &cobra.Command{
	Use:   "fivePaisaWS",
	Short: "Five Paisa Websocket",
	Long:  "Five Paisa Websocket",
	Run: func(cmd *cobra.Command, args []string) {
		fivepaisa.fivePaisaWS(fivePaisaTotp)
	},
}
