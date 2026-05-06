package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var whoamiVerbose bool

var whoamiCmd = &cobra.Command{
    Use:   "whoami",
    Short: "Display current identity context",
    Run: func(cmd *cobra.Command, args []string) {
        // Try loading agent config first
        if cfg, err := LoadAgentConfig(); err == nil && cfg.DID != "" {
            fmt.Printf("Mode: Agent\n")
            fmt.Printf("DID: %s\n", cfg.DID)
            fmt.Printf("Operator: %s\n", cfg.OperatorID)
            fmt.Printf("Server: %s\n", cfg.ServerURL)
            return
        }
        
        // Try loading operator config
        if cfg, err := LoadOperatorConfig(); err == nil && cfg.ID != "" {
            fmt.Printf("Mode: Operator\n")
            fmt.Printf("ID: %s\n", cfg.ID)
            fmt.Printf("Domain: %s\n", cfg.Domain)
            fmt.Printf("Verified: %s\n", cfg.DomainVerifiedAt)
            fmt.Printf("Server: %s\n", cfg.ServerURL)
            return
        }
        
        fmt.Fprintf(os.Stderr, "Not configured. Run 'auth4agent init' first.\n")
        os.Exit(1)
    },
}

func init() {
	whoamiCmd.Flags().BoolVar(&whoamiVerbose, "verbose", false, "show full configuration details")
	rootCmd.AddCommand(whoamiCmd)
}