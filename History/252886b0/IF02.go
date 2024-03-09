package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/sudeepbatra/alpha-hft/broker/breeze"
	"github.com/sudeepbatra/alpha-hft/cmd"
)

var (
	sessionToken string
)

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

func main() {
	cmd.Execute()
}
