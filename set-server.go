package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var setServerCmd = &cobra.Command{
	Use:   "set-server",
	Short: "Update server URL in configuration",
	
	Run: func(cmd *cobra.Command, args []string) {
		serverURL, _ := cmd.Flags().GetString("url")
		if serverURL == "" {
			fmt.Fprintf(os.Stderr, "server URL is required\n")
			os.Exit(1)
		}
		
		err := UpdateServerURL(serverURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to update config: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Server URL updated to: %s\n", serverURL)
	},
}

func init() {
	setServerCmd.Flags().String("url", "", "Auth4Agent server URL")
	setServerCmd.MarkFlagRequired("url")
	rootCmd.AddCommand(setServerCmd)
}