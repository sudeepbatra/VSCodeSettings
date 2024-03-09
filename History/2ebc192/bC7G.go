package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "alpha-hft",
	Short: "A AlgoTrading Bot",
	Long:  "My Application is a proprietary Algo Trading bot",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Error("Error while trying to ")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("error while trying to run execute on rootcmd.")
		fmt.Println(err)
		os.Exit(1)
	}
}
