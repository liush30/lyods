package ethereum

import "math/big"

// GetLastBlockNumber 获取上一个区块号
func GetLastBlockNumber(blockNumber *big.Int) *big.Int {
	return new(big.Int).Sub(blockNumber, big.NewInt(1))
}
