// register.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register operator or agent with auth4agents server",
}

var registerOperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Register an operator",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		publicKey, _ := cmd.Flags().GetString("public-key")
		server, _ := cmd.Flags().GetString("server")
		jsonOut, _ := cmd.Flags().GetBool("json")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client := NewClient(server)

		req := struct {
			Domain        string `json:"domain"`
			RootPublicKey string `json:"root_public_key"`
		}{
			Domain:        domain,
			RootPublicKey: publicKey,
		}

		var resp struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
			Status string `json:"status"`
		}

		if err := client.Post(ctx, "/api/v1/operators", req, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
			os.Exit(1)
		}

		if resp.ID == "" {
			fmt.Fprintf(os.Stderr, "Invalid server response: missing operator ID\n")
			os.Exit(1)
		}

		if jsonOut {
			_ = json.NewEncoder(os.Stdout).Encode(resp)
			return
		}

		fmt.Printf("✓ Operator registered\n")
		fmt.Printf("  ID: %s\n", resp.ID)
		fmt.Printf("  Domain: %s\n", resp.Domain)
		fmt.Printf("  Status: %s\n", resp.Status)

		fmt.Printf("\nNext:\n")
		fmt.Printf("  auth4agents verify-operator instructions --operator-id %s\n", resp.ID)
		fmt.Printf("  auth4agents verify-operator confirm --operator-id %s\n", resp.ID)
	},
}

var registerAgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Register an agent under an operator",
	Run: func(cmd *cobra.Command, args []string) {
		operatorID, _ := cmd.Flags().GetString("operator-id")
		publicKey, _ := cmd.Flags().GetString("public-key")
		server, _ := cmd.Flags().GetString("server")
		jsonOut, _ := cmd.Flags().GetBool("json")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client := NewClient(server)

		req := struct {
			AgentPublicKey string `json:"agent_public_key"`
		}{
			AgentPublicKey: publicKey,
		}

		var resp struct {
			DID    string `json:"did"`
			Status string `json:"status"`
		}

		path := fmt.Sprintf("/api/v1/operators/%s/agents", operatorID)

		if err := client.Post(ctx, path, req, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
			os.Exit(1)
		}

		if resp.DID == "" {
			fmt.Fprintf(os.Stderr, "Invalid server response: missing DID\n")
			os.Exit(1)
		}

		if jsonOut {
			_ = json.NewEncoder(os.Stdout).Encode(resp)
			return
		}

		fmt.Printf("✓ Agent registered\n")
		fmt.Printf("  DID: %s\n", resp.DID)
		fmt.Printf("  Status: %s\n", resp.Status)
	},
}

func init() {
	registerOperatorCmd.Flags().String("domain", "", "operator domain")
	registerOperatorCmd.Flags().StringP("public-key", "k", "", "public key (base64)")
	registerOperatorCmd.Flags().String("server", "http://localhost:8080", "auth4agents server URL")
	registerOperatorCmd.Flags().Bool("json", false, "output as JSON")
	registerOperatorCmd.MarkFlagRequired("domain")
	registerOperatorCmd.MarkFlagRequired("public-key")

	registerAgentCmd.Flags().String("operator-id", "", "operator ID")
	registerAgentCmd.Flags().StringP("public-key", "k", "", "public key (base64)")
	registerAgentCmd.Flags().String("server", "http://localhost:8080", "auth4agents server URL")
	registerAgentCmd.Flags().Bool("json", false, "output as JSON")
	registerAgentCmd.MarkFlagRequired("operator-id")
	registerAgentCmd.MarkFlagRequired("public-key")

	registerCmd.AddCommand(registerOperatorCmd, registerAgentCmd)
	rootCmd.AddCommand(registerCmd)
}