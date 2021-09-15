package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/thavlik/bvs/components/commissioner/pkg/commissioner"
	"github.com/thavlik/bvs/components/commissioner/pkg/storage/mongodb_storage"
)

type ServerArgs struct {
	Port              int
	MetricsPort       int
	NodePort          int
	NodeConfig        string
	NodeDatabasePath  string
	NodeSocketPath    string
	NodeHostAddr      string
	NodeTopology      string
	TokenName         string
	MongoDBHost       string
	MongoDBPort       int
	MongoDBDatabase   string
	MongoDBUsername   string
	MongoDBPassword   string
	MongoDBCACertPath string
}

var serverArgs ServerArgs

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if serverArgs.TokenName == "" {
			return fmt.Errorf("missing --token-name")
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
		if v, ok := os.LookupEnv("MONGODB_CACERT"); ok {
			if err := ioutil.WriteFile(
				serverArgs.MongoDBCACertPath,
				[]byte(v),
				0644,
			); err != nil {
				return err
			}
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
			serverArgs.MongoDBCACertPath,
		)
		if err != nil {
			return fmt.Errorf("storage: %v", err)
		}
		return commissioner.NewServer(
			serverArgs.TokenName,
			storage,
		).Start(
			serverArgs.Port,
			serverArgs.MetricsPort,
		)
	},
}

func init() {
	serverCmd.PersistentFlags().IntVar(&serverArgs.Port, "port", 80, "http service listener port")
	serverCmd.PersistentFlags().IntVar(&serverArgs.MetricsPort, "metrics-port", 0, "optional prometheus metrics port")
	serverCmd.PersistentFlags().StringVar(&serverArgs.TokenName, "token-name", "", "cardano NFT name")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBHost, "mongodb-host", "", "MongoDB service host")
	serverCmd.PersistentFlags().IntVar(&serverArgs.MongoDBPort, "mongodb-port", 27017, "MongoDB service port")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBDatabase, "mongodb-database", "default", "MongoDB database name")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBUsername, "mongodb-username", "", "MongoDB connection username")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBPassword, "mongodb-password", "", "MongoDB connection password")
	serverCmd.PersistentFlags().StringVar(&serverArgs.MongoDBCACertPath, "mongodb-cacert-path", "/ca.crt", "MongoDB connection CA certificate file path")
	ConfigureCommand(serverCmd)
}
