package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initOperator bool
var initOutput string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate DID keypair and DID document",
	Long:  `Generate Ed25519 keypair and optional DID document for agent or operator.`,
	Run: func(cmd *cobra.Command, args []string) {
		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating keys: %v\n", err)
			os.Exit(1)
		}

		publicKeyB64 := base64.StdEncoding.EncodeToString(pub)
		privateKeyB64 := base64.StdEncoding.EncodeToString(priv)

		output := map[string]interface{}{
			"public_key":  publicKeyB64,
			"private_key": privateKeyB64,
		}

		if !initOperator {
			// Agent mode
			output["did_template"] = "did:agent:{operator_domain}:{suffix}"
			output["message"] = "Register this public_key with your operator to get a DID"
			
			// Save agent config
			cfg := &AgentConfig{
				PrivateKey: privateKeyB64,
				ServerURL:  "http://localhost:8080",
			}
			
			if existing, err := LoadAgentConfig(); err == nil {
				cfg.DID = existing.DID
				cfg.OperatorID = existing.OperatorID
				cfg.ServerURL = existing.ServerURL
			}
			
			if err := SaveAgentConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not save agent config: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "✓ Agent config saved to ~/.agentauth/agent.json\n")
			}
		} else {
			// Operator mode
			output["key_type"] = "operator_root"
			output["message"] = "Keep private_key secret. Use public_key to register with Auth4Agent."
			
			// Save operator config
			cfg := &OperatorConfig{
				RootPublicKey:  publicKeyB64,
				RootPrivateKey: privateKeyB64,
				ServerURL:      "http://localhost:8080",
			}
			
			if existing, err := LoadOperatorConfig(); err == nil {
				cfg.ID = existing.ID
				cfg.Domain = existing.Domain
				cfg.DomainVerifiedAt = existing.DomainVerifiedAt
				cfg.ServerURL = existing.ServerURL
			}
			
			if err := SaveOperatorConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not save operator config: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "✓ Operator config saved to ~/.agentauth/operator.json\n")
			}
		}

		var outData []byte
		outData, err = json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling output: %v\n", err)
			os.Exit(1)
		}

		if initOutput != "" {
			err = os.WriteFile(initOutput, outData, 0600)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Keys saved to %s\n", initOutput)
		} else {
			fmt.Println(string(outData))
		}
	},
}

func init() {
	initCmd.Flags().BoolVar(&initOperator, "operator", false, "generate operator root key (default: agent key)")
	initCmd.Flags().StringVar(&initOutput, "output", "", "output file path")
	rootCmd.AddCommand(initCmd)
}