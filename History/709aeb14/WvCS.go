package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/logger"
)

func init() {
	smartAPILoginCmd.Flags().StringVarP(&smartApiTotp, "totp", "s", "", "specify the smart api totp from authenticator")
	rootCmd.AddCommand(smartAPILoginCmd)
}

var smartAPILoginCmd = &cobra.Command{
	Use:   "smartAPILogin",
	Short: "smart api client login command",
	Long:  "initializes the smart api client with access token and other creds required for further communication",
	Run: func(cmd *cobra.Command, args []string) {
		if smartApiTotp == "" {
			logger.Log.Error().Msg("no totp passed for smartAPILogin")
			return
		}
		smartapi.Login(smartApiTotp, "smartapi")
	},
}
