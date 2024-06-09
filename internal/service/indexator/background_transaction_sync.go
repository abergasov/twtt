package indexator

import (
	"strings"
	"twtt/internal/entities"
)

func (s *Service) ProcessTransactions() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case txList := <-s.txSync:
			s.sortTransactions(txList)
		}
	}
}

func (s *Service) sortTransactions(txList []*entities.Transaction) {
	saveData := make(map[string][]*entities.Transaction, len(txList))
	s.mu.Lock()
	for i := range txList {
		from := strings.ToLower(txList[i].From)
		to := strings.ToLower(txList[i].To)
		_, validFrom := s.observeAddresses[from]
		_, validTo := s.observeAddresses[to]
		if !(validFrom || validTo) {
			continue
		}
		if validFrom {
			if _, ok := saveData[from]; !ok {
				saveData[from] = make([]*entities.Transaction, 0, len(txList))
			}
			saveData[from] = append(saveData[from], txList[i])
		}
		if validTo {
			if _, ok := saveData[to]; !ok {
				saveData[to] = make([]*entities.Transaction, 0, len(txList))
			}
			saveData[to] = append(saveData[to], txList[i])
		}
	}
	s.mu.Unlock()
	if len(saveData) == 0 {
		return
	}
	if err := s.repo.Save(s.ctx, saveData); err != nil {
		s.log.Error("failed to save transactions", err)
	}
}
