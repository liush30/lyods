package check

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"log"
	"strings"
)

func IsERC20(contractABIJSON string) (bool, error) {

	abiJSON, err := readABIFromFile("pre/eip/abi/erc20/erc20.json")
	if err != nil {
		return false, err
	}

	implements := contractImplInput(abiJSON, contractABIJSON)
	if err != nil {
		return false, err
	}

	if implements {
		return true, nil
	}

	return false, nil
}
func contractImplInput(interfaceABIJSON, contractABIJSON string) bool {
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
		if !hasSameTypes(interfaceMethod.Inputs, contractMethod.Inputs) {
			return false
		}
	}

	return true
}
