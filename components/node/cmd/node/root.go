package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "node",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("please use a subcommand")
	},
}

func ConfigureCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
