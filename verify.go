package main

import "github.com/spf13/cobra"

var verifyToken string
var verifyServer string

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Decode and verify a token",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: implement verify logic
	},
}

func init() {
	verifyCmd.Flags().StringVar(&verifyToken, "token", "", "JWT token string")
	verifyCmd.Flags().StringVar(&verifyServer, "server", "", "Auth4Agent server URL")

	verifyCmd.MarkFlagRequired("token")

	rootCmd.AddCommand(verifyCmd)
}
