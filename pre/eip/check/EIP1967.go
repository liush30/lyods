package check

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/services/eth"
	"math/big"
)

// IsEIP1967 判断合约是否是eip1967合约
func IsEIP1967(ethClient *eth.EthClient, contractAddress string) (string, error) {
	proxyAddress, err := GetEIP1967ImplAddress(ethClient, contractAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get EIP1967 proxy implementation address: %v", err)
	}
	// 如果proxyAddress不为零地址，则说明该合约是符合EIP1967规范的合约，并使用的是逻辑合约地址，返回代理委托的逻辑合约地址
	if !isAddressEmpty(proxyAddress) {
		return proxyAddress, nil
	}
	//若未获取到逻辑合约地址，尝试获取信标合约地址
	proxyAddress, err = GetEIP1967BeaconAddress(ethClient, contractAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get EIP1967 beacon proxy address: %v", err)
	}
	// 如果proxyAddress不为零地址，说明该合约是符合EIP1967规范的合约，并使用的是信标合约地址，并通过信标合约地址查询代理合约地址
	if !isAddressEmpty(proxyAddress) {
		return proxyAddress, nil
	}
	return "", nil
}

func isAddressEmpty(address string) bool {
	return address == constants.ZERO_ADDRESS
}

func GetEIP1967ImplAddress(ethClient *eth.EthClient, contractAddress string) (string, error) {
	contractAddressObj := common.HexToAddress(contractAddress)
	keccakHash := crypto.Keccak256([]byte("eip1967.proxy.implementation"))
	bigIntValue := new(big.Int).SetBytes(keccakHash)
	bigIntValue = new(big.Int).Sub(bigIntValue, big.NewInt(1))

	var bytes32Value [32]byte
	copy(bytes32Value[:], bigIntValue.Bytes())
	storageSlot := common.BytesToHash(bytes32Value[:])

	proxyAddress, err := ethClient.StorageAt(context.Background(), contractAddressObj, storageSlot, nil)
	if err != nil {
		return "", fmt.Errorf("GetEIP1967ProxyAddress Error: %v", err)
	}

	return common.HexToAddress(common.BytesToHash(proxyAddress).String()).String(), nil
}

func GetEIP1967BeaconAddress(ethClient *eth.EthClient, contractAddress string) (string, error) {
	contractAddressObj := common.HexToAddress(contractAddress)
	keccakHash := crypto.Keccak256([]byte("eip1967.proxy.beacon"))
	bigIntValue := new(big.Int).SetBytes(keccakHash)
	bigIntValue = new(big.Int).Sub(bigIntValue, big.NewInt(1))

	var bytes32Value [32]byte
	copy(bytes32Value[:], bigIntValue.Bytes())
	storageSlot := common.BytesToHash(bytes32Value[:])

	proxyAddress, err := ethClient.StorageAt(context.Background(), contractAddressObj, storageSlot, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get EIP1967 beacon storage slot: %v", err)
	}

	beaconAddr := common.HexToAddress(common.BytesToHash(proxyAddress).String()).String()

	beaconABI, err := ethClient.GetContractAbiOnEth(beaconAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get beacon ABI: %v", err)
	}

	result, err := ethClient.CallContractMethod(beaconAddr, beaconABI, "implementation")
	if err != nil {
		return "", fmt.Errorf("failed to call beacon contract method: %v", err)
	}

	resultAddress := result[0].(common.Address)
	return resultAddress.String(), nil
}
