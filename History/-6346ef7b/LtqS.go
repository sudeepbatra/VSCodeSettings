package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/fivepaisa"
)

var fivePaisaTotp string

func init() {
	fivePaisaLoginCmd.Flags().StringVarP(&fivePaisaTotp, "fivePaisaTotp", "t", "", "Specify the Five Paisa TOTP")
	rootCmd.AddCommand(fivePaisaLoginCmd)
}

var fivePaisaLoginCmd = &cobra.Command{
	Use:   "fivePaisaLogin",
	Short: "Login into Five Paisa and generate the accessToken",
	Long:  "Login command for Five Paisa and takes as an argument the totp",
	Run: func(cmd *cobra.Command, args []string) {
		if fivePaisaTotp == "" {
			log.Warn("WARNING: No totp passed")
		}

		fivepaisa.FivePaisaLogin(fivePaisaTotp)
	},
}