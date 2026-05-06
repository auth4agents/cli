package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var scopesCmd = &cobra.Command{
	Use:   "agent-scopes",
	Short: "Manage agent scopes",
}

var scopesSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set allowed scopes",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := LoadAgentConfig()

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to load config: %v\n",
				err,
			)

			os.Exit(1)
		}

		operatorCfg, err := LoadOperatorConfig()

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to load operator config: %v\n",
				err,
			)

			os.Exit(1)
		}

		scopes, _ := cmd.Flags().
			GetStringSlice("scopes")

		client := NewClient(
			cfg.ServerURL,
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			15*time.Second,
		)

		defer cancel()

		path := fmt.Sprintf(
			"/v1/operators/%s/agents/%s/scopes",
			operatorCfg.ID,
			cfg.DID,
		)

		payload := map[string]interface{}{
			"allowed_scopes": scopes,
		}

		var out map[string]interface{}

		err = client.Patch(
			ctx,
			path,
			payload,
			&out,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"request failed: %v\n",
				err,
			)

			os.Exit(1)
		}

		data, _ := json.MarshalIndent(
			out,
			"",
			"  ",
		)

		fmt.Println(string(data))
	},
}

func init() {

	scopesSetCmd.Flags().
		StringSlice(
			"scopes",
			[]string{},
			"allowed scopes",
		)

	scopesCmd.AddCommand(
		scopesSetCmd,
	)

	rootCmd.AddCommand(
		scopesCmd,
	)
}