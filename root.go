package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var quiet bool

var rootCmd = &cobra.Command{
	Use:   "auth4agent",
	Short: "Auth4Agent CLI",
	Long:  "CLI for managing Auth4Agent identities, tokens, and operators",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.auth4agent/config.json", "config file path")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress non-error output")
}
