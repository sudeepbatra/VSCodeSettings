package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/fivepaisa"
)

func init() {
	rootCmd.AddCommand(fivePaisaWSCmd)
}

var fivePaisaWSCmd = &cobra.Command{
	Use:   "fivePaisaWS",
	Short: "Login into Five Paisa and generate the accessToken",
	Long:  "Login command for Five Paisa and takes as an argument the totp",
	Run: func(cmd *cobra.Command, args []string) {
		if fivePaisaTotp == "" {
			log.Warn("WARNING: No totp passed")
		}

		fivepaisa.Login(fivePaisaTotp)
	},
}
