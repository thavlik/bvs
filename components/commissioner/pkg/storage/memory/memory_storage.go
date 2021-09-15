package memory

import (
	"sync"

	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
)

type memoryStorage struct {
	v map[string]string
	l sync.Mutex
}

func NewMemoryStorage() storage.Storage {
	return &memoryStorage{
		v: make(map[string]string),
	}
}

func (s *memoryStorage) Store(key, value string) error {
	s.l.Lock()
	defer s.l.Unlock()
	s.v[key] = value
	return nil
}

func (s *memoryStorage) Retrieve(key string) (string, error) {
	s.l.Lock()
	defer s.l.Unlock()
	v, ok := s.v[key]
	if !ok {
		return "", storage.ErrKeyNotFound
	}
	return v, nil
}
