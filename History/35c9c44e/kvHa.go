package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sudeepbatra/alpha-hft/broker/cnbc"
)

func init() {
	rootCmd.AddCommand(cnbcQuotesCmd)
}

var cnbcQuotesCmd = &cobra.Command{
	Use:   "cnbcQuotes",
	Short: "Fetch and Save CNBC Quotes",
	Long:  "Fetch and Save CNBC Quotes",
	Run: func(cmd *cobra.Command, args []string) {
		cnbc.FetchQuoteDataAndSaveInDB()
	},
}
