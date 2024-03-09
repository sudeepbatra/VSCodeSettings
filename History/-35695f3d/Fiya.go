package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/breeze"
	"github.com/sudeepbatra/alpha-hft/logger"
)

var log = logger.GetLogger()

var (
	sessionToken string
)

func init() {
	breezeLoginCmd.Flags().StringVarP(&sessionToken, "sessionToken", "s", "", "Specify the sessionToken")
	rootCmd.Flags().BoolP("useProxy", "p", false, "Set to true to use a proxy")
	rootCmd.AddCommand(breezeLoginCmd)
}

var breezeLoginCmd = &cobra.Command{
	Use:   "loginBreeze",
	Short: "Login into Breeze and Stores the Token",
	Long:  "Login command for Breeze and takes as an argument the sessionToken, useProxy(optional)",
	Run: func(cmd *cobra.Command, args []string) {
		if sessionToken == "" {
			log.Warn("WARNING: No sessionToken passed")
		}
		useProxy, _ := cmd.Flags().GetBool("useProxy")
		log.Debug("Use Proxy flag:", useProxy)

		breeze.Login(sessionToken, useProxy)
	},
}
