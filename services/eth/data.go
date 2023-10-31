package eth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/utils"
	"math/big"
	"strconv"
)

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
