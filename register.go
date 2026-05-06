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
	Short: "Register operator or agent identity",
}

var registerOperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Register operator identity",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := LoadOperatorConfig()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"operator config not found: run init first\n",
			)

			os.Exit(1)
		}

		jsonOut, _ := cmd.Flags().GetBool("json")

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		defer cancel()

		client := NewClient(cfg.ServerURL)

		req := RegisterOperatorRequest{
			Domain:        cfg.Domain,
			RootPublicKey: cfg.RootPublicKey,
		}

		resp, err := client.RegisterOperator(
			ctx,
			req,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"registration failed: %v\n",
				err,
			)

			os.Exit(1)
		}

		cfg.ID = resp.ID

		err = SaveOperatorConfig(cfg)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"warning: failed to save config: %v\n",
				err,
			)
		}

		if jsonOut {
			_ = json.NewEncoder(os.Stdout).
				Encode(resp)

			return
		}

		fmt.Printf(
			"operator registered\n\n",
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

		fmt.Printf(
			"\nnext:\n",
		)

		fmt.Printf(
			"auth4agents verify-operator instructions --operator-id %s\n",
			resp.ID,
		)
	},
}

var registerAgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Register agent DID",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := LoadAgentConfig()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"agent config not found: run init first\n",
			)

			os.Exit(1)
		}

		if cfg.DIDDocument == nil {
			fmt.Fprintf(
				os.Stderr,
				"missing DID document\n",
			)

			os.Exit(1)
		}

		operatorID, _ := cmd.Flags().
			GetString("operator-id")

		if operatorID == "" {
			fmt.Fprintf(
				os.Stderr,
				"--operator-id is required\n",
			)

			os.Exit(1)
		}

		jsonOut, _ := cmd.Flags().
			GetBool("json")

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		defer cancel()

		client := NewClient(
			cfg.ServerURL,
		)

		req := RegisterAgentRequest{
			DID:            cfg.DID,
			DIDDocument:    cfg.DIDDocument,
			AgentPublicKey: cfg.PublicKey,
			OperatorID:     operatorID,
		}

		resp, err := client.RegisterAgent(
			ctx,
			req,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"registration failed: %v\n",
				err,
			)

			os.Exit(1)
		}

		cfg.OperatorID = operatorID

		err = SaveAgentConfig(cfg)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"warning: failed to save config: %v\n",
				err,
			)
		}

		if jsonOut {
			_ = json.NewEncoder(os.Stdout).
				Encode(resp)

			return
		}

		fmt.Printf(
			"agent registered\n\n",
		)

		fmt.Printf(
			"did: %s\n",
			resp.DID,
		)

		fmt.Printf(
			"status: %s\n",
			resp.Status,
		)
	},
}

func init() {

	registerOperatorCmd.Flags().
		Bool(
			"json",
			false,
			"output as JSON",
		)

	registerAgentCmd.Flags().
		String(
			"operator-id",
			"",
			"operator ID",
		)

	registerAgentCmd.Flags().
		Bool(
			"json",
			false,
			"output as JSON",
		)

	registerCmd.AddCommand(
		registerOperatorCmd,
		registerAgentCmd,
	)

	rootCmd.AddCommand(registerCmd)
}
