package evm

import (
	"context"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"io"
	"log"
	"lyods-adsTool/pkg/constants"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// EVMClient ethereum查询风险账号相关操作
type EVMClient struct {
	*ethclient.Client
	Key                []string
	RequestCountSecond int
	RequestCountDay    int
	LastRequestTime    time.Time
	KeyIndex           int
	HTTPClient         *http.Client
	Chain              string
}

// CreateEvmClient constants.HTTP_ALCHEMY_ETH
func CreateEvmClient(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("failed to connect:", err)
		return nil
	}
	return client
}
func CreateBnbClient() *ethclient.Client {
	client, err := ethclient.Dial(constants.HTTP_OMNIATECH_BSC)
	if err != nil {
		log.Fatal("fail to connect:", err)
		return nil
	}
	return client
}
func (e *EVMClient) CheckRequestStatus() bool {
	// 如果当天总请求超过最大请求，切换key
	if e.RequestCountDay >= constants.ETH_MAX_DAY {
		//切换key
		isChange := e.ChangeKey()
		//切换失败
		if !isChange {
			return false
		}
		e.RequestCountDay = 0
		e.RequestCountSecond = 0
	}
	//如果请求计数超过限制，等待1秒
	if e.RequestCountSecond >= constants.ETH_MAX_SECOND {
		//计算自上一次请求以来的时间间隔
		elapsed := time.Since(e.LastRequestTime)
		//若小于1秒，我们在 time.Sleep 中等待余下的时间。
		if elapsed < time.Second {
			sleepTime := time.Second - elapsed
			time.Sleep(sleepTime)
		}
		//重置计数器和时间戳
		e.RequestCountSecond = 0
		e.LastRequestTime = time.Now()
	}
	return true
}

func (e *EVMClient) SendHTTPRequest(url string) (*http.Response, error) {
	isNormalStatus := e.CheckRequestStatus()
	if !isNormalStatus {
		return nil, fmt.Errorf("the number of requests exceeds the limit")
	}

	for retry := 0; retry < 3; retry++ {
		resp, err := e.HTTPClient.Get(url)
		if err != nil {
			if retry < 2 {
				// 如果请求失败，等待3秒后重试
				time.Sleep(3 * time.Second)
			} else {
				return nil, fmt.Errorf("send http request error after 3 retries: %v", err)
			}
		} else if resp.StatusCode == http.StatusOK {
			return resp, nil
		} else {
			if retry < 2 {
				// 如果状态码不是200，等待3秒后重试
				time.Sleep(3 * time.Second)
			} else {
				return nil, fmt.Errorf("status code is %s after 3 retries", strconv.Itoa(resp.StatusCode))
			}
		}
	}

	return nil, fmt.Errorf("maximum number of retries reached")
}

// GetKey 获得key值
func (e *EVMClient) GetKey() string {
	return e.Key[e.KeyIndex]
}

// ChangeKey 切换key
func (e *EVMClient) ChangeKey() bool {
	index := e.KeyIndex + 1
	//index超出范围
	if index > len(e.Key)-1 {
		return false
	}
	e.KeyIndex = index
	return true
}

// getContractAbiUrl 返回etherScan中查询指定合约地址的abiUrl
func (e *EVMClient) getContractAbiUrl(addr string) string {
	switch e.Chain {
	case constants.CHAIN_ETH:
		return constants.ETH_ABI + e.GetKey() + "&address=" + addr
	case constants.CHAIN_BSC:
		return constants.BSC_ENDPOINTS + constants.BSC_ABI + e.GetKey() + "&address=" + addr
	case constants.CHAIN_ARB:
		return constants.ARB_ENDPOINTS + constants.ARB_ABI + e.GetKey() + "&address=" + addr
	default:
		return ""
	}
}

// GetContractAbiOnEth 获得合约地址的abi-GetRiskListOnContractAddr
func (e *EVMClient) GetContractAbiOnEth(addr string) (string, error) {
	var err error
	//发送http请求，查询到合约的abi
	resp, err := e.SendHTTPRequest(e.getContractAbiUrl(addr))
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

// CallContractMethod 调用
func (e *EVMClient) CallContractMethod(contractAddress string, contractABIJSON string, methodName string, args ...interface{}) ([]interface{}, error) {
	// 创建合约地址
	contractAddressObj := common.HexToAddress(contractAddress)
	var msg ethereum.CallMsg
	// 创建ABI对象（使用合约的ABI定义）
	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		return nil, fmt.Errorf("abi error :%v\n", err.Error())
	}
	var callData []byte
	if len(args) != 0 {
		callData, err = contractABI.Pack(methodName, args)
	} else {
		callData, err = contractABI.Pack(methodName)
	}
	// 调用查询方法
	if err != nil {
		return nil, fmt.Errorf("pack error :%v\n", err.Error())
	}
	msg = ethereum.CallMsg{
		To:   &contractAddressObj,
		Data: callData,
	}
	// 执行调用
	callResult, err := e.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call error :%v\n", err.Error())
	}
	resultInterface, err := contractABI.Unpack(methodName, callResult)
	if err != nil {
		return nil, fmt.Errorf("unpack error :%v\n", err.Error())
	}
	return resultInterface, nil
}

//func (e *EVMClient) GetTokenSymbol(contractAddress common.Address) (string, error) {
//	symbolData, err := e.CallContract(context.Background(), ethereum.CallMsg{
//		To:   &contractAddress,
//		Data: common.Hex2Bytes("0x95d89b41"), // 使用ERC-20标准的symbol函数签名
//	}, nil)
//	if err != nil {
//		return "", err
//	}
//	return string(symbolData), nil
//}
//
//func (e *EVMClient) GetTokenDecimal(contractAddress common.Address) (int, error) {
//	decimalData, err := e.CallContract(context.Background(), ethereum.CallMsg{
//		To:   &contractAddress,
//		Data: common.Hex2Bytes("0x313ce567"), // 使用ERC-20标准的decimals函数签名
//	}, nil)
//	if err != nil {
//		return 0, err
//	}
//
//}

// CallContractMethod 调用
//func (e *EVMClient) CallContractMethod(msg ethereum.CallMsg) ([]interface{}, error) {
//	// 创建合约地址
//	// 执行调用
//	callResult, err := e.CallContract(context.Background(), msg, nil)
//	if err != nil {
//		return nil, fmt.Errorf("call error :%v\n", err.Error())
//	}
//	resultInterface, err := contractABI.Unpack(methodName, callResult)
//	if err != nil {
//		return nil, fmt.Errorf("unpack error :%v\n", err.Error())
//	}
//	return resultInterface, nil
//}
//func GetCallMsgABI(contractAddress string, contractABIJSON string, methodName string, args ...interface{}) (ethereum.CallMsg, error) {
//	contractAddressObj := common.HexToAddress(contractAddress)
//	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))
//	if err != nil {
//		return ethereum.CallMsg{}, fmt.Errorf("abi error :%v\n", err.Error())
//	}
//	var callData []byte
//	if len(args) != 0 {
//		callData, err = contractABI.Pack(methodName, args)
//	} else {
//		callData, err = contractABI.Pack(methodName)
//	}
//	// 调用查询方法
//	if err != nil {
//		return ethereum.CallMsg{}, fmt.Errorf("pack error :%v\n", err.Error())
//	}
//	msg := ethereum.CallMsg{
//		To:   &contractAddressObj,
//		Data: callData,
//	}
//	return msg, nil
//}
//func GetCallMsg(contractAddress string) {
//
//}
