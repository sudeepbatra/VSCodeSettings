package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/controller"
)

var smartApiTotp string
var fivePaisaTotp string
var createTable bool
var populateOldHistoricData bool
var useBrokersSavedState bool
var useProxy bool

func init() {
	initialzeBotCmd.Flags().BoolVarP(&createTable, "createTable", "c", false, "Specify if you want to create initial tables")
	initialzeBotCmd.Flags().StringVarP(&smartApiTotp, "smartApiTotp", "s", "", "Specify the Smart Api TOTP")
	initialzeBotCmd.Flags().BoolVarP(&populateOldHistoricData, "populateHistoricData", "o", false, "Specify true to populate old historica data")
	initialzeBotCmd.Flags().BoolVarP(&useBrokersSavedState, "useBrokersSavedState", "u", false, "Specify true to use the existing saved broker state while intializing")
	initialzeBotCmd.Flags().BoolVarP(&useProxy, "useProxy", "u", false, "Specify true to use the existing saved broker state while intializing")
	rootCmd.AddCommand(startBotCmd)
	rootCmd.AddCommand(initialzeBotCmd)
}

var rootCmd = &cobra.Command{
	Use:   "alpha-hft",
	Short: "A AlgoTrading Bot",
	Long:  "My Application is a proprietary Algo Trading bot",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Error("Error while trying to help method on command", err)
		}
	},
}

var startBotCmd = &cobra.Command{
	Use:   "startBot",
	Short: "Command to Start ALPHA-HFT",
	Long:  "Boots up the prices, strategies and most importantly start making money",
	Run: func(cmd *cobra.Command, args []string) {
		controller.StartAlphaHft()
	},
}

var initialzeBotCmd = &cobra.Command{
	Use:   "initializeBot",
	Short: "Command to Initialize ALPHA-HFT",
	Long:  "Triggers login, populates table for instrument and historical data",
	Run: func(cmd *cobra.Command, args []string) {
		controller.InitializeAlphaHftSOD(smartApiTotp, createTable, populateOldHistoricData, useBrokersSavedState)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("error while trying to run execute on rootcmd.", err)
		os.Exit(1)
	}
}
