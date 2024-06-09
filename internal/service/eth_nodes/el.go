package ethnodes

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sync/atomic"
	"time"
	"twtt/internal/entities"
	"twtt/internal/logger"
	"twtt/internal/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type BlockNumber int64

const (
	vsn                  = "2.0"
	SafeBlockNumber      = BlockNumber(-4)
	FinalizedBlockNumber = BlockNumber(-3)
)

// type ELBalancer Balancer[*ethclient.Client]
type ELBalancer Balancer[string]

func NewELBalancer(
	log logger.AppLogger,
	rpcURLs []string,
	rpcFallbackURLs []string,
	retryCount int,
	timeout *time.Duration,
) *ELBalancer {
	balancer := NewBalancer[string](func() (*utils.RoundRobinBalancer[string], *utils.RoundRobinBalancer[string]) {
		return utils.NewRoundRobinBalancer(rpcURLs), utils.NewRoundRobinBalancer(rpcFallbackURLs)
	}, retryCount, retryCount, timeout, log)
	return (*ELBalancer)(balancer)
}

func (b *ELBalancer) BlockByHash(ctx context.Context, blockHash common.Hash) (*entities.Block, error) {
	fn := func(nodes *utils.RoundRobinBalancer[string]) func() (*entities.Block, error) {
		return func() (*entities.Block, error) {
			return b.curlBlock(ctx, nodes.Next(), "eth_getBlockByHash", true, blockHash)
		}
	}
	return Run(ctx, fn, (*Balancer[string])(b), "BlockByNumber")
}

func (b *ELBalancer) BlockByNumber(ctx context.Context, number *big.Int) (*entities.Block, error) {
	fn := func(nodes *utils.RoundRobinBalancer[string]) func() (*entities.Block, error) {
		return func() (*entities.Block, error) {
			return b.curlBlock(ctx, nodes.Next(), "eth_getBlockByNumber", true, toBlockNumArg(number))
		}
	}
	return Run(ctx, fn, (*Balancer[string])(b), "BlockByNumber")
}

func (b *ELBalancer) curlBlock(ctx context.Context, rpcURL, method string, fullTx bool, param any) (*entities.Block, error) {
	type blockWrapper struct {
		Jsonrpc string          `json:"jsonrpc"`
		ID      int             `json:"id"`
		Result  *entities.Block `json:"result"`
	}
	res, code, err := utils.PostCurl[blockWrapper](ctx, rpcURL, map[string]any{
		"method":  method,
		"id":      atomic.AddUint64(&b.requestID, 1),
		"jsonrpc": vsn,
		"params":  []any{param, fullTx},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("failed to get block, wrong code: %d", code)
	}
	return res.Result, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	finalized := big.NewInt(int64(FinalizedBlockNumber))
	if number.Cmp(finalized) == 0 {
		return "finalized"
	}
	safe := big.NewInt(int64(SafeBlockNumber))
	if number.Cmp(safe) == 0 {
		return "safe"
	}
	return hexutil.EncodeBig(number)
}
