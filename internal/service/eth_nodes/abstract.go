package ethnodes

import (
	"context"
	"math/big"
	"twtt/internal/entities"

	"github.com/ethereum/go-ethereum/common"
)

type ELNode interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*entities.Block, error)
	BlockByHash(ctx context.Context, blockHash common.Hash) (*entities.Block, error)
}
