package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/sudeepbatra/alpha-hft/broker/breeze"
)

var (
	sessionToken string
)

var rootCmd = &cobra.Command{
	Use:   "alpha-hft",
	Short: "A AlgoTrading Bot",
	Long:  "My Application is a propreitary Algo Trading bot",
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no subcommand is provided
		cmd.Help()
	},
}

var breezeLoginCmd = &cobra.Command{
	Use:   "loginBreeze",
	Short: "Login into Breeze and Stores the Token",
	Run: func(cmd *cobra.Command, args []string) {
		if sessionToken == "" {
			log.Println("WARNING: No sessionToken passed")
		}
		breeze.Login(sessionToken)
	},
}

func init() {
	breezeLoginCmd.Flags().StringVarP(&sessionToken, "sessionToken", "s", "", "Specify the sessionToken")
	rootCmd.AddCommand(breezeLoginCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
