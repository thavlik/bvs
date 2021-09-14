package main

import (
	"github.com/spf13/cobra"
	"github.com/thavlik/bvs/components/gateway/pkg/gateway"
)

type ServerArgs struct {
	Port int
}

var serverArgs ServerArgs

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return gateway.NewServer(
			log,
		).Listen(serverArgs.Port)
	},
}

func init() {
	serverCmd.PersistentFlags().IntVarP(&serverArgs.Port, "port", "p", 80, "listener port")
	ConfigureCommand(serverCmd)
}
