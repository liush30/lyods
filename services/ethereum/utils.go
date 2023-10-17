package ethereum

import (
	"github.com/buger/jsonparser"
	"io"
	"log"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
)

// getContractAbiUrl 返回etherScan中查询指定合约地址的abiUrl
func getContractAbiUrl(addr string) string {
	return constants.API_ETH_ABI + addr
}

// GetContractAbiOnEth 获得合约地址的abi-GetRiskListOnContractAddr
func GetContractAbiOnEth(addr string) (string, error) {
	var err error
	//发送http请求，查询到合约的abi
	resp, err := utils.SendHTTPRequest(getContractAbiUrl(addr), constants.HTTP_GET, nil)
	if err != nil {
		log.Printf("GetContractAbi: Do Error->%v\n", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	//读取数据
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		log.Printf("GetContractAbi: Io Read Error->%v\n", err.Error())
		return "", err
	}
	//获取abi
	abi, err := jsonparser.GetString(body, "result")
	if err != nil {
		log.Printf("GetContractAbi: Fail get result->%v\n", err.Error())
		return "", err
	}
	return abi, nil
}
