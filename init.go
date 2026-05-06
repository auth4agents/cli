package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initOperator bool
var initOutput string
var initDomain string
var initServer string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate identity, keys, and DID document",

	Run: func(cmd *cobra.Command, args []string) {

		pub, priv, err := GenerateEd25519Keypair()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to generate keypair: %v\n",
				err,
			)

			os.Exit(1)
		}

		publicKeyB64 := EncodePublicKey(pub)
		privateKeyB64 := EncodePrivateKey(priv)

		if initOperator {

			err := initializeOperator(
				publicKeyB64,
				privateKeyB64,
			)

			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"failed to initialize operator: %v\n",
					err,
				)

				os.Exit(1)
			}

			return
		}

		err = initializeAgent(
			pub,
			publicKeyB64,
			privateKeyB64,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to initialize agent: %v\n",
				err,
			)

			os.Exit(1)
		}
	},
}

func initializeAgent(
	publicKey []byte,
	publicKeyB64 string,
	privateKeyB64 string,
) error {

	if initDomain == "" {
		return fmt.Errorf(
			"--domain is required for agent initialization",
		)
	}

	did := GenerateDID(
		initDomain,
		publicKey,
	)

	doc := GenerateDIDDocument(
		did,
		publicKey,
	)

	keyPath, err := SavePrivateKeyFile(
		"agent.key",
		privateKeyB64,
	)

	if err != nil {
		return err
	}

	cfg := &AgentConfig{
		DID:       did,
		ServerURL: initServer,
		PublicKey: publicKeyB64,

		KeyConfig: KeyConfig{
			Provider: DefaultKeyProvider,
			Ref:      keyPath,
		},

		DIDDocument: doc,
	}

	err = SaveAgentConfig(cfg)
	if err != nil {
		return err
	}

	output := map[string]interface{}{
		"mode":         "agent",
		"did":          did,
		"public_key":   publicKeyB64,
		"did_document": doc,
		"server":       initServer,
	}

	return writeInitOutput(output)
}

func initializeOperator(
	publicKeyB64 string,
	privateKeyB64 string,
) error {

	keyPath, err := SavePrivateKeyFile(
		"operator.key",
		privateKeyB64,
	)

	if err != nil {
		return err
	}

	cfg := &OperatorConfig{
		Domain:        NormalizeDomain(initDomain),
		ServerURL:     initServer,
		RootPublicKey: publicKeyB64,

		KeyConfig: KeyConfig{
			Provider: DefaultKeyProvider,
			Ref:      keyPath,
		},
	}

	err = SaveOperatorConfig(cfg)
	if err != nil {
		return err
	}

	output := map[string]interface{}{
		"mode":            "operator",
		"domain":          cfg.Domain,
		"root_public_key": publicKeyB64,
		"server":          initServer,
	}

	return writeInitOutput(output)
}

func writeInitOutput(
	output interface{},
) error {

	data, err := json.MarshalIndent(
		output,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	if initOutput != "" {

		err = os.WriteFile(
			initOutput,
			data,
			0600,
		)

		if err != nil {
			return err
		}

		fmt.Printf(
			"identity written to %s\n",
			initOutput,
		)

		return nil
	}

	fmt.Println(string(data))

	return nil
}

func init() {

	initCmd.Flags().BoolVar(
		&initOperator,
		"operator",
		false,
		"initialize operator identity",
	)

	initCmd.Flags().StringVar(
		&initDomain,
		"domain",
		"",
		"operator domain",
	)

	initCmd.Flags().StringVar(
		&initServer,
		"server",
		"http://localhost:8080",
		"Auth4Agents server URL",
	)

	initCmd.Flags().StringVar(
		&initOutput,
		"output",
		"",
		"write output to file",
	)

	rootCmd.AddCommand(initCmd)
}
