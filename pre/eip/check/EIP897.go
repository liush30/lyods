package check

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"lyods-adsTool/services/eth"
	"os"
	"strings"
)

func IsEIP897(ethClient *eth.EthClient, contractAddress, contractABIJSON string) (string, error) {
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
			result, err := ethClient.CallContractMethod(contractAddress, contractABIJSON, "implementation")
			if err != nil {
				log.Println("fail call contract method:", err.Error())
				continue
			}
			return result[0].(common.Address).String(), nil
		}
	}

	return "", nil
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
