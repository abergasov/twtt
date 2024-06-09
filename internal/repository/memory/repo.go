package memory

import (
	"context"
	"fmt"
	"sync"
	"twtt/internal/entities"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string][]*entities.Transaction
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string][]*entities.Transaction, 1_000),
	}
}

func (s *Storage) Save(_ context.Context, transactions map[string][]*entities.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for addr := range transactions {
		if _, ok := s.data[addr]; !ok {
			s.data[addr] = make([]*entities.Transaction, 0, len(transactions[addr]))
		}
		s.data[addr] = append(s.data[addr], transactions[addr]...)
	}
	return nil
}

func (s *Storage) GetTransactions(_ context.Context, addr string) ([]*entities.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.data[addr]; !ok {
		return nil, fmt.Errorf("no transactions for address %s", addr)
	}
	return s.data[addr], nil
}
