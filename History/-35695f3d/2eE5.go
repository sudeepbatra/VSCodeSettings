package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/breeze"
)

var (
	sessionToken string
)

func init() {
	breezeLoginCmd.Flags().StringVarP(&sessionToken, "sessionToken", "s", "", "Specify the sessionToken")
	rootCmd.AddCommand(breezeLoginCmd)
}

var breezeLoginCmd = &cobra.Command{
	Use:   "loginBreeze",
	Short: "Login into Breeze and Stores the Token",
	Long:  "Login command for Breeze and takes as an argument the sessionToken, useProxy(optional)",
	Run: func(cmd *cobra.Command, args []string) {
		if sessionToken == "" {
			log.Println("WARNING: No sessionToken passed")
		}
		useProxy, _ := cmd.Flags().GetBool("useProxy")

		breeze.Login(sessionToken)
	},
}
