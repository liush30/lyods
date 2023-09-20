package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"lyods-adsTool/tool"
	"regexp"
	"strings"
)

// GetChainForAddress 根据地址识别所在链
func GetChainForAddress(address string) string {
	if strings.HasPrefix(address, "0x") && len(address) == 42 {
		return "Ethereum"
	} else if strings.HasPrefix(address, "0x") && len(address) == 43 {
		return "Binance Smart Chain"
	} else if (strings.HasPrefix(address, "X") || strings.HasPrefix(address, "D")) && len(address) == 58 {
		return "Algorand"
	} else if strings.HasPrefix(address, "S") && len(address) == 55 {
		return "Solana"
	} else if strings.HasPrefix(address, "1") || strings.HasPrefix(address, "3") {
		return "Bitcoin" // Bitcoin主网
	} else if strings.HasPrefix(address, "bc1") || strings.HasPrefix(address, "tb1") {
		return "Bitcoin" // Bitcoin SegWit地址
	} else if strings.HasPrefix(address, "L") && len(address) == 34 {
		return "Litecoin"
	} else if strings.HasPrefix(address, "X") && len(address) == 95 {
		return "Monero"
	} else if strings.HasPrefix(address, "0x") && len(address) == 40 {
		return "Ethereum Classic"
	} else if strings.HasPrefix(address, "t1") && len(address) == 34 {
		return "Zcash"
	} else if strings.HasPrefix(address, "X") && len(address) == 35 {
		return "Dash"
	} else {
		return "Unknown Chain"
	}
}

// IsValidAddress 判断地址是否合法
func IsValidAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

// IsContractAddress 判断地址是否为合约地址-以太坊
func IsContractAddress(addressStr string) (bool, error) {
	var address common.Address
	//验证地址是否合法
	if IsValidAddress(addressStr) {
		address = common.HexToAddress(addressStr)
	}
	bytecode, err := tool.EthClient.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal("获取字节码失败:", err)
		return false, err
	}
	//bytecode>0，说明是合约地址
	return len(bytecode) > 0, nil
}
