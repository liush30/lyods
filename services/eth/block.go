package eth

import (
	"context"
	"math/big"
)

// GetLastBlockNumber 获取上一个区块号
func GetLastBlockNumber(blockNumber *big.Int) *big.Int {
	return new(big.Int).Sub(blockNumber, big.NewInt(1))
}
func (e *EthClient) GetLatestBlockNumber() (int64, error) {
	// 获取最新的区块头
	header, err := e.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Int64(), nil
}
