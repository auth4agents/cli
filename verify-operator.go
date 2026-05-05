// verify-operator.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var verifyOperatorCmd = &cobra.Command{
	Use:   "verify-operator",
	Short: "Verify operator domain ownership",
	Long:  "Get DNS TXT record instructions and verify domain ownership",
}

var verifyStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check domain verification status",
	Run: func(cmd *cobra.Command, args []string) {
		operatorID, _ := cmd.Flags().GetString("operator-id")
		server, _ := cmd.Flags().GetString("server")

		if operatorID == "" {
			fmt.Fprintf(os.Stderr, "Error: --operator-id is required\n")
			os.Exit(1)
		}

		client := NewClient(server)
		path := fmt.Sprintf("/api/v1/operators/%s", operatorID)

		var resp struct {
			ID               string  `json:"id"`
			Domain           string  `json:"domain"`
			DomainVerifiedAt *string `json:"domain_verified_at"`
			Status           string  `json:"status"`
		}

		err := client.Get(context.Background(), path, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if resp.DomainVerifiedAt != nil {
			fmt.Printf("✓ Domain verified: %s\n", resp.Domain)
			fmt.Printf("  Verified at: %s\n", *resp.DomainVerifiedAt)
		} else {
			fmt.Printf("✗ Domain not verified: %s\n", resp.Domain)
			fmt.Printf("  Run 'agentauth verify instructions --operator-id %s' to get DNS record\n", operatorID)
		}
	},
}

var verifyInstructionsCmd = &cobra.Command{
	Use:   "instructions",
	Short: "Get DNS TXT record instructions for domain verification",
	Run: func(cmd *cobra.Command, args []string) {
		operatorID, _ := cmd.Flags().GetString("operator-id")
		server, _ := cmd.Flags().GetString("server")

		if operatorID == "" {
			fmt.Fprintf(os.Stderr, "Error: --operator-id is required\n")
			os.Exit(1)
		}

		client := NewClient(server)
		path := fmt.Sprintf("/api/v1/operators/verify/%s", operatorID)

		var resp struct {
			Domain          string `json:"domain"`
			TXTRecordName   string `json:"txt_record_name"`
			TXTRecordValue  string `json:"txt_record_value"`
			Message         string `json:"message"`
		}

		err := client.Get(context.Background(), path, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n=== Domain Verification Instructions ===\n\n")
		fmt.Printf("Domain: %s\n\n", resp.Domain)
		fmt.Printf("1. Add this TXT record to your DNS:\n\n")
		fmt.Printf("   Name:  %s\n", resp.TXTRecordName)
		fmt.Printf("   Value: %s\n\n", resp.TXTRecordValue)
		fmt.Printf("2. Wait for DNS propagation (a few minutes)\n\n")
		fmt.Printf("3. Run: agentauth verify confirm --operator-id %s\n\n", operatorID)
	},
}

var verifyConfirmCmd = &cobra.Command{
	Use:   "confirm",
	Short: "Confirm domain verification by checking DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		operatorID, _ := cmd.Flags().GetString("operator-id")
		server, _ := cmd.Flags().GetString("server")

		if operatorID == "" {
			fmt.Fprintf(os.Stderr, "Error: --operator-id is required\n")
			os.Exit(1)
		}

		client := NewClient(server)
		path := fmt.Sprintf("/api/v1/operators/verify/%s", operatorID)

		var resp struct {
			Verified   bool   `json:"verified"`
			Message    string `json:"message"`
			VerifiedAt string `json:"verified_at"`
		}

		err := client.Post(context.Background(), path, nil, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if resp.Verified {
			fmt.Printf("✓ %s\n", resp.Message)
			if resp.VerifiedAt != "" {
				fmt.Printf("  Verified at: %s\n", resp.VerifiedAt)
			}
		} else {
			fmt.Printf("✗ %s\n", resp.Message)
			os.Exit(1)
		}
	},
}

func init() {
	// Status command
	verifyStatusCmd.Flags().String("operator-id", "", "Operator ID")
	verifyStatusCmd.Flags().String("server", "http://localhost:8080", "AgentAuth server URL")
	verifyStatusCmd.MarkFlagRequired("operator-id")

	// Instructions command
	verifyInstructionsCmd.Flags().String("operator-id", "", "Operator ID")
	verifyInstructionsCmd.Flags().String("server", "http://localhost:8080", "AgentAuth server URL")
	verifyInstructionsCmd.MarkFlagRequired("operator-id")

	// Confirm command
	verifyConfirmCmd.Flags().String("operator-id", "", "Operator ID")
	verifyConfirmCmd.Flags().String("server", "http://localhost:8080", "AgentAuth server URL")
	verifyConfirmCmd.MarkFlagRequired("operator-id")

	verifyOperatorCmd.AddCommand(verifyStatusCmd, verifyInstructionsCmd, verifyConfirmCmd)
	rootCmd.AddCommand(verifyOperatorCmd)
}