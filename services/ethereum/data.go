package ethereum

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"
)

// convertDataItem 据数据项的类型执行不同的转换
func convertDataItem(dataItem interface{}, dataType string) (string, bool) {
	switch dataType {
	case "address":
		if addressV, addressOk := dataItem.(common.Address); addressOk {
			return addressV.Hex(), true
		}
	case "bytes32":
		if bytes32V, ok := dataItem.([32]byte); ok {
			return common.BytesToHash(bytes32V[:]).Hex(), true
		}
	case "uint256", "int256":
		if uint256V, uint256Ok := dataItem.(*big.Int); uint256Ok {
			return uint256V.String(), true
		}
	case "bool":
		if boolValue, boolOk := dataItem.(bool); boolOk {
			return strconv.FormatBool(boolValue), true
		}
	case "bytes[]":
		if bytesValue, bytesOK := dataItem.([]byte); bytesOK {
			return hex.EncodeToString(bytesValue), true
		}
	case "string":
		if stringValue, stringOk := dataItem.(string); stringOk {
			return stringValue, true
		}
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64":
		if intValue, intOk := dataItem.(int64); intOk {
			return strconv.FormatInt(intValue, 10), true
		}
	}
	return "", false
}
