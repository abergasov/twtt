package indexator

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
	"twtt/internal/entities"
	"twtt/internal/logger"
)

func (s *Service) ObserveBlocks() {
	s.log.Info("start observe")
	defer s.log.Info("stop observe")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	s.FetchBlocks()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.FetchBlocks()
		}
	}
}

func (s *Service) FetchBlocks() {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	block, err := s.balanceEL.BlockByNumber(ctx, nil)
	if err != nil {
		s.log.Error("unable to get block by number", err)
		return
	}

	currentBlock := atomic.LoadUint64(&s.currentBlock)
	if block.NumberU64() == currentBlock {
		return
	}
	s.log.Info("new block", logger.WithBlockNumber(block.NumberU64()))
	if block.NumberU64()-currentBlock > 1 {
		diff := block.NumberU64() - (currentBlock + 1)
		if currentBlock == 0 {
			diff = block.NumberU64() - 10
		}
		diff = min(diff, 10)

		s.log.Info("load missed blocks", logger.WithUnt64("missed", diff))
		var (
			wg      sync.WaitGroup
			mu      sync.Mutex
			errList []error
			blocks  = make([]*entities.Block, 0, diff)
		)

		// load maximum 10 blocks in case of big gap
		// for production it should be changed to bigger number
		for i := block.NumberU64(); i > block.NumberU64()-diff; i-- {
			wg.Add(1)
			go func(loadBlockID uint64) {
				defer wg.Done()
				s.log.Info("load missed block", logger.WithBlockNumber(loadBlockID))
				b, errL := s.balanceEL.BlockByNumber(ctx, big.NewInt(int64(loadBlockID)))
				mu.Lock()
				defer mu.Unlock()
				if errL != nil {
					errList = append(errList, errL)
					return
				}
				blocks = append(blocks, b)
			}(i)
		}
		wg.Wait()
		if len(errList) > 0 {
			s.log.Error("failed to load missed blocks", fmt.Errorf("%v", errList))
			return
		}
		for i := range blocks {
			s.txSync <- blocks[i].Transactions
		}
	}
	s.txSync <- block.Transactions
	atomic.StoreUint64(&s.currentBlock, block.NumberU64())
}
