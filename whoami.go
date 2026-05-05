package main

import "github.com/spf13/cobra"

var whoamiVerbose bool

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display current DID and operator context",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: implement whoami logic
	},
}

func init() {
	whoamiCmd.Flags().BoolVar(&whoamiVerbose, "verbose", false, "show full configuration details")

	rootCmd.AddCommand(whoamiCmd)
}