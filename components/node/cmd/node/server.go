package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thavlik/bvs/components/node/pkg/node"
)

type ServerArgs struct {
	Port         int
	Config       string
	DatabasePath string
	SocketPath   string
	HostAddr     string
	Topology     string
}

var serverArgs ServerArgs

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if serverArgs.Config == "" {
			return fmt.Errorf("missing --config")
		}
		if serverArgs.DatabasePath == "" {
			return fmt.Errorf("missing --database-path")
		}
		if serverArgs.SocketPath == "" {
			return fmt.Errorf("missing --socket-path")
		}
		if serverArgs.HostAddr == "" {
			return fmt.Errorf("missing --host-addr")
		}
		if serverArgs.Topology == "" {
			return fmt.Errorf("missing --topology")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return node.NewServer(log).Start(
			serverArgs.Port,
			serverArgs.Config,
			serverArgs.DatabasePath,
			serverArgs.SocketPath,
			serverArgs.HostAddr,
			serverArgs.Topology,
		)
	},
}

func init() {
	serverCmd.PersistentFlags().IntVarP(&serverArgs.Port, "port", "p", 80, "listener port")
	serverCmd.PersistentFlags().StringVar(&serverArgs.Config, "config", "", "config path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.DatabasePath, "database-path", "", "database path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.SocketPath, "socket-path", "", "socket path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.HostAddr, "host-addr", "", "host address")
	serverCmd.PersistentFlags().StringVar(&serverArgs.Topology, "topology", "", "topology file path")
	ConfigureCommand(serverCmd)
}
