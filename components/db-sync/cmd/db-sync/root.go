package main

import (
	"github.com/spf13/cobra"
	"github.com/thavlik/bvs/components/db-sync/pkg/db_sync"
)

type RootArgs struct {
	NodeAddr string
}

var rootArgs RootArgs

var rootCmd = &cobra.Command{
	Use: "db-sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		return db_sync.Start(rootArgs.NodeAddr)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootArgs.NodeAddr, "node-addr", "localhost:2100", "cardano-node socat address")
}
