package normal

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"log"
	"lyods-adsTool/pkg/constants"
	"net/http"
)

// getContractAbiUrl 返回etherScan中查询指定合约地址的abiUrl
func getContractAbiUrl(addr string) string {
	return constants.API_ETH_ABI + addr
}

// GetContractAbiOnEth 获得合约地址的abi-GetRiskListOnContractAddr
func GetContractAbiOnEth(client *http.Client, addr string) (string, error) {
	var err error
	//发送http请求，查询到合约的abi
	resp, err := client.Get(getContractAbiUrl(addr))
	if err != nil {
		log.Printf("GetContractAbi: Fail request,http status is %v, do error ->%v\n", resp.StatusCode, err.Error())
		return "", err
	}
	defer resp.Body.Close()
	//读取数据
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		log.Printf("GetContractAbi: Io Read Error->%v\n", err.Error())
		return "", err
	}
	abi, err := jsonparser.GetString(body, "result")
	if err != nil {
		log.Printf("GetContractAbi: Fail get result->%v\n", err.Error())
		return "", err
	}
	if abi == "Contract source code not verified" {
		log.Printf("GetContractAbi:Contract source code not verified\n")
		return "", fmt.Errorf("GetContractAbi:Contract source code not verified")
	}
	return abi, nil
}
