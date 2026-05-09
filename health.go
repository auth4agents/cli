package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var healthJSON bool

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check server health status",
	
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := LoadAgentConfig()
		if err != nil {
			// Try loading operator config
			opCfg, err := LoadOperatorConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "no config found\n")
				os.Exit(1)
			}
			cfg = &AgentConfig{ServerURL: opCfg.ServerURL}
		}
		
		client := NewClient(cfg.ServerURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		var healthResp map[string]interface{}
		
		// Try new health endpoint first
		err = client.Get(ctx, "/health", &healthResp)
		if err != nil {
			// Fallback to old method
			fmt.Fprintf(os.Stderr, "Server unreachable: %v\n", err)
			os.Exit(1)
		}
		
		if healthJSON {
			json.NewEncoder(os.Stdout).Encode(healthResp)
			return
		}
		
		fmt.Printf("Server status: %s\n", healthResp["status"])
		if ts, ok := healthResp["timestamp"]; ok {
			fmt.Printf("Timestamp: %v\n", ts)
		}
	},
}

func init() {
	healthCmd.Flags().BoolVar(&healthJSON, "json", false, "output JSON")
	rootCmd.AddCommand(healthCmd)
}