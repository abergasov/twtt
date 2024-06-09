package repository

import (
	"context"
	"twtt/internal/entities"
)

type Storage interface {
	Save(ctx context.Context, transactions map[string][]*entities.Transaction) error
	GetTransactions(ctx context.Context, addr string) ([]*entities.Transaction, error)
}
