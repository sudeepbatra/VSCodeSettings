package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/angelone_smartapi"
)

func init() {
	rootCmd.AddCommand(smartApiWSCmd)
}

var smartApiWSCmd = &cobra.Command{
	Use:   "angelOneSmartApiWSLogin",
	Short: "smart api login command for smartapi alpha-hft",
	Long:  `smart api login command for smartapi alpha-hft`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Invoking Angel One Smart API Session Token...")
		angelone_smartapi.AngelOneWebSocket()
		log.Info("Finished executing Angel One Smart API Session")
	},
}
