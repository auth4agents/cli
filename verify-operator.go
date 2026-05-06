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

var verifyOperatorStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check operator verification status",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := LoadOperatorConfig()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"operator config not found\n",
			)

			os.Exit(1)
		}

		if cfg.ID == "" {
			fmt.Fprintf(
				os.Stderr,
				"operator not registered\n",
			)

			os.Exit(1)
		}

		jsonOut, _ := cmd.Flags().
			GetBool("json")

		client := NewClient(
			cfg.ServerURL,
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		defer cancel()

		path := fmt.Sprintf(
			"/v1/operators/%s",
			cfg.ID,
		)

		var resp struct {
			ID               string  `json:"id"`
			Domain           string  `json:"domain"`
			Status           string  `json:"status"`
			DomainVerifiedAt *string `json:"domain_verified_at"`
		}

		err = client.Get(
			ctx,
			path,
			&resp,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"request failed: %v\n",
				err,
			)

			os.Exit(1)
		}

		if resp.DomainVerifiedAt != nil {
			cfg.DomainVerifiedAt = *resp.DomainVerifiedAt

			_ = SaveOperatorConfig(cfg)
		}

		if jsonOut {

			_ = json.NewEncoder(os.Stdout).
				Encode(resp)

			return
		}

		fmt.Printf(
			"operator verification status\n\n",
		)

		fmt.Printf(
			"id: %s\n",
			resp.ID,
		)

		fmt.Printf(
			"domain: %s\n",
			resp.Domain,
		)

		fmt.Printf(
			"status: %s\n",
			resp.Status,
		)

		if resp.DomainVerifiedAt != nil {

			fmt.Printf(
				"verified_at: %s\n",
				*resp.DomainVerifiedAt,
			)

			return
		}

		fmt.Printf(
			"verified: false\n",
		)
	},
}

var verifyOperatorInstructionsCmd = &cobra.Command{
	Use:   "instructions",
	Short: "Get DNS TXT verification instructions",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := LoadOperatorConfig()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"operator config not found\n",
			)

			os.Exit(1)
		}

		if cfg.ID == "" {
			fmt.Fprintf(
				os.Stderr,
				"operator not registered\n",
			)

			os.Exit(1)
		}

		jsonOut, _ := cmd.Flags().
			GetBool("json")

		client := NewClient(
			cfg.ServerURL,
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		defer cancel()

		path := fmt.Sprintf(
			"/v1/operators/verify/%s",
			cfg.ID,
		)

		var resp struct {
			Domain         string `json:"domain"`
			TXTRecordName  string `json:"txt_record_name"`
			TXTRecordValue string `json:"txt_record_value"`
			Message        string `json:"message"`
		}

		err = client.Get(
			ctx,
			path,
			&resp,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"request failed: %v\n",
				err,
			)

			os.Exit(1)
		}

		if jsonOut {

			_ = json.NewEncoder(os.Stdout).
				Encode(resp)

			return
		}

		fmt.Printf(
			"dns verification instructions\n\n",
		)

		fmt.Printf(
			"domain: %s\n\n",
			resp.Domain,
		)

		fmt.Printf(
			"record_type: TXT\n",
		)

		fmt.Printf(
			"name: %s\n",
			resp.TXTRecordName,
		)

		fmt.Printf(
			"value: %s\n\n",
			resp.TXTRecordValue,
		)

		fmt.Printf(
			"after propagation run:\n",
		)

		fmt.Printf(
			"auth4agents verify-operator confirm\n",
		)
	},
}

var verifyOperatorConfirmCmd = &cobra.Command{
	Use:   "confirm",
	Short: "Confirm DNS verification",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := LoadOperatorConfig()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"operator config not found\n",
			)

			os.Exit(1)
		}

		if cfg.ID == "" {
			fmt.Fprintf(
				os.Stderr,
				"operator not registered\n",
			)

			os.Exit(1)
		}

		jsonOut, _ := cmd.Flags().
			GetBool("json")

		client := NewClient(
			cfg.ServerURL,
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			15*time.Second,
		)

		defer cancel()

		path := fmt.Sprintf(
			"/v1/operators/verify/%s",
			cfg.ID,
		)

		var resp struct {
			Verified   bool   `json:"verified"`
			Message    string `json:"message"`
			VerifiedAt string `json:"verified_at"`
		}

		err = client.post(
			ctx,
			path,
			nil,
			&resp,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"verification failed: %v\n",
				err,
			)

			os.Exit(1)
		}

		if resp.Verified {

			cfg.DomainVerifiedAt = resp.VerifiedAt

			_ = SaveOperatorConfig(cfg)
		}

		if jsonOut {

			_ = json.NewEncoder(os.Stdout).
				Encode(resp)

			return
		}

		if !resp.Verified {

			fmt.Printf(
				"verification failed\n\n",
			)

			fmt.Printf(
				"%s\n",
				resp.Message,
			)

			os.Exit(1)
		}

		fmt.Printf(
			"operator verified\n\n",
		)

		fmt.Printf(
			"verified_at: %s\n",
			resp.VerifiedAt,
		)
	},
}

func init() {

	verifyOperatorStatusCmd.Flags().
		Bool(
			"json",
			false,
			"output JSON",
		)

	verifyOperatorInstructionsCmd.Flags().
		Bool(
			"json",
			false,
			"output JSON",
		)

	verifyOperatorConfirmCmd.Flags().
		Bool(
			"json",
			false,
			"output JSON",
		)

	verifyOperatorCmd.AddCommand(
		verifyOperatorStatusCmd,
		verifyOperatorInstructionsCmd,
		verifyOperatorConfirmCmd,
	)

	rootCmd.AddCommand(
		verifyOperatorCmd,
	)
}
