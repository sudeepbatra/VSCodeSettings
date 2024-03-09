package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of alpha-hft",
	Long:  `All software has versions. This is alpha-hft's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Alpha HFT Beta version v0.1 -- HEAD")
	},
}
