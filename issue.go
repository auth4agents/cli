package main

import "github.com/spf13/cobra"

var issueScope string
var issueAud string
var issueTTL string
var issueDID string
var issueServer string

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Fetch a token for an agent",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: implement issue logic
	},
}

func init() {
	issueCmd.Flags().StringVar(&issueScope, "scope", "", "requested scope")
	issueCmd.Flags().StringVar(&issueAud, "aud", "", "target service audience")
	issueCmd.Flags().StringVar(&issueTTL, "ttl", "1h", "token lifetime")
	issueCmd.Flags().StringVar(&issueDID, "did", "", "agent DID")
	issueCmd.Flags().StringVar(&issueServer, "server", "", "Auth4Agent server URL")

	issueCmd.MarkFlagRequired("scope")
	issueCmd.MarkFlagRequired("aud")

	rootCmd.AddCommand(issueCmd)
}
