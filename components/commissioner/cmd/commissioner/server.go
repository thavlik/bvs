package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thavlik/bvs/components/commissioner/pkg/commissioner"
	"github.com/thavlik/bvs/components/commissioner/pkg/storage/memory"
)

type ServerArgs struct {
	Port             int
	NodePort         int
	NodeConfig       string
	NodeDatabasePath string
	NodeSocketPath   string
	NodeHostAddr     string
	NodeTopology     string
	TokenName        string
}

var serverArgs ServerArgs

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if serverArgs.TokenName == "" {
			return fmt.Errorf("missing --token-name")
		}
		if serverArgs.NodeConfig == "" {
			return fmt.Errorf("missing --node-config")
		}
		if serverArgs.NodeDatabasePath == "" {
			return fmt.Errorf("missing --node-database-path")
		}
		if serverArgs.NodeSocketPath == "" {
			return fmt.Errorf("missing --node-socket-path")
		}
		if serverArgs.NodeHostAddr == "" {
			return fmt.Errorf("missing --node-host-addr")
		}
		if serverArgs.NodeTopology == "" {
			return fmt.Errorf("missing --node-topology")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return commissioner.NewServer(
			serverArgs.TokenName,
			memory.NewMemoryStorage(),
			log,
		).Start(
			serverArgs.Port,
			serverArgs.NodePort,
			serverArgs.NodeConfig,
			serverArgs.NodeDatabasePath,
			serverArgs.NodeSocketPath,
			serverArgs.NodeHostAddr,
			serverArgs.NodeTopology,
		)
	},
}

func init() {
	serverCmd.PersistentFlags().IntVar(&serverArgs.Port, "port", 80, "http service listener port")
	serverCmd.PersistentFlags().IntVar(&serverArgs.NodePort, "node-port", 1337, "cardano-node listener port")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeConfig, "node-config", "", "cardano-node config path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeDatabasePath, "node-database-path", "", "cardano-node database path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeSocketPath, "node-socket-path", "", "cardano-node socket path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeHostAddr, "node-host-addr", "", "cardano-node host address")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeTopology, "node-topology", "", "cardano-node topology file path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.TokenName, "token-name", "", "cardano NFT name")
	ConfigureCommand(serverCmd)
}
