package idempotency

import (
	"sync"
	"time"
)

type Store struct {
	responses map[string]time.Time
	mutex     sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		responses: make(map[string]time.Time),
	}
}

func (s *Store) Exists(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// ideally a cleanup would exist
	expiry, exists := s.responses[key]
	if expiry.Before(time.Now()) {
		delete(s.responses, key)
		return false
	}

	return exists
}

func (s *Store) Set(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.responses[key] = time.Now().Add(10 * time.Minute)
}

func (s *Store) Remove(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.responses, key)
}
