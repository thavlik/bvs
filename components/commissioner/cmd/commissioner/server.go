package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thavlik/bvs/components/commissioner/pkg/commissioner"
	"github.com/thavlik/bvs/components/commissioner/pkg/storage/mongodb_storage"
)

type ServerArgs struct {
	Port             int
	MetricsPort      int
	NodePort         int
	NodeConfig       string
	NodeDatabasePath string
	NodeSocketPath   string
	NodeHostAddr     string
	NodeTopology     string
	TokenName        string
	MongoDBHost      string
	MongoDBPort      int
	MongoDBDatabase  string
	MongoDBUsername  string
	MongoDBPassword  string
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
		if serverArgs.MongoDBHost == "" {
			return fmt.Errorf("missing --mongodb-host")
		}
		if v, ok := os.LookupEnv("MONGODB_USERNAME"); ok {
			serverArgs.MongoDBUsername = v
		}
		if v, ok := os.LookupEnv("MONGODB_PASSWORD"); ok {
			serverArgs.MongoDBPassword = v
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		storage, err := mongodb_storage.NewMongoDBStorage(
			serverArgs.MongoDBUsername,
			serverArgs.MongoDBPassword,
			serverArgs.MongoDBHost,
			serverArgs.MongoDBPort,
			serverArgs.MongoDBDatabase,
		)
		if err != nil {
			return fmt.Errorf("storage: %v", err)
		}
		return commissioner.NewServer(
			serverArgs.TokenName,
			storage,
			log,
		).Start(
			serverArgs.Port,
			serverArgs.MetricsPort,
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
	serverCmd.PersistentFlags().IntVar(&serverArgs.MetricsPort, "metrics-port", 0, "optional prometheus metrics port")
	serverCmd.PersistentFlags().IntVar(&serverArgs.NodePort, "node-port", 1337, "cardano-node listener port")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeConfig, "node-config", "", "cardano-node config path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeDatabasePath, "node-database-path", "", "cardano-node database path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeSocketPath, "node-socket-path", "", "cardano-node socket path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeHostAddr, "node-host-addr", "", "cardano-node host address")
	serverCmd.PersistentFlags().StringVar(&serverArgs.NodeTopology, "node-topology", "", "cardano-node topology file path")
	serverCmd.PersistentFlags().StringVar(&serverArgs.TokenName, "token-name", "", "cardano NFT name")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBHost, "mongodb-host", "", "MongoDB service host")
	serverCmd.PersistentFlags().IntVar(&serverArgs.MongoDBPort, "mongodb-port", 27017, "MongoDB service port")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBDatabase, "mongodb-database", "default", "MongoDB database name")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBUsername, "mongodb-username", "", "MongoDB connection username")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBPassword, "mongodb-password", "", "MongoDB connection password")
	ConfigureCommand(serverCmd)
}
