package storage

import (
	"fmt"
)

var ErrKeyNotFound = fmt.Errorf("key not found")

type Election struct {
	ID               string
	SigningKey       string
	VerificationKey  string
	InvalidHereafter int
	MintingScript    string
	PolicyID         string
	Protocol         string
}

type Minter struct {
	ID              string
	SigningKey      string
	VerificationKey string
	Address         string
}

type Storage interface {
	StoreElection(*Election) error
	RetrieveElection(id string) (*Election, error)

	StoreMinter(*Minter) error
	RetrieveMinter(id string) (*Minter, error)
}
