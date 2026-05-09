package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var revokeDID string
var revokeServer string
var revokeReason string
var revokeJSON bool
var revokeToken string
var revokeJTI string

var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke agent credentials or tokens",
}

// Revoke a specific token
var revokeTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Revoke a specific JWT token",
	
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := LoadAgentConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
			os.Exit(1)
		}
		
		serverURL := cfg.ServerURL
		if revokeServer != "" {
			serverURL = revokeServer
		}
		
		client := NewClient(serverURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		
		payload := map[string]interface{}{
			"token":  revokeToken,
			"reason": revokeReason,
		}
		
		if revokeJTI != "" {
			payload["jti"] = revokeJTI
		}
		
		var resp map[string]interface{}
		
		err = client.Post(ctx, "/v1/revoke/token", payload, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "revocation failed: %v\n", err)
			os.Exit(1)
		}
		
		if revokeJSON {
			json.NewEncoder(os.Stdout).Encode(resp)
			return
		}
		
		fmt.Printf("Token revoked successfully\n")
		if jti, ok := resp["jti"]; ok {
			fmt.Printf("JTI: %v\n", jti)
		}
		if revokedAt, ok := resp["revoked_at"]; ok {
			fmt.Printf("Revoked at: %v\n", revokedAt)
		}
	},
}

// Revoke an entire agent (permanent)
var revokeAgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Permanently revoke an agent (requires operator authentication)",
	
	Run: func(cmd *cobra.Command, args []string) {
		operatorCfg, err := LoadOperatorConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "operator config not found: %v\n", err)
			os.Exit(1)
		}
		
		if operatorCfg.ID == "" {
			fmt.Fprintf(os.Stderr, "operator not registered. Run 'auth4agent register operator' first.\n")
			os.Exit(1)
		}
		
		agentCfg, err := LoadAgentConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "agent config not found: %v\n", err)
			os.Exit(1)
		}
		
		agentID := revokeDID
		if agentID == "" {
			// Try to get from config
			if agentCfg.DID != "" {
				// Extract agent ID from DID (format: did:agent:domain:hash)
				// For now, we need the agent's database ID, not DID
				fmt.Fprintf(os.Stderr, "agent ID (not DID) is required. Use --id flag or check config\n")
				os.Exit(1)
			}
		}
		
		serverURL := operatorCfg.ServerURL
		if revokeServer != "" {
			serverURL = revokeServer
		}
		
		client := NewClient(serverURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		
		path := fmt.Sprintf("/v1/operators/%s/agents/%s/revoke", operatorCfg.ID, agentID)
		
		payload := map[string]interface{}{
			"reason": revokeReason,
		}
		
		var resp map[string]interface{}
		
		err = client.Post(ctx, path, payload, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "revocation failed: %v\n", err)
			os.Exit(1)
		}
		
		if revokeJSON {
			json.NewEncoder(os.Stdout).Encode(resp)
			return
		}
		
		fmt.Printf("Agent revoked permanently\n")
		if agentDid, ok := resp["agent_did"]; ok {
			fmt.Printf("DID: %v\n", agentDid)
		}
		if revokedAt, ok := resp["revoked_at"]; ok {
			fmt.Printf("Revoked at: %v\n", revokedAt)
		}
	},
}

// Suspend an agent (temporary)
var suspendAgentCmd = &cobra.Command{
	Use:   "suspend",
	Short: "Temporarily suspend an agent (requires operator authentication)",
	
	Run: func(cmd *cobra.Command, args []string) {
		operatorCfg, err := LoadOperatorConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "operator config not found: %v\n", err)
			os.Exit(1)
		}
		
		if operatorCfg.ID == "" {
			fmt.Fprintf(os.Stderr, "operator not registered. Run 'auth4agent register operator' first.\n")
			os.Exit(1)
		}
		
		agentID, _ := cmd.Flags().GetString("agent-id")
		if agentID == "" {
			fmt.Fprintf(os.Stderr, "--agent-id is required\n")
			os.Exit(1)
		}
		
		serverURL := operatorCfg.ServerURL
		if revokeServer != "" {
			serverURL = revokeServer
		}
		
		client := NewClient(serverURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		
		path := fmt.Sprintf("/v1/operators/%s/agents/%s/suspend", operatorCfg.ID, agentID)
		
		payload := map[string]interface{}{
			"reason": revokeReason,
		}
		
		var resp map[string]interface{}
		
		err = client.Post(ctx, path, payload, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "suspension failed: %v\n", err)
			os.Exit(1)
		}
		
		if revokeJSON {
			json.NewEncoder(os.Stdout).Encode(resp)
			return
		}
		
		fmt.Printf("Agent suspended temporarily\n")
		if suspendedAt, ok := resp["suspended_at"]; ok {
			fmt.Printf("Suspended at: %v\n", suspendedAt)
		}
	},
}

// Reactivate a suspended agent
var reactivateAgentCmd = &cobra.Command{
	Use:   "reactivate",
	Short: "Reactivate a suspended agent (requires operator authentication)",
	
	Run: func(cmd *cobra.Command, args []string) {
		operatorCfg, err := LoadOperatorConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "operator config not found: %v\n", err)
			os.Exit(1)
		}
		
		if operatorCfg.ID == "" {
			fmt.Fprintf(os.Stderr, "operator not registered. Run 'auth4agent register operator' first.\n")
			os.Exit(1)
		}
		
		agentID, _ := cmd.Flags().GetString("agent-id")
		if agentID == "" {
			fmt.Fprintf(os.Stderr, "--agent-id is required\n")
			os.Exit(1)
		}
		
		serverURL := operatorCfg.ServerURL
		if revokeServer != "" {
			serverURL = revokeServer
		}
		
		client := NewClient(serverURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		
		path := fmt.Sprintf("/v1/operators/%s/agents/%s/reactivate", operatorCfg.ID, agentID)
		
		var resp map[string]interface{}
		
		err = client.Post(ctx, path, nil, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reactivation failed: %v\n", err)
			os.Exit(1)
		}
		
		if revokeJSON {
			json.NewEncoder(os.Stdout).Encode(resp)
			return
		}
		
		fmt.Printf("Agent reactivated successfully\n")
		if reactivatedAt, ok := resp["reactivated_at"]; ok {
			fmt.Printf("Reactivated at: %v\n", reactivatedAt)
		}
	},
}

// List revoked/suspended agents
var listRevokedCmd = &cobra.Command{
	Use:   "list",
	Short: "List revoked or suspended agents",
	
	Run: func(cmd *cobra.Command, args []string) {
		operatorCfg, err := LoadOperatorConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "operator config not found: %v\n", err)
			os.Exit(1)
		}
		
		if operatorCfg.ID == "" {
			fmt.Fprintf(os.Stderr, "operator not registered\n")
			os.Exit(1)
		}
		
		serverURL := operatorCfg.ServerURL
		if revokeServer != "" {
			serverURL = revokeServer
		}
		
		client := NewClient(serverURL)
		
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		
		// Get all agents for this operator
		path := fmt.Sprintf("/v1/operators/%s/agents", operatorCfg.ID)
		
		var agents []struct {
			ID         string     `json:"id"`
			DID        string     `json:"did"`
			Status     string     `json:"status"`
			RevokedAt  *time.Time `json:"revoked_at"`
			SuspendedAt *time.Time `json:"suspended_at"`
		}
		
		err = client.Get(ctx, path, &agents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list agents: %v\n", err)
			os.Exit(1)
		}
		
		revokedOrSuspended := []interface{}{}
		for _, agent := range agents {
			if agent.Status == "revoked" || agent.Status == "suspended" {
				revokedOrSuspended = append(revokedOrSuspended, agent)
			}
		}
		
		if revokeJSON {
			json.NewEncoder(os.Stdout).Encode(revokedOrSuspended)
			return
		}
		
		if len(revokedOrSuspended) == 0 {
			fmt.Printf("No revoked or suspended agents found\n")
			return
		}
		
		fmt.Printf("Revoked/Suspended Agents:\n\n")
		for _, agent := range revokedOrSuspended {
			a := agent.(struct {
				ID         string     `json:"id"`
				DID        string     `json:"did"`
				Status     string     `json:"status"`
				RevokedAt  *time.Time `json:"revoked_at"`
				SuspendedAt *time.Time `json:"suspended_at"`
			})
			fmt.Printf("ID: %s\n", a.ID)
			fmt.Printf("DID: %s\n", a.DID)
			fmt.Printf("Status: %s\n", a.Status)
			if a.RevokedAt != nil {
				fmt.Printf("Revoked at: %s\n", a.RevokedAt)
			}
			if a.SuspendedAt != nil {
				fmt.Printf("Suspended at: %s\n", a.SuspendedAt)
			}
			fmt.Println()
		}
	},
}

func init() {
	// Token revocation flags
	revokeTokenCmd.Flags().StringVar(&revokeToken, "token", "", "JWT token to revoke")
	revokeTokenCmd.Flags().StringVar(&revokeJTI, "jti", "", "JWT ID to revoke")
	revokeTokenCmd.Flags().StringVar(&revokeReason, "reason", "revoked by user", "revocation reason")
	revokeTokenCmd.Flags().BoolVar(&revokeJSON, "json", false, "output JSON")
	revokeTokenCmd.Flags().StringVar(&revokeServer, "server", "", "Auth4Agent server URL")
	revokeTokenCmd.MarkFlagsOneRequired("token", "jti")
	
	// Agent revocation flags
	revokeAgentCmd.Flags().StringVar(&revokeDID, "id", "", "agent ID to revoke")
	revokeAgentCmd.Flags().StringVar(&revokeReason, "reason", "revoked by operator", "revocation reason")
	revokeAgentCmd.Flags().BoolVar(&revokeJSON, "json", false, "output JSON")
	revokeAgentCmd.Flags().StringVar(&revokeServer, "server", "", "Auth4Agent server URL")
	revokeAgentCmd.MarkFlagRequired("id")
	
	// Suspend flags
	suspendAgentCmd.Flags().String("agent-id", "", "agent ID to suspend")
	suspendAgentCmd.Flags().StringVar(&revokeReason, "reason", "suspended by operator", "suspension reason")
	suspendAgentCmd.Flags().BoolVar(&revokeJSON, "json", false, "output JSON")
	suspendAgentCmd.Flags().StringVar(&revokeServer, "server", "", "Auth4Agent server URL")
	suspendAgentCmd.MarkFlagRequired("agent-id")
	
	// Reactivate flags
	reactivateAgentCmd.Flags().String("agent-id", "", "agent ID to reactivate")
	reactivateAgentCmd.Flags().BoolVar(&revokeJSON, "json", false, "output JSON")
	reactivateAgentCmd.Flags().StringVar(&revokeServer, "server", "", "Auth4Agent server URL")
	reactivateAgentCmd.MarkFlagRequired("agent-id")
	
	// List flags
	listRevokedCmd.Flags().BoolVar(&revokeJSON, "json", false, "output JSON")
	listRevokedCmd.Flags().StringVar(&revokeServer, "server", "", "Auth4Agent server URL")
	
	// Build command tree
	revokeCmd.AddCommand(revokeTokenCmd, revokeAgentCmd, suspendAgentCmd, reactivateAgentCmd, listRevokedCmd)
	rootCmd.AddCommand(revokeCmd)
}