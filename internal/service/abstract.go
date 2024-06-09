package service

import "twtt/internal/entities"

type Parser interface {
	// GetCurrentBlock last parsed block
	GetCurrentBlock() uint64
	// Subscribe add address to observer
	Subscribe(address string) bool
	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) ([]entities.Transaction, error)
}
