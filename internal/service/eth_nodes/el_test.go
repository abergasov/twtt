package ethnodes_test

import (
	"math/big"
	"testing"
	"time"
	"twtt/internal/entities"
	ethnodes "twtt/internal/service/eth_nodes"
	testhelpers "twtt/internal/test_helpers"
	"twtt/internal/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestELBalancer_TransactionByHash(t *testing.T) {
	// given
	testhelpers.InitMainetRPC(t)
	container := testhelpers.GetClean(t)
	clientBalancer := ethnodes.NewELBalancer(
		container.Logger,
		testhelpers.MainetRPC.Values(),
		testhelpers.MainetRPC.Values(),
		3,
		utils.ToPointer(100*time.Millisecond),
	)

	t.Run("should return transaction by hash", func(t *testing.T) {
		// when
		block, errB := clientBalancer.BlockByHash(container.Ctx, common.HexToHash("0xc8edea36c004659d96827d066341a7e215a211b0607050be745eee283b5d9397"))
		require.NoError(t, errB)

		// then
		checkBlock(t, block)
	})
	t.Run("should return block by number", func(t *testing.T) {
		// when
		block, errB := clientBalancer.BlockByNumber(container.Ctx, big.NewInt(1696337))
		require.NoError(t, errB)

		// then
		checkBlock(t, block)
	})
}

func checkBlock(t *testing.T, block *entities.Block) {
	require.NotNil(t, block)
	require.Equal(t, uint64(1696337), block.NumberU64())
	require.Equal(t, "0xc8edea36c004659d96827d066341a7e215a211b0607050be745eee283b5d9397", block.Hash)
	require.Len(t, block.Transactions, 54)
}
