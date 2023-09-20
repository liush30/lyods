package tool

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"lyods-adsTool/pkg/constants"
)

var EthClient *ethclient.Client

func init() {
	EthClient = createEthClient()
}

// 创建ethereum客户端
func createEthClient() *ethclient.Client {
	//连接infura网关
	client, err := ethclient.Dial(constants.URL_INFRUA)
	if err != nil {
		log.Fatal("连接失败:", err)
		return nil
	}
	return client
}

// 获取地址发出的交易数量（out方向）
func getAddrOutTransCount(addStr string) {
	address := common.HexToAddress(addStr)
	EthClient.NonceAt(context.Background(), address, nil)
}

// 根据合约地址查询风险名单信息
func getListByContractAddr(address string) {
	//根据合约地址查询合约所有交易信息
}
