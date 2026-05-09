package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var issueScope string
var issueAudience string
var issueTTL string
var issueServer string
var issueJSON bool

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Exchange DID proof for scoped JWT",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := LoadAgentConfig()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"agent config not found: run init first\n",
			)

			os.Exit(1)
		}

		if cfg.DID == "" {
			fmt.Fprintf(
				os.Stderr,
				"missing DID in config\n",
			)

			os.Exit(1)
		}

		privateKeyB64, err := LoadPrivateKeyFile(
			cfg.KeyConfig.Ref,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to load private key: %v\n",
				err,
			)

			os.Exit(1)
		}

		privateKey, err := DecodePrivateKey(
			privateKeyB64,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to decode private key: %v\n",
				err,
			)

			os.Exit(1)
		}

		serverURL := cfg.ServerURL

		if issueServer != "" {
			serverURL = issueServer
		}

		client := NewClient(serverURL)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			15*time.Second,
		)

		defer cancel()

		challengeResp, err := client.GetChallenge(
			ctx,
			cfg.DID,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"challenge request failed: %v\n",
				err,
			)

			os.Exit(1)
		}

		signature, err := SignChallenge(
			privateKey,
			challengeResp.Challenge,
			challengeResp.Nonce,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to sign challenge: %v\n",
				err,
			)

			os.Exit(1)
		}

		proof := NewProof(
			cfg.DID,
			challengeResp.Challenge,
			challengeResp.Nonce,
			signature,
		)

		req := BuildTokenExchangeRequest(
			cfg.DID,
			issueScope,
			issueAudience,
			issueTTL,
			proof,
		)

		tokenResp, err := client.ExchangeToken(ctx, TokenRequest{
			DID:       req.DID,
			Challenge: req.Proof.Challenge,
			Nonce:     req.Proof.Nonce,
			Signature: req.Proof.Signature,
			Scope:     req.Scope,
			Audience:  req.Audience,
			TTL:       req.TTL,
		})

		if err != nil {
			// Enhanced error messages for new status codes
			errMsg := err.Error()

			if strings.Contains(errMsg, "410") || strings.Contains(errMsg, "challenge expired") {
				fmt.Fprintf(os.Stderr, "Challenge expired. Please request a new challenge.\n")
			} else if strings.Contains(errMsg, "409") || strings.Contains(errMsg, "already used") {
				fmt.Fprintf(os.Stderr, "Challenge already used. Please request a new challenge.\n")
			} else if strings.Contains(errMsg, "404") || strings.Contains(errMsg, "not found") {
				fmt.Fprintf(os.Stderr, "Challenge not found. Please request a new challenge.\n")
			} else if strings.Contains(errMsg, "agent is not active") || strings.Contains(errMsg, "agent has been revoked") {
				fmt.Fprintf(os.Stderr, "Agent is revoked or suspended. Cannot issue token.\n")
			} else if strings.Contains(errMsg, "scope not allowed") {
				fmt.Fprintf(os.Stderr, "Scope '%s' is not allowed for this agent.\n", issueScope)
			} else {
				fmt.Fprintf(os.Stderr, "Token exchange failed: %v\n", err)
			}

			os.Exit(1)
		}

		if issueJSON {

			output := map[string]interface{}{
				"token":      tokenResp.Token,
				"token_type": tokenResp.TokenType,
				"expires_at": tokenResp.ExpiresAt,
			}

			_ = json.NewEncoder(os.Stdout).
				Encode(output)

			return
		}

		fmt.Printf(
			"token issued\n\n",
		)

		fmt.Printf(
			"did: %s\n",
			cfg.DID,
		)

		fmt.Printf(
			"scope: %s\n",
			issueScope,
		)

		fmt.Printf(
			"audience: %s\n",
			issueAudience,
		)

		fmt.Printf(
			"expires_at: %s\n\n",
			tokenResp.ExpiresAt,
		)

		fmt.Printf(
			"%s\n",
			tokenResp.Token,
		)
	},
}

func init() {

	issueCmd.Flags().
		StringVar(
			&issueScope,
			"scope",
			"",
			"requested scope",
		)

	issueCmd.Flags().
		StringVar(
			&issueAudience,
			"aud",
			"",
			"target audience",
		)

	issueCmd.Flags().
		StringVar(
			&issueTTL,
			"ttl",
			"1h",
			"token lifetime",
		)

	issueCmd.Flags().
		StringVar(
			&issueServer,
			"server",
			"",
			"Auth4Agents server URL",
		)

	issueCmd.Flags().
		BoolVar(
			&issueJSON,
			"json",
			false,
			"output JSON",
		)

	_ = issueCmd.MarkFlagRequired(
		"scope",
	)

	_ = issueCmd.MarkFlagRequired(
		"aud",
	)

	rootCmd.AddCommand(issueCmd)
}
