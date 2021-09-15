package storage

import (
	"fmt"
	"time"
)

var ErrKeyNotFound = fmt.Errorf("key not found")

type Election struct {
	ID              string
	SigningKey      string
	VerificationKey string
	Deadline        time.Time
}

type Minter struct {
	ID         string
	SigningKey string
}

type Storage interface {
	StoreElection(*Election) error
	RetrieveElection(id string) (*Election, error)

	StoreMinter(*Minter) error
	RetrieveMinter(id string) (*Minter, error)
}
