package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "BVS gateway microservice",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("please use a subcommand")
	},
}

func ConfigureCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
