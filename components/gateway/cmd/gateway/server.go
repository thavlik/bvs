package main

import (
	"github.com/spf13/cobra"
	"github.com/thavlik/bvs/components/gateway/pkg/gateway"
)

type ServerArgs struct {
	Port        int
	MetricsPort int
}

var serverArgs ServerArgs

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return gateway.NewServer().Start(
			serverArgs.Port,
			serverArgs.MetricsPort,
		)
	},
}

func init() {
	serverCmd.PersistentFlags().IntVar(&serverArgs.Port, "port", 80, "http service listener port")
	serverCmd.PersistentFlags().IntVar(&serverArgs.MetricsPort, "metrics-port", 0, "optional prometheus metrics port")
	ConfigureCommand(serverCmd)
}
