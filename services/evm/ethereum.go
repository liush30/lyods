package evm

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"io"
	"log"
	"lyods-adsTool/db"
	"lyods-adsTool/domain"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"math/big"
	"strconv"
	"strings"
)

type ContractMap struct {
	ABI      string
	IsErc20  bool
	Decimal  int
	hasProxy bool
	ProxyABI string
	IsCheck  bool
}
type field struct {
	constantKeys []string
	field        interface{}
}

// GetTxList 查询指定外部账户的所有交易信息
func (e *EVMClient) GetTxList(c *es.ElasticClient, cbClient *ChainBaseClient, addr, startBlock string) ([]domain.EsTrans, string, error) {
	//查询该地址是否为合约地址
	isContract, err := e.IsContractAddress(addr)
	fmt.Println("GetTxListOnEth:", addr, "isContract:", isContract)
	if err != nil {
		return nil, "", fmt.Errorf("GetTxListOnEth IsContractAddress error: %v", err)
	}
	// 获取URL以获取指定账户的所有交易信息
	url := e.getNormalUrl(addr, startBlock, e.GetKey())
	//判断url是否为空
	//if url ==""{
	//
	//}
	//log.Println("GetTxListOnEth:", url)
	// 发送HTTP请求
	resp, err := e.SendHTTPRequest(url)
	if err != nil {
		return nil, "", fmt.Errorf("GetTxListOnEth request error: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		return nil, "", fmt.Errorf("GetTxListOnEth io read error: %v", err)
	}

	taskList := make([]domain.EsTrans, 0)
	contractAbiMap := make(map[common.Address]ContractMap)
	hashCount := 0        //记录获取的交易总量
	lastBlockNumber := "" //最后一个交易的区块数
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Printf("Error processing transaction: %v", err)
			return
		}

		// 处理交易信息
		trans, err := processTrans(value, addr)
		if err != nil {
			log.Printf("Error processing transaction: %v", err)
			return
		}
		var toAddressIsContract bool
		if !isContract {
			//判断to是否为合约地址
			toAddressIsContract, err = e.IsContractAddress(trans.Out[0].Addr)
			if err != nil {
				log.Printf("Error getting contract address: %v", err)
				return
			}
		}
		var traceTran []domain.InternalTxn
		//查询合约地址且交易未出错的交易或to地址为合约地址的交易中的内部交易信息
		if isContract && trans.IsError == "0" || toAddressIsContract && trans.IsError == "0" {
			// 查询交易内部的trace tx信息
			traceTran, err = e.GetTraceTransaction(cbClient, trans.Hash)
			if err != nil {
				log.Printf("Error getting trace transaction: %v", err)
				return
			}
		}
		if len(traceTran) == 1 {
			trans.InternalTx = nil
		} else {
			trans.InternalTx = traceTran
		}

		db := db.GetDb()
		defer db.Close()

		// 获取交易内日志信息和交易内部erc20转账交易
		logs, erc20Txs, err := e.ParseTransReceiptByHash(db, trans.Hash, &contractAbiMap)
		if err != nil {
			log.Printf("Error parsing %s receipt: \n%v", trans.Hash, err)
			return
		}
		trans.Logs = logs
		trans.Erc20Txn = erc20Txs
		trans.Chain = constants.CHAIN_ETH
		blockNumberInt, _ := big.NewInt(0).SetString(trans.BlockHeight, 10)
		//获取当时交易时的余额
		balanceInt, err := e.GetBalance(addr, blockNumberInt)
		if err != nil {
			log.Println("GetBalance Error:", err.Error())
			return
		}
		balanceFloat, _, err := WeiToEth(balanceInt)
		if err != nil {
			log.Println("WeiToEth Error:", err.Error())
			return
		}
		trans.Balance = balanceFloat
		// 查询交易内部的ERC20转账交易信息
		taskList = append(taskList, trans)
		hashCount++
		//获得区块号
		if hashCount == constants.ETH_MAX_TRANS {
			lastBlockNumber = trans.BlockHeight
		}
		log.Println("hashCount:", hashCount, "hash:", trans.Hash)
		err = c.Insert(constants.ES_TRANSACTION, strings.ToLower(trans.Hash), trans)
		if err != nil {
			log.Println("Insert Transaction Error:", err.Error())
			return
		}

	}, "result")

	if err != nil {
		log.Println("ArrayEach result:", err.Error())
		return nil, "", err
	}

	return taskList, lastBlockNumber, nil
}

// 处理交易基本信息，将信息存储于 EsTrans 中并返回
func processTrans(value []byte, addr string) (domain.EsTrans, error) {
	if len(value) == 0 {
		return domain.EsTrans{}, errors.New("VALUE_IS_NIL")
	}

	var trans domain.EsTrans
	trans.Out = make([]domain.OutTrans, 1)
	trans.Inputs = make([]domain.InputsTrans, 1)

	type field struct {
		constantKey string
		field       interface{}
	}
	fields := []field{
		{constants.TO_KEY, &trans.Out[0].Addr},
		{constants.HASH_KEY, &trans.Hash},
		{constants.GAS_USED_KEY, &trans.GasUsed},
		{constants.IS_ERROR_KEY, &trans.IsError},
		{constants.CONTRACT_ADDR_KEY, &trans.ContractAddress},
		{constants.FUNCTION_NAME_KEY, &trans.FunctionName},
		{constants.METHOD_ID_KEY, &trans.MethodId},
		{constants.CONFIRMATIONS_KEY, &trans.Confirmations},
		{constants.CUMULATIVE_GAS_USED_KEY, &trans.CumulativeGasUsed},
		{constants.GAS_PRICE_KEY, &trans.GasPrice},
		//{constants.TIME_STAMP_KEY, &trans.Time},
		{constants.BLOCK_NUMBER_KEY, &trans.BlockHeight},
		{constants.BLOCK_HASH_KEY, &trans.BlockHash},
		{constants.TRANSACTION_INDEX_KEY, &trans.TxIndex},
		{constants.FROM_KEY, &trans.Inputs[0].Addr},
		{constants.VALUE_KEY, &trans.Value},
		{constants.INPUT_KEY, &trans.Inputs[0].Witness},
		//{constants.TIME_STAMP_KEY, &trans.Time},
	}

	for _, f := range fields {
		key := f.constantKey
		val, err := jsonparser.GetString(value, key)
		if err != nil {
			return trans, fmt.Errorf("fail to get '%s': %v", key, err)
		}

		switch v := f.field.(type) {
		case *string:
			*v = val
		case *float64:
			// 将值转换为 *big.Int
			intVal, success := new(big.Int).SetString(val, 10)
			if !success {
				return domain.EsTrans{}, fmt.Errorf("failed to convert '%s' to *big.Int", key)
			}

			// 调用 WeiToEth 函数将 *big.Int 转换为 *big.Float
			floatVal, floatStr, err := WeiToEth(intVal)
			if err != nil {
				return domain.EsTrans{}, fmt.Errorf("failed to convert '%s' to *big.Float", key)
			}
			*v = floatVal
			// 设置 *big.Float 值
			trans.ValueText = floatStr
		case *int64: // 处理 int64 类型字段
			intVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return domain.EsTrans{}, fmt.Errorf("failed to convert '%s' to int64", key)
			}
			*v = intVal
		}
	}
	// 设置额外字段
	trans.Address = addr
	timeStamp, err := jsonparser.GetString(value, constants.TIME_STAMP_KEY)
	if err != nil {
		return trans, fmt.Errorf("fail to get '%s': %v", constants.TIME_STAMP_KEY, err)
	}
	var timeBigInt big.Int
	_, success := timeBigInt.SetString(timeStamp, 10)
	if !success {
		return trans, fmt.Errorf("failed to convert '%s' to *big.Int", constants.TIME_STAMP_KEY)
	}
	trans.Time = utils.UnixToTime(timeBigInt.Int64())

	return trans, nil
}

// GetTraceTransaction 获取指定交易hash的trace交易信息
func (e *EVMClient) GetTraceTransaction(cbClient *ChainBaseClient, hash string) ([]domain.InternalTxn, error) {
	var iTxList []domain.InternalTxn
	res, err := cbClient.SendHTTPRequest(hash, e.Chain)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("Io Read Error:", err)
		return nil, err
	}
	//遍历结果集，将查询到的信息存储于iTxList
	var continueLoop = true
	if len(body) == 0 {
		return nil, nil
	}
	//log.Println(string(body))
	_, arrayErr := jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if !continueLoop {
			return // 如果已经发生错误，退出循环
		}
		iTx, err := processTraceTrans(value)
		if err != nil {
			log.Printf("%s processTraceTrans error:%v", hash, err.Error())
			continueLoop = false
			return
		}
		iTxList = append(iTxList, iTx)

	}, "result")
	if arrayErr != nil {
		log.Println("ArrayEach result:", arrayErr.Error())
		return nil, arrayErr
	}
	return iTxList, nil
}
func processTraceTrans(value []byte) (domain.InternalTxn, error) {
	//获得交易类型
	typeStr, err := jsonparser.GetString(value, constants.TYPE_KEY)
	if err != nil {
		return domain.InternalTxn{}, fmt.Errorf("fail to get '%s': %v", constants.TYPE_KEY, err)
	}
	var internalTx domain.InternalTxn
	internalTx.Type = typeStr
	if typeStr == constants.TRACR_TYPE_CALL {
		return processCallTraceTrans(value, &internalTx)
	}
	return processCreateTrace(value, &internalTx)
}
func processCreateTrace(value []byte, internalTx *domain.InternalTxn) (domain.InternalTxn, error) {
	fields := []field{
		{[]string{constants.ACTION_KEY, constants.INIT_KEY}, &internalTx.Init},
		{[]string{constants.ACTION_KEY, constants.FROM_KEY}, &internalTx.FromAddr},
		{[]string{constants.ACTION_KEY, constants.VALUE_KEY}, &internalTx.Value},
		{[]string{constants.RESULT_KEY, constants.ADDRESS_KEY}, &internalTx.Address},
		{[]string{constants.RESULT_KEY, constants.CODE_KEY}, &internalTx.Code},
		{[]string{constants.TRACE_ADDRESS_KEY}, &internalTx.TraceAddressInt},
		{[]string{constants.SUBTRACES_KEY}, &internalTx.SubTraces},
	}
	for _, f := range fields {
		val, _, _, err := jsonparser.Get(value, f.constantKeys...)
		if err != nil {
			return domain.InternalTxn{}, fmt.Errorf("fail get '%s': %v", f.constantKeys, err)
		}
		switch v := f.field.(type) {
		case *string:
			*v = string(val)
		case *int64:
			intVal, err := strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				return domain.InternalTxn{}, fmt.Errorf("fail to parse int64 from '%s': %v", string(val), err)
			}
			*v = intVal
		case *float64:
			valueInt, err := utils.HexToBigInt(big.NewInt(0).SetBytes(val).String())
			if err != nil {
				return domain.InternalTxn{}, fmt.Errorf("fail convert '%s' to *big.Int:%v", string(val), err)
			}
			// 调用 WeiToEth 函数将 *big.Int 转换为 *big.Float
			floatVal, floatStr, err := WeiToEth(valueInt)
			if err != nil {
				return domain.InternalTxn{}, fmt.Errorf("failed to convert '%s' to *big.Float:%v", string(val), err)
			}
			*v = floatVal
			// 设置 *big.Float 值
			internalTx.ValueText = floatStr
		case *[]int64:
			// 解码 JSON 数组
			err = json.Unmarshal(val, v)
			if err != nil {
				return domain.InternalTxn{}, fmt.Errorf("fail unmarshal '%s': %v", string(val), err)
			}
		}
	}
	internalTx.TraceAddress = utils.JoinInt64SliceToString(internalTx.TraceAddressInt, "_")
	internalTx.Id = fmt.Sprintf("call_%s", internalTx.TraceAddress)
	return *internalTx, nil
}

// 解析内部交易类型为cal的内部trace交易信息
func processCallTraceTrans(value []byte, internalTx *domain.InternalTxn) (domain.InternalTxn, error) {
	fields := []field{
		{[]string{constants.ACTION_KEY, constants.CALL_TYPE_KEY}, &internalTx.CallType},
		{[]string{constants.ACTION_KEY, constants.FROM_KEY}, &internalTx.FromAddr},
		{[]string{constants.ACTION_KEY, constants.TO_KEY}, &internalTx.ToAddr},
		{[]string{constants.ACTION_KEY, constants.VALUE_KEY}, &internalTx.Value},
		{[]string{constants.ACTION_KEY, constants.INPUT_KEY}, &internalTx.InputTx},
		{[]string{constants.RESULT_KEY, constants.OUTPUT_KEY}, &internalTx.OutputTx},
		{[]string{constants.TRACE_ADDRESS_KEY}, &internalTx.TraceAddressInt},
		{[]string{constants.SUBTRACES_KEY}, &internalTx.SubTraces},
	}
	for _, f := range fields {
		val, _, _, err := jsonparser.Get(value, f.constantKeys...)
		if err != nil {
			return domain.InternalTxn{}, fmt.Errorf("fail get '%s': %v", f.constantKeys, err)
		}
		switch v := f.field.(type) {
		case *string:
			*v = string(val)
		case *float64:
			valueInt, err := utils.HexToBigInt(big.NewInt(0).SetBytes(val).String())
			if err != nil {
				return domain.InternalTxn{}, fmt.Errorf("fail convert '%s' to *big.Int:%v", string(val), err)
			}
			// 调用 WeiToEth 函数将 *big.Int 转换为 *big.Float
			floatVal, floatStr, err := WeiToEth(valueInt)
			if err != nil {
				return domain.InternalTxn{}, fmt.Errorf("failed to convert '%s' to *big.Float:%v", string(val), err)
			}
			*v = floatVal
			// 设置 *big.Float 值
			internalTx.ValueText = floatStr
		case *[]int64:
			// 解码 JSON 数组
			err = json.Unmarshal(val, v)
			if err != nil {
				return domain.InternalTxn{}, fmt.Errorf("fail unmarshal '%s': %v", string(val), err)
			}
		case *int64:
			intVal, err := strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				return domain.InternalTxn{}, fmt.Errorf("fail to parse int64 from '%s': %v", string(val), err)
			}
			*v = intVal
		}

	}

	internalTx.TraceAddress = utils.JoinInt64SliceToString(internalTx.TraceAddressInt, "_")
	internalTx.Id = fmt.Sprintf("call_%s", internalTx.TraceAddress)
	return *internalTx, nil
}

// bnb根据地址获取内部交易请求url
func getInterUrlBnb(addr string) string {
	return constants.API_BNB_INTRANS + addr
}

// ethereum根据地址获取内部交易请求url
func getInternalEthUrl(addr string) string {
	return constants.API_ETH_INTRANS + addr
}

// ethereum根据地址获取普通交易请求url
func (e *EVMClient) getNormalUrl(addr, startBlock, key string) string {
	switch e.Chain {
	case constants.CHAIN_ETH:
		return constants.ETH_ADDR_ETHSCAN + startBlock + "&address=" + addr + "&apikey=" + key
	case constants.CHAIN_BSC:
		return constants.BSC_ENDPOINTS + constants.BSC_TX_ADDR + startBlock + "&address=" + addr + "&apikey=" + key
	default:
		return ""
	}

}

// arbitrum根据地址获得内部交易url
func getInterUrlArb(addr string) string {
	return constants.API_ARB_INTRANS + addr
}

// getInternalTxn 根据事件信息进行解析，存储到InternalTxn结构体中
// eventName 事件名称
// eventNameToValueByAddress 事件相关参数信息
func (e *EVMClient) getInternalTxn(interParam *InternalTxnParam, internalTxn *domain.InternalTxn) {
	///如果不存在address类型的参数,则跳过参数解析
	if _, ok := interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]; !ok {
		return
	}
	//仅存在address类型的参数
	if interParam.length == 1 && len(interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]) == 1 {
		//var changeMoney *big.Int
		//var v string
		//var err error
		//判断address是from还是to
		//获取交易前后金额差
		if interParam.isErc20 {
			//for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
			//changeMoney, err = e.GetERC20TokenBalanceChange(interParam.contractAddress, v, interParam.blockNumber)
			//if err != nil {
			//	log.Fatal("Fail GetERC20TokenBalanceChange:", err)
			//}
			//}
		} else {
			//for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
			//changeMoney, err = e.GetBalanceChange(interParam.blockNumber, v)
			//if err != nil {
			//	log.Fatal("Fail GetERC20TokenBalanceChange:", err)
			//}
			//}
		}
		//判断changeMoney是否大于0
		//if changeMoney.Cmp(big.NewInt(0)) > 0 {
		//	internalTxn.ToAddr = v
		//	internalTxn.FromAddr = interParam.contractAddress
		//	internalTxn.ActualAmount = *changeMoney
		//} else {
		//	internalTxn.FromAddr = v
		//	internalTxn.ToAddr = interParam.contractAddress
		//	internalTxn.ActualAmount = *changeMoney
		//}
		return
	} else if interParam.length == 2 && len(interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]) == 2 {

	}
	//首先对常规相关转账的行为进行判断
	//1.Transfer 事件解析-Transfer(address indexed from, address indexed to, uint256 value)
	if strings.EqualFold(interParam.eventName, constants.EVENT_NAME_TRANSFER) {
		fromValue, fromOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS][constants.EVENT_PARAM_FROM]
		toValue, toOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS][constants.EVENT_PARAM_TO]
		//判断是否存在(address indexed from, address indexed to, uint256 value)结构
		if len(interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]) == 2 && len(interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT]) == 1 {
			//from,to都存在,根据字段名称存储
			if fromOk && toOk {
				internalTxn.FromAddr = fromValue
				internalTxn.ToAddr = toValue
				//to字段存在，from字段不存在
			} else if fromOk && !toOk {
				internalTxn.FromAddr = fromValue
				delete(interParam.eventNameToValueByAddress, constants.EVENT_PARAM_FROM)
				for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
					internalTxn.ToAddr = v
				}
				//from字段存在，to字段存在
			} else if !fromOk && toOk {
				internalTxn.ToAddr = toValue
				delete(interParam.eventNameToValueByAddress, constants.EVENT_PARAM_TO)
				for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
					internalTxn.FromAddr = v
				}
				//from,to都不存在
			} else {
				//判断地址中哪个是from,哪个是to
				if interParam.isErc20 {
					for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
						changeMoney, err := e.GetERC20TokenBalanceChange(interParam.contractAddress, v, interParam.blockNumber)
						if err != nil {
							log.Fatal("Fail GetERC20TokenBalanceChange:", err)
						}
						//判断changeMoney是否大于0
						if changeMoney.Cmp(big.NewInt(0)) > 0 {
							internalTxn.ToAddr = v
						} else {
							internalTxn.FromAddr = v
						}
					}
				} else {
					for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
						changeMoney, err := e.GetBalanceChange(interParam.blockNumber, v)
						if err != nil {
							log.Fatal("Fail GetERC20TokenBalanceChange:", err)
						}
						//判断changeMoney是否大于0
						if changeMoney.Cmp(big.NewInt(0)) > 0 {
							internalTxn.ToAddr = v
						} else {
							internalTxn.FromAddr = v
						}
					}
				}
			}
			//判断internalTxn.FromAddr、internalTxn.ToAddr是否为空
			if len(internalTxn.FromAddr) == 0 || len(internalTxn.ToAddr) == 0 {
				log.Fatal("Fail internalTxn.FromAddr or internalTxn.ToAddr is empty")
			}
			//for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT] {
			//将v转成big.int类型
			//value, _ := new(big.Int).SetString(v, 10)
			//internalTxn.Amount = *value
			//}
			return
		} else {

		}
		// Deposit (index_topic_1 address sender, uint256 value)
	} else if strings.EqualFold(interParam.eventName, constants.EVENT_NAME_DEPOSIT) {

	}

}

// ParseTransReceiptByHash 根据交易hash解析交易中的log信息
func (e *EVMClient) ParseTransReceiptByHash(dbClient *sql.DB, hash string, contractAbiMap *map[common.Address]ContractMap) ([]domain.Logs, []domain.Erc20Txn, error) {
	//根据交易哈希查询交易的receipt信息
	receipt, err := e.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get receipt info for %v: %v", hash, err)
	}
	//若log为空 ，直接返回空
	if len(receipt.Logs) == 0 {
		return nil, nil, nil
	}
	var resultList []domain.Logs      //log信息列表
	var erc20TxList []domain.Erc20Txn //erc20转账信息列表
	//查询log信息
	if receipt.Status == 1 { //receipt status 1表示执行成功，0表示执行失败
		resultList, erc20TxList, err = e.getLogsInReceipt(dbClient, receipt.Logs, contractAbiMap)
		if err != nil {
			return nil, nil, fmt.Errorf("fail get logs info:%v", err)
		}
	}
	return resultList, erc20TxList, nil
}
func (e *EVMClient) getLogsInReceipt(dbClient *sql.DB, logs []*types.Log, contractAbiMap *map[common.Address]ContractMap) ([]domain.Logs, []domain.Erc20Txn, error) {
	var resultList []domain.Logs
	var erc20TxList []domain.Erc20Txn
	//遍历receipt log
	for _, logInfo := range logs {
		contractMap, err := e.getContractMapByLogAddress(dbClient, logInfo.Address, contractAbiMap)
		if err != nil {
			return nil, nil, err
		}
		var resultLog domain.Logs
		var erc20Tx domain.Erc20Txn
		//若合约未被验证，记录未被解析的log数据，并退出循环，解析下一个log
		if string(contractMap.ABI) == constants.ABI_NO {
			resultLog, erc20Tx, err = e.parseNoVerifyLog(logInfo, &contractMap)
			if err != nil {
				return nil, nil, fmt.Errorf("fail parse log:%v", err.Error())
			}
		} else {
			resultLog, erc20Tx, err = e.parseVerifyLog(logInfo, &contractMap)
			if err != nil {
				return nil, nil, fmt.Errorf("fail parse log:%v", err.Error())
			}
		}
		resultList = append(resultList, resultLog)
		erc20TxList = append(erc20TxList, erc20Tx)
	}
	return resultList, erc20TxList, nil
}
func (e *EVMClient) getContractMapByLogAddress(dbClient *sql.DB, logAddr common.Address, contractAbiMap *map[common.Address]ContractMap) (ContractMap, error) {
	if contractMap, ok := (*contractAbiMap)[logAddr]; ok {
		return contractMap, nil
	}

	contractMap, err := e.checkIsErc20AndABI(dbClient, logAddr.String())
	if err != nil {
		return ContractMap{}, err
	}
	//存储到contractAbiMap
	(*contractAbiMap)[logAddr] = contractMap
	return contractMap, nil
}
func (e *EVMClient) checkIsErc20AndABI(dbClient *sql.DB, addr string) (ContractMap, error) {
	//判断该地址是否存在数据库中
	isExits, err := db.ExistsToken(dbClient, addr, constants.DB_CHAIN_ETH)
	if err != nil {
		return ContractMap{}, fmt.Errorf("ExistsToken:failed to get ABI info from the database: %v", err)
	}
	var abi string
	var decimal int
	//若count>0,说明该地址是erc20合约
	if isExits {
		abi, decimal, err = db.QueryAbiAndCheckByAddress(dbClient, addr, constants.DB_CHAIN_ETH)
		if err != nil {
			return ContractMap{}, fmt.Errorf("QueryAbiAndCheckByAddress:failed to get ABI info from the database: %v", err)
		}
		return ContractMap{ABI: abi, IsErc20: true, Decimal: decimal}, nil
	}

	//若不存在,则通过地址查询合约的ABI信息
	contractAbiStr, err := e.GetContractAbiOnEth(addr)
	if err != nil {
		return ContractMap{}, err
	} else if contractAbiStr == constants.ABI_NO { //合约未被验证
		return ContractMap{ABI: contractAbiStr, IsErc20: false}, nil
	}
	//根据合约abi信息判断是否是erc20合约
	isErc20, err := IsERC20(contractAbiStr)
	if err != nil {
		return ContractMap{}, err
	}
	if isErc20 {
		return ContractMap{
			ABI:     contractAbiStr,
			IsErc20: true,
			Decimal: decimal,
		}, nil
	}
	return ContractMap{ABI: contractAbiStr, IsErc20: false}, nil
}

// 获取未验证的合约的log信息
func (e *EVMClient) parseNoVerifyLog(logInfo *types.Log, contractMap *ContractMap) (domain.Logs, domain.Erc20Txn, error) {
	var topicsInfo []domain.TopicsValStruct
	if logInfo.Topics[0].String() == constants.SIGN_TRANSFER {
		return TransferLog(logInfo, contractMap)
	}
	//遍历topics数组，将数据存储于topicsInfo中
	for i, iTopicsList := range logInfo.Topics[1:] {
		topicsInfo = append(topicsInfo, domain.TopicsValStruct{
			Key:   fmt.Sprintf("topic%d", i),
			Value: iTopicsList.String(),
		})
	}
	//将data解析为16进制数据，存储与topicsInfo中
	topicsInfo = append(topicsInfo, domain.TopicsValStruct{
		Key:   "dataHex",
		Value: common.BytesToHash(logInfo.Data).Hex(),
	})
	return domain.Logs{
		Address:   logInfo.Address.String(),
		EventInfo: logInfo.Topics[0].String(), //topics[0]是事件签名
		Topics:    topicsInfo,
	}, domain.Erc20Txn{}, nil
}

// 获取已验证合约的log信息

func (e *EVMClient) parseVerifyLog(logInfo *types.Log, contractMap *ContractMap) (domain.Logs, domain.Erc20Txn, error) {
	//若该地址是一个代理合约，先尝试使用代理合约abi进行解析
	var abiContractMap string
	if contractMap.hasProxy {
		abiContractMap = contractMap.ProxyABI
	} else {
		abiContractMap = contractMap.ABI
	}
	var contractABI abi.ABI
	//根据获得的abi信息，生成合约abi对象
	contractABI, err := abi.JSON(bytes.NewReader([]byte(abiContractMap)))
	if err != nil {
		return domain.Logs{}, domain.Erc20Txn{}, fmt.Errorf("fail create abi:%v", err)
	}
	evenSig := logInfo.Topics[0] //默认topics的第一个参数为事件的签名
	if evenSig.String() == constants.SIGN_TRANSFER {
		return TransferLog(logInfo, contractMap)
	}
	var eventInfo *abi.Event
	//获取事件签名，并根据事件签名和abi解析出事件
	eventInfo, err = contractABI.EventByID(evenSig)
	if err != nil {
		var abiOther string
		//若代理合约解析事件失败，尝试使用原合约abi解析事件
		if contractMap.hasProxy {
			abiOther = contractMap.ABI
			//如果还未检查该地址是否为代理合约
		} else if !contractMap.IsCheck {
			//尝试获取代理合约abi进行解析
			isProxy, proxyAddress, err := e.IsProxyContract(logInfo.Address.String(), contractMap.ABI)
			if err != nil {
				return domain.Logs{}, domain.Erc20Txn{}, fmt.Errorf("IsProxyContract:fail check %s proxy:%v", logInfo.Address, err)
			}
			if isProxy {
				contractMap.hasProxy = true
				proxyAbi, err := e.GetContractAbiOnEth(proxyAddress)
				if err != nil {
					return domain.Logs{}, domain.Erc20Txn{}, fmt.Errorf("GetContractAbiOnEth：fail get %s proxy abi:%v", logInfo.Address, err)
				}
				contractMap.ProxyABI = proxyAbi
				if proxyAbi == constants.ABI_NO {
					return e.parseNoVerifyLog(logInfo, contractMap)
				} else {
					abiOther = proxyAbi
				}
			} else {
				return e.parseNoVerifyLog(logInfo, contractMap)
			}
			contractMap.IsCheck = true
		} else if contractMap.IsCheck && !contractMap.hasProxy {
			return e.parseNoVerifyLog(logInfo, contractMap)
		} else {
			return domain.Logs{}, domain.Erc20Txn{}, fmt.Errorf("ParseTransReceiptByHash:Fail to get event info by %s->%v", logInfo.Topics[0].String(), err.Error())
		}
		contractABI, err = abi.JSON(bytes.NewReader([]byte(abiOther)))
		if err != nil {
			return domain.Logs{}, domain.Erc20Txn{}, fmt.Errorf("fail parse abi by abi and proxy abi:%v", err)
		}
		eventInfo, err = contractABI.EventByID(evenSig)
		if err != nil {
			return domain.Logs{}, domain.Erc20Txn{}, fmt.Errorf("ParseTransReceiptByHash:Fail to get event info by %s->%v", logInfo.Topics[0].String(), err.Error())
		}
	}
	// 存储事件信息，以参数名称和类型键值对形式存储
	eventNameToValue := make(map[string]string)
	var paramName, paramType []string
	topicIndex := 1

	// 遍历事件参数
	var nameParams []string
	for _, input := range eventInfo.Inputs {
		inputType := input.Type.String()
		inputName := input.Name
		param := fmt.Sprintf("%s %s", inputType, inputName)
		nameParams = append(nameParams, param)
		//若事件参数被标记为indexed，参数值存在于topics数组中
		if input.Indexed {
			//解析数据并将解析好的数据存储到eventName中
			eventNameToValue[inputName] = hexToData(inputType, logInfo.Topics[topicIndex])
			topicIndex++
			//若参数未被标记为indexed，数据将会存在data中,先将事件参数的名字和类型进行存储
		} else {
			paramName = append(paramName, inputName)
			paramType = append(paramType, inputType)
		}
	}
	// 构建事件信息字符串
	eventInfoStr := eventInfo.RawName + "(" + strings.Join(nameParams, ", ") + ")"
	//若data不为空，则解析data数据
	if len(logInfo.Data) >= 0 && len(paramName) > 0 && len(paramType) > 0 {
		//解析data数据-没有被标记为indexed的参数的值则存储于log.data中，需解析出值后存储于键值对中
		err = parseData(paramName, paramType, eventInfo.Name, logInfo.Data, contractABI, &eventNameToValue)
		if err != nil {
			return domain.Logs{}, domain.Erc20Txn{}, fmt.Errorf("parseData: %s", err.Error())
		}
	}
	return domain.Logs{
		Address:   logInfo.Address.String(),
		EventInfo: eventInfoStr,
		Topics:    mapToTopicsValStruct(eventNameToValue),
	}, domain.Erc20Txn{}, nil
}

// GetNonceForTransaction 获得指定交易时，调用者的nonce值
func (e *EVMClient) GetNonceForTransaction(txHash string) (uint64, error) {
	// 查询交易
	tx, _, err := e.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return 0, err
	}
	return tx.Nonce(), nil
}
