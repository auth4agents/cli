// verify-operator.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var verifyOperatorCmd = &cobra.Command{
	Use:   "verify-operator",
	Short: "Verify operator domain ownership",
}

var verifyStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check domain verification status",
	Run: func(cmd *cobra.Command, args []string) {
		operatorID, _ := cmd.Flags().GetString("operator-id")
		server, _ := cmd.Flags().GetString("server")
		jsonOut, _ := cmd.Flags().GetBool("json")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client := NewClient(server)
		path := fmt.Sprintf("/api/v1/operators/%s", operatorID)

		var resp struct {
			ID               string  `json:"id"`
			Domain           string  `json:"domain"`
			DomainVerifiedAt *string `json:"domain_verified_at"`
			Status           string  `json:"status"`
		}

		if err := client.Get(ctx, path, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
			os.Exit(1)
		}

		if jsonOut {
			_ = json.NewEncoder(os.Stdout).Encode(resp)
			return
		}

		if resp.DomainVerifiedAt != nil {
			fmt.Printf("✓ Domain verified\n")
			fmt.Printf("  Domain: %s\n", resp.Domain)
			fmt.Printf("  Verified at: %s\n", *resp.DomainVerifiedAt)
		} else {
			fmt.Printf("✗ Domain not verified\n")
			fmt.Printf("  Domain: %s\n", resp.Domain)
			fmt.Printf("  Next:\n")
			fmt.Printf("    auth4agents verify-operator instructions --operator-id %s\n", operatorID)
		}
	},
}

var verifyInstructionsCmd = &cobra.Command{
	Use:   "instructions",
	Short: "Get DNS TXT record instructions",
	Run: func(cmd *cobra.Command, args []string) {
		operatorID, _ := cmd.Flags().GetString("operator-id")
		server, _ := cmd.Flags().GetString("server")
		jsonOut, _ := cmd.Flags().GetBool("json")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client := NewClient(server)
		path := fmt.Sprintf("/api/v1/operators/verify/%s", operatorID)

		var resp struct {
			Domain         string `json:"domain"`
			TXTRecordName  string `json:"txt_record_name"`
			TXTRecordValue string `json:"txt_record_value"`
			Message        string `json:"message"`
		}
		println(resp.TXTRecordValue)
		if err := client.Get(ctx, path, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
			os.Exit(1)
		}

		if jsonOut {
			_ = json.NewEncoder(os.Stdout).Encode(resp)
			return
		}

		fmt.Printf("Domain: %s\n", resp.Domain)
		fmt.Printf("\nAdd DNS TXT record:\n")
		fmt.Printf("  Name:  %s\n", resp.TXTRecordName)
		fmt.Printf("  Value: %s\n", resp.TXTRecordValue)
		fmt.Printf("\nThen run:\n")
		fmt.Printf("  auth4agents verify-operator confirm --operator-id %s\n", operatorID)
	},
}

var verifyConfirmCmd = &cobra.Command{
	Use:   "confirm",
	Short: "Confirm domain verification",
	Run: func(cmd *cobra.Command, args []string) {
		operatorID, _ := cmd.Flags().GetString("operator-id")
		server, _ := cmd.Flags().GetString("server")
		jsonOut, _ := cmd.Flags().GetBool("json")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client := NewClient(server)
		path := fmt.Sprintf("/api/v1/operators/verify/%s", operatorID)

		var resp struct {
			Verified   bool   `json:"verified"`
			Message    string `json:"message"`
			VerifiedAt string `json:"verified_at"`
		}

		if err := client.Post(ctx, path, nil, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
			os.Exit(1)
		}

		if jsonOut {
			_ = json.NewEncoder(os.Stdout).Encode(resp)
			return
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
	verifyStatusCmd.Flags().String("operator-id", "", "operator ID")
	verifyStatusCmd.Flags().String("server", "http://localhost:8080", "auth4agents server URL")
	verifyStatusCmd.Flags().Bool("json", false, "output as JSON")
	verifyStatusCmd.MarkFlagRequired("operator-id")

	verifyInstructionsCmd.Flags().String("operator-id", "", "operator ID")
	verifyInstructionsCmd.Flags().String("server", "http://localhost:8080", "auth4agents server URL")
	verifyInstructionsCmd.Flags().Bool("json", false, "output as JSON")
	verifyInstructionsCmd.MarkFlagRequired("operator-id")

	verifyConfirmCmd.Flags().String("operator-id", "", "operator ID")
	verifyConfirmCmd.Flags().String("server", "http://localhost:8080", "auth4agents server URL")
	verifyConfirmCmd.Flags().Bool("json", false, "output as JSON")
	verifyConfirmCmd.MarkFlagRequired("operator-id")

	verifyOperatorCmd.AddCommand(verifyStatusCmd, verifyInstructionsCmd, verifyConfirmCmd)
	rootCmd.AddCommand(verifyOperatorCmd)
}