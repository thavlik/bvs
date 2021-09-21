package main

import (
	"github.com/spf13/cobra"
	"github.com/thavlik/bvs/components/db-sync/pkg/db_sync"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RootArgs struct {
	NodeAddr string
}

var rootArgs RootArgs

var rootCmd = &cobra.Command{
	Use: "db-sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		k8sConfig, err := rest.InClusterConfig()
		if err != nil {
			return err
		}
		cl, err := client.New(k8sConfig, client.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			return err
		}
		return db_sync.Start(rootArgs.NodeAddr, cl)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootArgs.NodeAddr, "node-addr", "localhost:2100", "cardano-node socat address")
}
