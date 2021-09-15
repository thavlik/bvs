package memory

import (
	"sync"

	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
)

type memoryStorage struct {
	v map[string]*storage.Election
	l sync.Mutex
}

func NewMemoryStorage() storage.Storage {
	return &memoryStorage{
		v: make(map[string]*storage.Election),
	}
}

func (s *memoryStorage) StoreElection(e *storage.Election) error {
	s.l.Lock()
	defer s.l.Unlock()
	s.v[e.ID] = e
	return nil
}

func (s *memoryStorage) RetrieveElection(key string) (*storage.Election, error) {
	s.l.Lock()
	defer s.l.Unlock()
	v, ok := s.v[key]
	if !ok {
		return nil, storage.ErrKeyNotFound
	}
	return v, nil
}
