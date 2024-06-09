package entities

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Transaction struct {
	BlockHash   common.Hash  `json:"blockHash"`
	BlockNumber *hexutil.Big `json:"blockNumber"`
	From        string       `json:"from"`
	To          string       `json:"to"`
	Gas         string       `json:"gas"`
	GasPrice    string       `json:"gas_price"`
	Nonce       string       `json:"nonce"`
	Value       string       `json:"value"`
	Data        string       `json:"input"`
}
