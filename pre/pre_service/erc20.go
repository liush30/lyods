package pre_service

import (
	"database/sql"
	"fmt"
	"log"
	"lyods-adsTool/db"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pre/eip/check"
	"lyods-adsTool/services/eth"
)

func processContractAndStore(dbClient *sql.DB, e *eth.EthClient, addr, key, ABIStr string) (bool, error) {
	proxyAddress := ""
	var err error
	switch key {
	case "eip1967":
		proxyAddress, err = check.IsEIP1967(e, addr)
	case "eip1822":
		proxyAddress, err = check.IsEIP1822(e, addr)
	case "eip897":
		proxyAddress, err = check.IsEIP897(e, addr, ABIStr)
	case "OpenZeppelin's Unstructured":
		proxyAddress, err = check.IsZeppelinsUnStructStorage(e, addr)
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
	isErc20, err := check.IsERC20(ABIStr)
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

func GetABIToDbOnEth(dbClient *sql.DB, e *eth.EthClient) error {
	addressList, err := db.GetContractAddressAll(dbClient, constants.DB_CHAIN_ETH)
	if err != nil {
		return err
	}
	for _, addr := range addressList {
		//获取addr Abi
		ABIStr, err := e.GetContractAbiOnEth(addr)
		if err != nil {
			log.Printf("%s fail get contract address abi: %s\n", addr, err.Error())
			continue
			//若合约是未被验证的状态则直接退出此次循环
		} else if ABIStr == constants.ABI_NO {
			log.Printf("%s contract source code not verified\n", addr)
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
				log.Println(err.Error())
			}
		}
	}

	return nil
}
