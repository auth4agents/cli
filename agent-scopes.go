package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var scopesCmd = &cobra.Command{
	Use:   "agent-scopes",
	Short: "Manage agent scopes",
}

var scopesSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set allowed scopes for an agent",

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := LoadAgentConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load agent config: %v\n", err)
			os.Exit(1)
		}

		scopes, _ := cmd.Flags().GetStringSlice("scopes")
		
		client := NewClient(cfg.ServerURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Use operator ID from config
		path := fmt.Sprintf("/v1/operators/%s/agents/scopes", cfg.OperatorID)
		
		payload := map[string]interface{}{
			"did":            cfg.DID,
			"allowed_scopes": scopes,
		}
		
		var out map[string]interface{}
		
		err = client.Patch(ctx, path, payload, &out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
			os.Exit(1)
		}
		
		data, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(data))
	},
}

var scopesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List allowed scopes for an agent",
	
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := LoadAgentConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load agent config: %v\n", err)
			os.Exit(1)
		}
		
		client := NewClient(cfg.ServerURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		// Get agent details including scopes
		path := fmt.Sprintf("/v1/agent/%s", cfg.DID)
		
		var agent struct {
			AllowedScopes []string `json:"allowed_scopes"`
		}
		
		err = client.Get(ctx, path, &agent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Allowed scopes for %s:\n", cfg.DID)
		for _, scope := range agent.AllowedScopes {
			fmt.Printf("  - %s\n", scope)
		}
		
		if len(agent.AllowedScopes) == 0 {
			fmt.Println("  (none)")
		}
	},
}

func init() {
	// Only add flags for the set command
	scopesSetCmd.Flags().StringSlice("scopes", []string{}, "allowed scopes (comma-separated)")
	scopesSetCmd.MarkFlagRequired("scopes")
	
	// The list command doesn't need any flags
	// Add both commands to scopesCmd
	scopesCmd.AddCommand(scopesSetCmd, scopesListCmd)
	
	// Add scopesCmd to rootCmd
	rootCmd.AddCommand(scopesCmd)
}