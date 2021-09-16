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
	client    *mongo.Client
	db        *mongo.Database
	elections *mongo.Collection
	minters   *mongo.Collection
}

type mongoElection struct {
	ID               string `bson:"_id"`
	SigningKey       string `bson:"signingKey"`
	VerificationKey  string `bson:"verificationKey"`
	PolicyID         string `bson:"policyId"`
	MintingScript    string `bson:"mintingScript"`
	InvalidHereafter int    `bson:"invalidHereafter"`
	Protocol         string `bson:"protocol"`
}

type mongoMinter struct {
	ID              string `bson:"_id"`
	SigningKey      string `bson:"signingKey"`
	VerificationKey string `bson:"verificationKey"`
}

func NewMongoDBStorage(
	username,
	password,
	host string,
	port int,
	database,
	cacertPath string,
) (storage.Storage, error) {
	if port != 27017 {
		return nil, fmt.Errorf("by virtue of a limitation in mongodb, only port 27017 is supported")
	}
	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/%s?tls=true&tlsCAFile=%s",
		username,
		password,
		host,
		database,
		cacertPath,
	)
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}
	db := client.Database(database)
	return &mongoStorage{
		client,
		db,
		db.Collection("elections"),
		db.Collection("minters"),
	}, nil
}

func (s *mongoStorage) StoreElection(v *storage.Election) error {
	if _, err := s.elections.InsertOne(context.TODO(), &mongoElection{
		ID:               v.ID,
		SigningKey:       v.SigningKey,
		VerificationKey:  v.VerificationKey,
		PolicyID:         v.PolicyID,
		MintingScript:    v.MintingScript,
		InvalidHereafter: v.InvalidHereafter,
		Protocol:         v.Protocol,
	}); err != nil {
		return fmt.Errorf("mongo insert: %v", err)
	}
	return nil
}

func (s *mongoStorage) RetrieveElection(id string) (*storage.Election, error) {
	result := s.elections.FindOne(context.TODO(), map[string]interface{}{
		"_id": id,
	})
	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("mongo find: %v", err)
	}
	v := &mongoElection{}
	if err := result.Decode(v); err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}
	return &storage.Election{
		ID:               v.ID,
		SigningKey:       v.SigningKey,
		VerificationKey:  v.VerificationKey,
		PolicyID:         v.PolicyID,
		MintingScript:    v.MintingScript,
		InvalidHereafter: v.InvalidHereafter,
		Protocol:         v.Protocol,
	}, nil
}

func (s *mongoStorage) StoreMinter(v *storage.Minter) error {
	if _, err := s.minters.InsertOne(context.TODO(), &mongoMinter{
		ID:              v.ID,
		SigningKey:      v.SigningKey,
		VerificationKey: v.VerificationKey,
	}); err != nil {
		return fmt.Errorf("mongo insert: %v", err)
	}
	return nil
}

func (s *mongoStorage) RetrieveMinter(id string) (*storage.Minter, error) {
	result := s.minters.FindOne(context.TODO(), map[string]interface{}{
		"_id": id,
	})
	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("mongo find: %v", err)
	}
	v := &mongoMinter{}
	if err := result.Decode(v); err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}
	return &storage.Minter{
		ID:              v.ID,
		SigningKey:      v.SigningKey,
		VerificationKey: v.VerificationKey,
	}, nil
}
