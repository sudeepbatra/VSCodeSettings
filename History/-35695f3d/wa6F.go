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

func init() {
	breezeLoginCmd.Flags().StringVarP(&sessionToken, "sessionToken", "s", "", "Specify the sessionToken")
	rootCmd.AddCommand(breezeLoginCmd)
}
