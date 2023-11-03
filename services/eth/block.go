package eth

import (
	"context"
)

func (e *EthClient) GetLatestBlockNumber() (int64, error) {
	// 获取最新的区块头
	header, err := e.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Int64(), nil
}
