package eth

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"lyods-adsTool/config"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"math/big"
	"strconv"
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
func WeiToEth(wei *big.Int) *big.Float {
	// 创建一个 big.Float 类型的值，用于表示 ETH
	eth := new(big.Float)

	// 创建一个 big.Int 类型的值，表示 10^18（1 ETH 对应的 wei 数量）
	weiPerEth := new(big.Int)
	weiPerEth.Exp(big.NewInt(10), big.NewInt(18), nil)

	// 将 wei 转换为 ETH
	eth.SetPrec(256) // 设置精度，可根据需要调整
	eth.SetInt(wei)
	eth.Quo(eth, new(big.Float).SetInt(weiPerEth))

	return eth
}
func ConvertTokenValue(value *big.Int, decimal int) float64 {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	valueInDecimal := new(big.Float).SetInt(value)
	scaleInDecimal := new(big.Float).SetInt(scale)
	result := new(big.Float).Quo(valueInDecimal, scaleInDecimal)
	f, _ := result.Float64()
	return f
}
func IsAddressEmpty(address string) bool {
	return address == constants.ZERO_ADDRESS
}

// convertDataItem 据数据项的类型执行不同的转换
func interfaceToData(dataItem interface{}, dataType string) (string, bool) {
	var dataResult string
	var dataOk bool
	switch dataType {
	case "address":
		if addressV, addressOk := dataItem.(common.Address); addressOk {
			dataResult, dataOk = addressV.Hex(), true
		}
	case "bytes32":
		if bytes32V, ok := dataItem.([32]byte); ok {
			dataResult, dataOk = common.BytesToHash(bytes32V[:]).Hex(), true
		}
	case "uint256", "int256":
		if uint256V, uint256Ok := dataItem.(*big.Int); uint256Ok {
			dataResult, dataOk = uint256V.String(), true
		}
	case "bool":
		if boolValue, boolOk := dataItem.(bool); boolOk {
			dataResult, dataOk = strconv.FormatBool(boolValue), true
		}
	case "bytes[]":
		if bytesValue, bytesOK := dataItem.([]byte); bytesOK {
			dataResult, dataOk = hex.EncodeToString(bytesValue), true
		}
	case "string":
		if stringValue, stringOk := dataItem.(string); stringOk {
			dataResult, dataOk = stringValue, true
		}
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64":
		if intValue, intOk := dataItem.(int64); intOk {
			dataResult, dataOk = strconv.FormatInt(intValue, 10), true
		}
	}
	return dataResult, dataOk
}
func hexToData(inputType string, hash common.Hash) string {
	var dataResult string
	switch inputType {
	case "address":
		dataResult = utils.ConvertAddressEth(hash)
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64":
		dataResult = hash.Big().String()
	default:
		dataResult = hash.String()
	}
	return dataResult
}

// 解析data数据
func parseData(paramName, paramType []string, eventName string, dataLog []byte, contractABI abi.ABI, dataResult *map[string]string) error {

	// 检查参数的有效性
	if len(paramName) != len(paramType) {
		return errors.New("invalid input data: paramName and paramType must have the same length")
	}

	// 解析data数据
	dataInter, err := contractABI.Unpack(eventName, dataLog)
	if err != nil {
		return fmt.Errorf("failed to unpack data: %s", err.Error())
	}
	if len(paramName) != len(dataInter) {
		return errors.New("invalid input data: paramName and dataInter must have the same length")
	}

	// 解析data数据，将参数名与值以键值对的形式存储到paramInfo中
	for i, dataItem := range dataInter {
		dataType := paramType[i]
		dataName := paramName[i]

		// 转换数据项
		convertedValue, ok := interfaceToData(dataItem, dataType)
		if !ok {
			// 转换失败
			return fmt.Errorf("failed to convert %s type. Parameter name is %s", dataType, dataName)
		}
		(*dataResult)[dataName] = convertedValue
	}

	return nil
}

// j将解析后的键值对信息转存为TopicsValStruct
func mapToTopicsValStruct(topics map[string]string) []domain.TopicsValStruct {
	var resultList []domain.TopicsValStruct
	for key, value := range topics {
		topicsValStruct := domain.TopicsValStruct{
			Key:   key,
			Value: value,
		}
		resultList = append(resultList, topicsValStruct)
	}
	return resultList
}

// GetLastBlockNumber 获取上一个区块号
func GetLastBlockNumber(blockNumber *big.Int) *big.Int {
	return new(big.Int).Sub(blockNumber, big.NewInt(1))
}