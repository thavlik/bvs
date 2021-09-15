package storage

import "fmt"

var ErrKeyNotFound = fmt.Errorf("key not found")

type Storage interface {
	Store(key, value string) error
	Retrieve(key string) (string, error)
}
