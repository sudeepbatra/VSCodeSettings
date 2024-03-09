package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(smartApiCmd)
}

var smartApiCmd = &cobra.Command{
	Use:   "Smart API Login command",
	Short: "smart api login command for smartapi alpha-hft",
	Long:  `smart api login command for smartapi alpha-hft`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Info("Alpha HFT Beta version v0.1 -- HEAD")
	},
}
