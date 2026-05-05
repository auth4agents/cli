// register.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var regServer string

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register operator or agent with AgentAuth server",
}

var registerOperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Register an operator",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		publicKey, _ := cmd.Flags().GetString("public-key")
		server, _ := cmd.Flags().GetString("server")

		client := NewClient(server)
		
		req := struct {
			Domain        string `json:"domain"`
			RootPublicKey string `json:"root_public_key"`
		}{Domain: domain, RootPublicKey: publicKey}

		var resp struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
			Status string `json:"status"`
		}

		err := client.Post(context.Background(), "/api/v1/operators", req, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Operator registered: %s\n", resp.ID)
		fmt.Print("NOTE: Save this operator ID. Use the returned operator ID to register agents under this operator.\n")
	},
}

var registerAgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Register an agent under an operator",
	Run: func(cmd *cobra.Command, args []string) {
		operatorID, _ := cmd.Flags().GetString("operator-id")
		publicKey, _ := cmd.Flags().GetString("public-key")
		server, _ := cmd.Flags().GetString("server")

		if operatorID == "" {
			fmt.Fprintf(os.Stderr, "Error: --operator-id is required\n")
			os.Exit(1)
		}
		if publicKey == "" {
			fmt.Fprintf(os.Stderr, "Error: --public-key is required\n")
			os.Exit(1)
		}

		client := NewClient(server)

		req := struct {
			AgentPublicKey string `json:"agent_public_key"`
		}{AgentPublicKey: publicKey}

		var resp struct {
			DID    string `json:"did"`
			Status string `json:"status"`
		}

		path := fmt.Sprintf("/api/v1/operators/%s/agents", operatorID)
		err := client.Post(context.Background(), path, req, &resp)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintf(os.Stderr, "Failed to register agent. Possible reason:\n\t1. Operator not found \n \t2. Invalid public key format \n\t3. Server error\n\t4. Operator domain not verified\n")
			os.Exit(1)
		}

		if resp.DID == "" {
			fmt.Fprintf(os.Stderr, "Error: No DID returned from server\n")
			os.Exit(1)
		}

		fmt.Printf("✓ Agent registered successfully\n")
		fmt.Printf("  DID: %s\n", resp.DID)
		fmt.Printf("  Status: %s\n", resp.Status)
	},
}

func init() {
	// Operator flags
	registerOperatorCmd.Flags().String("domain", "", "operator domain")
	registerOperatorCmd.Flags().String("public-key", "", "public key (base64)")
	registerOperatorCmd.MarkFlagRequired("domain")
	registerOperatorCmd.MarkFlagRequired("public-key")

	// Agent flags
	registerAgentCmd.Flags().String("operator-id", "", "operator ID")
	registerAgentCmd.Flags().String("public-key", "", "public key (base64)")
	registerAgentCmd.MarkFlagRequired("operator-id")
	registerAgentCmd.MarkFlagRequired("public-key")

	// Global flag for both
	registerOperatorCmd.Flags().String("server", "http://localhost:8080", "AgentAuth server URL")
	registerAgentCmd.Flags().String("server", "http://localhost:8080", "AgentAuth server URL")

	registerCmd.AddCommand(registerOperatorCmd, registerAgentCmd)
	rootCmd.AddCommand(registerCmd)
}