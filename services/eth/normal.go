package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"lyods-adsTool/config"
	"lyods-adsTool/pkg/constants"
	"math/big"
)

// InternalTxnParam InternalTxn解析参数
type InternalTxnParam struct {
	eventName                 string                       //事件名称
	contractAddress           string                       //合约地址
	eventNameToValueByAddress map[string]map[string]string //事件相关参数信息
	blockNumber               *big.Int                     //区块高度
	isErc20                   bool                         //是否为erc20
	length                    int                          //参数个数
	token                     string                       //token name
	tokenDecimal              int                          //token decimal
	//internalTxn *domain.InternalTxn,
}

// getNormalTransUrl 返回etherScan中查询指定地址的普通交易信息列表Url
func getNormalTransUrl(addr string) string {
	return constants.API_ETH_TRANS + addr
}

func getTraceTransactionUrl() string {
	return constants.URL_CHAINBASE + config.CHAINBASE_KEY
}

// IsContractAddress 判断地址是否为合约地址-以太坊
func (e *EthClient) IsContractAddress(addressStr string) (bool, error) {
	var address common.Address
	//获取字节码信息
	bytecode, err := e.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Println("Fail get byte code:", err)
		return false, err
	}
	//bytecode>0，说明是合约地址
	return len(bytecode) > 0, nil
}
