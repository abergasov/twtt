package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/shopspring/decimal"
)

const ZeroAddress = "0x0000000000000000000000000000000000000000"

func WeiFromETHString(eth string) *big.Int {
	amount, _ := decimal.NewFromString(eth)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(18))
	return amount.Mul(mul).BigInt()
}

func GWeiFromETHString(eth string) *big.Int {
	amount, _ := decimal.NewFromString(eth)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(9))
	return amount.Mul(mul).BigInt()
}

func ETHFromWei(wei *big.Int) string {
	amount := decimal.NewFromBigInt(wei, 0)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(18))
	return amount.Div(mul).String()
}

func ETHFromGWei(wei *big.Int) string {
	amount := decimal.NewFromBigInt(wei, 0)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(9))
	return amount.Div(mul).String()
}

func CustomFromWei(wei *big.Int, decimals int) string {
	amount := decimal.NewFromBigInt(wei, 0)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	return amount.Div(mul).String()
}

func CustomToWei(amount float64, decimals int) *big.Int {
	amountDecimal := decimal.NewFromFloat(amount)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	return amountDecimal.Mul(mul).BigInt()
}

// ValidateETHAddress validates ETH address. As ETH address is array of bytes, it is not possible to properly validate it
// simply by checking if it is valid hex string. So, we check if it is not empty and not zero address.
func ValidateETHAddress(address string) error {
	if address == "" {
		return fmt.Errorf("address is empty")
	}

	if common.HexToAddress(address).String() == ZeroAddress {
		return fmt.Errorf("address is invalid")
	}
	if !common.IsHexAddress(address) {
		return fmt.Errorf("address is invalid")
	}
	return nil
}
