package indexator

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"twtt/internal/entities"
	"twtt/internal/logger"
	"twtt/internal/repository"
	ethnodes "twtt/internal/service/eth_nodes"
)

type Service struct {
	ctx       context.Context
	log       logger.AppLogger
	balanceEL *ethnodes.ELBalancer
	repo      repository.Storage

	currentBlock uint64
	txSync       chan []*entities.Transaction

	observeAddresses map[string]struct{}
	mu               sync.RWMutex
}

func NewService(
	ctx context.Context,
	log logger.AppLogger,
	balanceEL *ethnodes.ELBalancer,
	repo repository.Storage,
) *Service {
	return &Service{
		ctx:              ctx,
		repo:             repo,
		log:              log.With(logger.WithService("indexator")),
		balanceEL:        balanceEL,
		txSync:           make(chan []*entities.Transaction, 100),
		observeAddresses: make(map[string]struct{}, 1_000),
	}
}

func (s *Service) Run() {
	go s.ObserveBlocks()
	go s.ProcessTransactions()
}

func (s *Service) GetCurrentBlock() uint64 {
	return atomic.LoadUint64(&s.currentBlock)
}

func (s *Service) GetTransactions(address string) ([]*entities.Transaction, error) {
	address = strings.ToLower(address)
	return s.repo.GetTransactions(s.ctx, address)
}

func (s *Service) Subscribe(address string) bool {
	address = strings.ToLower(address)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.observeAddresses[address]; ok {
		return false // already subscribed
	}
	s.observeAddresses[address] = struct{}{}
	return true
}
