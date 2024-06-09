package entities

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Block struct {
	Number       *hexutil.Big `json:"number"`
	Hash         string       `json:"hash"`
	number       uint64
	Transactions []*Transaction `json:"transactions"`
}

func (b *Block) NumberU64() uint64 {
	if b.number > 0 {
		return b.number
	}
	b.number = (*big.Int)(b.Number).Uint64()
	return b.number
}
