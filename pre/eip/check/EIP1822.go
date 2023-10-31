package check

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"lyods-adsTool/services/eth"
)

// IsEIP1822 判断合约是否是eip1822标准代理合约
func IsEIP1822(ethClient *eth.EthClient, contractAddress string) (string, error) {
	implAddress, err := getStorageValue(ethClient, contractAddress, "PROXIABLE")
	if err != nil {
		return "", err
	}
	return implAddress, nil
}
func getStorageValue(ethClient *eth.EthClient, contractAddress string, key string) (string, error) {
	contractAddressObj := common.HexToAddress(contractAddress)
	keccakHash := crypto.Keccak256([]byte(key))
	storageSlot := common.BytesToHash(keccakHash)
	by, err := ethClient.StorageAt(context.Background(), contractAddressObj, storageSlot, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get storage value: %v", err)
	}
	implAddress := common.HexToAddress(common.BytesToHash(by).String()).String()
	if isAddressEmpty(implAddress) {
		return "", nil
	}
	return common.HexToAddress(common.BytesToHash(by).String()).String(), nil
}
