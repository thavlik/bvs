package memory

import (
	"sync"

	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
)

type memoryStorage struct {
	elections map[string]*storage.Election
	minters   map[string]*storage.Minter
	l         sync.Mutex
}

func NewMemoryStorage() storage.Storage {
	return &memoryStorage{
		elections: make(map[string]*storage.Election),
		minters:   make(map[string]*storage.Minter),
	}
}

func (s *memoryStorage) StoreElection(e *storage.Election) error {
	s.l.Lock()
	defer s.l.Unlock()
	s.elections[e.ID] = e
	return nil
}

func (s *memoryStorage) RetrieveElection(key string) (*storage.Election, error) {
	s.l.Lock()
	defer s.l.Unlock()
	v, ok := s.elections[key]
	if !ok {
		return nil, storage.ErrKeyNotFound
	}
	return v, nil
}

func (s *memoryStorage) StoreMinter(v *storage.Minter) error {
	s.l.Lock()
	defer s.l.Unlock()
	s.minters[v.ID] = v
	return nil
}

func (s *memoryStorage) RetrieveMinter(key string) (*storage.Minter, error) {
	s.l.Lock()
	defer s.l.Unlock()
	v, ok := s.minters[key]
	if !ok {
		return nil, storage.ErrKeyNotFound
	}
	return v, nil
}
