package pre_service

import (
	"database/sql"
	"fmt"
	"log"
	"lyods-adsTool/db"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/services/evm"
)

func processContractAndStore(dbClient *sql.DB, e *evm.EthClient, addr, key, ABIStr string) (bool, error) {
	proxyAddress := ""
	var err error
	switch key {
	case "eip1967":
		proxyAddress, err = e.IsEIP1967(addr)
	case "eip1822":
		proxyAddress, err = e.IsEIP1822(addr)
	case "eip897":
		proxyAddress, err = e.IsEIP897(addr, ABIStr)
	case "OpenZeppelin's Unstructured":
		proxyAddress, err = e.IsZeppelinsUnStructStorage(addr)
	}

	if err != nil {
		return false, fmt.Errorf("%s fail check %s: %s\n", addr, key, err.Error())
	}

	if proxyAddress != "" {
		proxyABI, err := e.GetContractAbiOnEth(proxyAddress)
		if err != nil {
			return false, fmt.Errorf("%s fail get proxy contract(%s) address abi: %s\n", addr, key, err.Error())
		}

		err = db.SaveContractABI(dbClient, constants.DB_CHAIN_ETH, addr, proxyABI, proxyAddress)
		if err != nil {
			return false, fmt.Errorf("%s fail save contract abi: %s\n", addr, err.Error())
		}
		return true, nil
	}
	return false, nil
}
func processErc20(dbClient *sql.DB, addr, ABIStr string) (bool, error) {
	//验证合约abi是否符合erc20规范
	isErc20, err := evm.IsERC20(ABIStr)
	if err != nil {
		return false, fmt.Errorf("%s fail check erc20:%v", addr, err)
	}
	//如果是erc20合约，直接将abi存储到数据库
	if isErc20 {
		err = db.SaveContractABI(dbClient, constants.DB_CHAIN_ETH, addr, ABIStr, "")
		if err != nil {
			return false, fmt.Errorf("%s fail save contract abi:%v", addr, err)
		}
	}
	return isErc20, nil
}

// GetABIToDbOnEth 查询数据库中的erc20信息，存储ethereum erc20 abi信息
func GetABIToDbOnEth(dbClient *sql.DB, e *evm.EthClient, chain string) error {
	addressList, err := db.GetContractAddressAll(dbClient, chain)
	if err != nil {
		return fmt.Errorf("fail get contract address list by db: %v", err)
	}
	for _, addr := range addressList {
		log.Println(addr)
		//获取addr Abi
		ABIStr, err := e.GetContractAbiOnEth(addr)
		if err != nil {
			log.Printf("%s fail get contract address abi: %s\n", addr, err.Error())
			continue
			//若合约是未被验证的状态则直接存储到数据库
		} else if ABIStr == constants.ABI_NO {
			err = db.SaveContractABI(dbClient, chain, addr, ABIStr, "")
			if err != nil {
				return fmt.Errorf("%s contract source code not verified，fail save contract abi:%v", addr, err)
			}
			continue
		}
		isErc20, err := processErc20(dbClient, addr, ABIStr)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		//如果不是erc20合约
		if !isErc20 {
			success, err := processContractAndStore(dbClient, e, addr, "OpenZeppelin's Unstructured", "")
			if success {
				continue
			}
			success, err = processContractAndStore(dbClient, e, addr, "eip1967", "")
			if success {
				continue
			}
			success, err = processContractAndStore(dbClient, e, addr, "eip1822", "")
			if success {
				continue
			}
			success, err = processContractAndStore(dbClient, e, addr, "eip897", ABIStr)
			if success {
				continue
			}
			if err != nil {
				log.Println("fail check proxy contract:", err.Error())
			}
		}
	}

	return nil
}
