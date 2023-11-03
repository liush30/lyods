package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"os"
	"strings"
)

func (e *EthClient) IsProxyContract(contractAddress, contractABIJSON string) (bool, string, error) {
	//判断该地址是否是符合eip1822标准的代理合约
	proxyAddress, err := e.IsZeppelinsUnStructStorage(contractAddress)
	if err != nil {
		return false, "", fmt.Errorf("fail check IsZeppelinsUnStructStorage: %s", err)
	}
	if proxyAddress != "" {
		return true, proxyAddress, nil
	}
	proxyAddress, err = e.IsEIP1967(contractAddress)
	if err != nil {
		return false, "", fmt.Errorf("fail check IsEIP1967: %s", err)
	}
	if proxyAddress != "" {
		return true, proxyAddress, nil
	}
	proxyAddress, err = e.IsEIP1822(contractAddress)
	if err != nil {
		return false, "", fmt.Errorf("fail check IsEIP1822: %s", err)
	}
	if proxyAddress != "" {
		return true, proxyAddress, nil
	}
	proxyAddress, err = e.IsEIP897(contractAddress, contractABIJSON)
	if err != nil {
		return false, "", fmt.Errorf("fail check IsEIP897: %s", err)
	}
	if proxyAddress != "" {
		return true, proxyAddress, nil
	}
	return false, "", nil
}
func (e *EthClient) IsEIP897(contractAddress, contractABIJSON string) (string, error) {
	abiFiles := []string{
		"AppProxyUpgradeableAbi.json",
		"AppProxyPinned.json",
		"KernelProxy.json",
		"eip897.json",
		"zepplinOSAbi.json",
	}

	for _, fileName := range abiFiles {
		abiJSON, err := readABIFromFile("pre/eip/abi/eip-897/" + fileName)
		if err != nil {
			return "", err
		}

		implements := contractImplementsInterface(abiJSON, contractABIJSON)
		if err != nil {
			return "", err
		}

		if implements {
			result, err := e.CallContractMethod(contractAddress, contractABIJSON, "implementation")
			if err != nil {
				log.Println("fail call contract method:", err.Error())
				continue
			}
			return result[0].(common.Address).String(), nil
		}
	}

	return "", nil
}

// IsEIP1822 判断合约是否是eip1822标准代理合约
func (e *EthClient) IsEIP1822(contractAddress string) (string, error) {
	implAddress, err := getStorageValue(e, contractAddress, "PROXIABLE")
	if err != nil {
		return "", err
	}
	return implAddress, nil
}

// IsEIP1967 判断合约是否是eip1967合约
func (e *EthClient) IsEIP1967(contractAddress string) (string, error) {
	proxyAddress, err := e.GetEIP1967ImplAddress(contractAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get EIP1967 proxy implementation address: %v", err)
	}
	// 如果proxyAddress不为零地址，则说明该合约是符合EIP1967规范的合约，并使用的是逻辑合约地址，返回代理委托的逻辑合约地址
	if !IsAddressEmpty(proxyAddress) {
		return proxyAddress, nil
	}
	//若未获取到逻辑合约地址，尝试获取信标合约地址
	proxyAddress, err = e.GetEIP1967BeaconAddress(contractAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get EIP1967 beacon proxy address: %v", err)
	}
	// 如果proxyAddress不为零地址，说明该合约是符合EIP1967规范的合约，并使用的是信标合约地址，并通过信标合约地址查询代理合约地址
	if !IsAddressEmpty(proxyAddress) {
		return proxyAddress, nil
	}
	return "", nil
}

func (e *EthClient) GetEIP1967ImplAddress(contractAddress string) (string, error) {
	contractAddressObj := common.HexToAddress(contractAddress)
	keccakHash := crypto.Keccak256([]byte("eip1967.proxy.implementation"))
	bigIntValue := new(big.Int).SetBytes(keccakHash)
	bigIntValue = new(big.Int).Sub(bigIntValue, big.NewInt(1))

	var bytes32Value [32]byte
	copy(bytes32Value[:], bigIntValue.Bytes())
	storageSlot := common.BytesToHash(bytes32Value[:])

	proxyAddress, err := e.StorageAt(context.Background(), contractAddressObj, storageSlot, nil)
	if err != nil {
		return "", fmt.Errorf("GetEIP1967ProxyAddress Error: %v", err)
	}

	return common.HexToAddress(common.BytesToHash(proxyAddress).String()).String(), nil
}

func (e *EthClient) GetEIP1967BeaconAddress(contractAddress string) (string, error) {
	contractAddressObj := common.HexToAddress(contractAddress)
	keccakHash := crypto.Keccak256([]byte("eip1967.proxy.beacon"))
	bigIntValue := new(big.Int).SetBytes(keccakHash)
	bigIntValue = new(big.Int).Sub(bigIntValue, big.NewInt(1))

	var bytes32Value [32]byte
	copy(bytes32Value[:], bigIntValue.Bytes())
	storageSlot := common.BytesToHash(bytes32Value[:])

	proxyAddress, err := e.StorageAt(context.Background(), contractAddressObj, storageSlot, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get EIP1967 beacon storage slot: %v", err)
	}

	beaconAddr := common.HexToAddress(common.BytesToHash(proxyAddress).String()).String()

	beaconABI, err := e.GetContractAbiOnEth(beaconAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get beacon ABI: %v", err)
	}

	result, err := e.CallContractMethod(beaconAddr, beaconABI, "implementation")
	if err != nil {
		return "", fmt.Errorf("failed to call beacon contract method: %v", err)
	}

	resultAddress := result[0].(common.Address)
	return resultAddress.String(), nil
}
func IsERC20(contractABIJSON string) (bool, error) {

	abiJSON, err := readABIFromFile("pre/eip/abi/erc20/erc20.json")
	if err != nil {
		return false, err
	}

	implements, err := contractImplInput(abiJSON, contractABIJSON)
	if err != nil {
		return false, err
	}

	if implements {
		return true, nil
	}

	return false, nil
}
func contractImplInput(interfaceABIJSON, contractABIJSON string) (bool, error) {
	interfaceABI, err := abi.JSON(strings.NewReader(interfaceABIJSON))
	if err != nil {
		return false, fmt.Errorf("fail parse interface abi: %s", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		return false, fmt.Errorf("fail parse contract abi: %s", err)
	}

	// 创建一个映射以便快速查找合约ABI中的方法
	contractMethods := make(map[string]abi.Method, len(contractABI.Methods))
	for _, method := range contractABI.Methods {
		contractMethods[method.Name] = method
	}

	for _, interfaceMethod := range interfaceABI.Methods {
		contractMethod, found := contractMethods[interfaceMethod.Name]

		if !found {
			return false, nil
		}

		// 检查函数的输入和输出类型是否匹配
		if !hasSameTypes(interfaceMethod.Inputs, contractMethod.Inputs) {
			return false, nil
		}
	}

	return true, nil
}

// IsZeppelinsUnStructStorage 判断合约是否符合openzeppelin 非结构化存储代理模式
func (e *EthClient) IsZeppelinsUnStructStorage(contractAddress string) (string, error) {
	implAddress, err := getStorageValue(e, contractAddress, "org.zeppelinos.proxy.implementation")
	if err != nil {
		return "", err
	}
	return implAddress, nil
}
func getStorageValue(ethClient *EthClient, contractAddress string, key string) (string, error) {
	contractAddressObj := common.HexToAddress(contractAddress)
	keccakHash := crypto.Keccak256([]byte(key))
	storageSlot := common.BytesToHash(keccakHash)
	by, err := ethClient.StorageAt(context.Background(), contractAddressObj, storageSlot, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get storage value: %v", err)
	}
	implAddress := common.HexToAddress(common.BytesToHash(by).String()).String()
	if IsAddressEmpty(implAddress) {
		return "", nil
	}
	return common.HexToAddress(common.BytesToHash(by).String()).String(), nil
}

func contractImplementsInterface(interfaceABIJSON, contractABIJSON string) bool {
	interfaceABI, err := abi.JSON(strings.NewReader(interfaceABIJSON))
	if err != nil {
		log.Fatal(err)
	}

	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个映射以便快速查找合约ABI中的方法
	contractMethods := make(map[string]abi.Method, len(contractABI.Methods))
	for _, method := range contractABI.Methods {
		contractMethods[method.Name] = method
	}

	for _, interfaceMethod := range interfaceABI.Methods {
		contractMethod, found := contractMethods[interfaceMethod.Name]

		if !found {
			return false
		}

		// 检查函数的输入和输出类型是否匹配
		if !hasSameTypes(interfaceMethod.Inputs, contractMethod.Inputs) || !hasSameTypes(interfaceMethod.Outputs, contractMethod.Outputs) {
			return false
		}
	}

	return true
}

func hasSameTypes(a, b abi.Arguments) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Type.String() != b[i].Type.String() {
			return false
		}
	}

	return true
}

func readABIFromFile(filename string) (string, error) {
	// 读取文件内容
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 将文件内容转换为字符串
	abiData := string(file)
	return abiData, nil
}
