package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/nse"
)

func init() {
	rootCmd.AddCommand(nseCorporateActionsCmd)
}

var nseCorporateActionsCmd = &cobra.Command{
	Use:   "cnbcQuotes",
	Short: "NSE Corporate Actions",
	Long:  "NSE Corporate Actions",
	Run: func(cmd *cobra.Command, args []string) {
		nse.FetchAndSaveNSECorporateActions()
	},
}
