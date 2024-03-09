package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/cnbc"
	"github.com/sudeepbatra/alpha-hft/broker/nse"
)

func init() {
	rootCmd.AddCommand(nseCorporateActionsCmd)
}

var nseCorporateActionsCmd = &cobra.Command{
	Use:   "cnbcQuotes",
	Short: "Fetch and Save CNBC Quotes",
	Long:  "Fetch and Save CNBC Quotes",
	Run: func(cmd *cobra.Command, args []string) {
		cnbc.FetchQuoteDataAndSaveInDB()
		nse.FetchAndSaveNSECorporateActions()
	},
}
