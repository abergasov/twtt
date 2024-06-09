package utils_test

import (
	"github.com/google/uuid"
	"math/big"
	"strings"
	"testing"
	"twtt/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestWeiFromETHString(t *testing.T) {
	table := map[string]string{
		"0.02": "20000000000000000",
		"0.1":  "100000000000000000",
	}
	for eth, wei := range table {
		res := utils.WeiFromETHString(eth).String()
		require.Equal(t, wei, res)
	}
}

func TestGWeiFromETHString(t *testing.T) {
	table := map[string]string{
		"0.02": "20000000",
		"0.1":  "100000000",
	}
	for eth, gwei := range table {
		res := utils.GWeiFromETHString(eth).String()
		require.Equal(t, gwei, res)
	}
}

func TestETHFromWei(t *testing.T) {
	a, ok := big.NewInt(0).SetString("20000000000000000", 10)
	require.True(t, ok)
	b, ok := big.NewInt(0).SetString("100000000000000000", 10)
	require.True(t, ok)
	table := map[*big.Int]string{
		a: "0.02",
		b: "0.1",
	}
	for wei, eth := range table {
		res := utils.ETHFromWei(wei)
		require.Equal(t, eth, res)
	}
}

func TestETHFromGWei(t *testing.T) {
	a, ok := big.NewInt(0).SetString("20000000", 10)
	require.True(t, ok)
	b, ok := big.NewInt(0).SetString("100000000", 10)
	require.True(t, ok)
	table := map[*big.Int]string{
		a: "0.02",
		b: "0.1",
	}
	for wei, eth := range table {
		res := utils.ETHFromGWei(wei)
		require.Equal(t, eth, res)
	}
}

func TestCustomFromWei(t *testing.T) {
	table := map[string]string{
		"2815394107": "2815.394107",
	}
	for wei, eth := range table {
		a, ok := big.NewInt(0).SetString(wei, 10)
		require.True(t, ok)
		res := utils.CustomFromWei(a, 6)
		require.Equal(t, eth, res)
	}
}

func TestCustomToWei(t *testing.T) {
	table := map[float64]string{
		2815.394107: "2815394107",
	}
	for val, wei := range table {
		a, ok := big.NewInt(0).SetString(wei, 10)
		require.True(t, ok)
		res := utils.CustomToWei(val, 6)
		require.Equal(t, a, res)
	}
}

func TestValidateETHAddress(t *testing.T) {
	validAddresses := []string{
		"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		"0xFEa7a6a0B346362BF88A9e4A88416B77a57D6c2A",
		"0x130966628846BFd36ff31a822705796e8cb8C18D",
		"0x82f0B8B456c1A451378467398982d4834b6829c1",
		"0xfE19F0B51438fd612f6FD59C1dbB3eA319f433Ba",
		"0xAf5191B0De278C7286d6C7CC6ab6BB8A73bA2Cd6",
		"0x6694340fc020c5E6B96567843da2df01b2CE1eb6",
		"0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590",
		"0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590",
		"0x296F55F8Fb28E498B858d0BcDA06D955B2Cb3f97",
		"0xB0D502E938ed5f4df2E681fE6E419ff29631d62b",
		"0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590",
		"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		"0xFA5Ed56A203466CbBC2430a43c66b9D8723528E7",
		"0xE0B52e49357Fd4DAf2c15e02058DCE6BC0057db4",
		"0xAEC8318a9a59bAEb39861d10ff6C7f7bf1F96C57",
		"0x4b1E2c2762667331Bc91648052F646d1b0d35984",
		"0xC16B81Af351BA9e64C1a069E3Ab18c244A1E3049",
		"0xFA5Ed56A203466CbBC2430a43c66b9D8723528E7",
		"0xf1dDcACA7D17f8030Ab2eb54f2D9811365EFe123",
		"0x840b25c87B626a259CA5AC32124fA752F0230a72",
		"0x14C00080F97B9069ae3B4Eb506ee8a633f8F5434",
		"0x0c1EBBb61374dA1a8C57cB6681bF27178360d36F",
		"0x2297aebd383787a160dd0d9f71508148769342e3",
		"0x2297aEbD383787A160DD0d9F71508148769342E3",
		"0x2297aEbD383787A160DD0d9F71508148769342E3",
		"0x152b9d0FdC40C096757F570A51E494bd4b943E50",
		"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		"0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		"0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8",
		"0x7F5c764cBc14f9669B88837ca1490cCa17c31607",
		"0x04068DA6C83AFCFA0e13ba15A6696662335D5B75",
		"0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E",
		"0xa4151B2B3e269645181dCcF2D426cE75fcbDeca9",
		"0x9b5fae311A4A4b9d838f301C9c27b55d19BAa4Fb",
		"0x6B175474E89094C44Da98b954EedeAC495271d0F",
		"0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
		"0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1",
		"0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1",
		"0x8D11eC38a3EB5E956B052f67Da8Bdc9bef8Abf3E",
		"0xd586E7F844cEa2F87f50152665BCbc2C279D8d70",
		"0xdAC17F958D2ee523a2206206994597C13D831ec7",
		"0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
		"0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
		"0x94b008aA00579c1307B0EF2c499aD98a8ce58e58",
		"0x1B27A9dE6a775F98aaA5B90B62a4e2A0B84DbDd9",
		"0x9702230A8Ea53601f5cD2dc00fDBc13d4dF4A8c7",
	}
	for _, addr := range validAddresses {
		require.NoError(t, utils.ValidateETHAddress(addr))
		require.NoError(t, utils.ValidateETHAddress(strings.ToLower(addr)))
	}
	require.Error(t, utils.ValidateETHAddress(""))
	require.Error(t, utils.ValidateETHAddress(utils.ZeroAddress))
	require.Error(t, utils.ValidateETHAddress(uuid.NewString()))
	require.Error(t, utils.ValidateETHAddress("abc"))
}
