package mongodb_storage

import (
	"context"
	"fmt"

	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoStorage struct {
	client *mongo.Client
}

func NewMongoDBStorage(
	username,
	password,
	host string,
	port int,
	database string,
) (storage.Storage, error) {
	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s:%d/%s?w=majority",
		username,
		password,
		host,
		port,
		database,
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}
	return &mongoStorage{client}, nil
}

func (s *mongoStorage) Store(key, value string) error {
	return nil
}

func (s *mongoStorage) Retrieve(key string) (string, error) {
	return "", nil
}
