package check

import (
	"lyods-adsTool/services/eth"
)

// IsZeppelinsUnStructStorage 判断合约是否符合openzeppelin 非结构化存储代理模式
func IsZeppelinsUnStructStorage(ethClient *eth.EthClient, contractAddress string) (string, error) {
	implAddress, err := getStorageValue(ethClient, contractAddress, "org.zeppelinos.proxy.implementation")
	if err != nil {
		return "", err
	}
	return implAddress, nil
}
