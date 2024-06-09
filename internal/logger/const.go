package logger

import (
	_ "embed"
	"fmt"
)

func WithTransaction(txHash string) StringWith {
	return StringWith{Key: "_transaction_hash", Val: txHash}
}

func WithBlockNumber(blockNumber uint64) StringWith {
	return StringWith{Key: "_block_id", Val: fmt.Sprint(blockNumber)}
}

func WithService(serviceName string) StringWith {
	return StringWith{Key: "_service", Val: serviceName}
}
