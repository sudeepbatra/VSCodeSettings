package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/angelone_smartapi"
)

func init() {
	rootCmd.AddCommand(smartApiCmd)
}

var smartApiCmd = &cobra.Command{
	Use:   "angelOneSmartApiLogin",
	Short: "smart api login command for smartapi alpha-hft",
	Long:  `smart api login command for smartapi alpha-hft`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Invoking Angel One Smart API Session Token...")
		angelone_smartapi.SmartApiSessionToken()
		log.Info("Finished executing Angel One Smart API Session")
	},
}
