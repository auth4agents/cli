package main

import "github.com/spf13/cobra"

var revokeDID string
var revokeServer string

var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke an agent's credentials",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: implement revoke logic
	},
}

func init() {
	revokeCmd.Flags().StringVar(&revokeDID, "did", "", "agent DID to revoke")
	revokeCmd.Flags().StringVar(&revokeServer, "server", "", "Auth4Agent server URL")

	revokeCmd.MarkFlagRequired("did")

	rootCmd.AddCommand(revokeCmd)
}
