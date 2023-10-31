// Package eth 查询ethereum交易信息及层级信息
package eth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"log"
	"lyods-adsTool/db"
	"lyods-adsTool/domain"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

// GetTxListOnEth 查询指定外部账户的所有交易信息
func (e *EthClient) GetTxListOnEth(addr, startBlock string, c *es.ElasticClient) ([]domain.EsTrans, string, error) {
	// 获取URL以获取指定账户的所有交易信息
	url := getNormalUrlEth(addr, startBlock, e.GetKey())

	// 发送HTTP请求
	resp, err := e.SendHTTPRequest(url)
	if err != nil {
		return nil, "", fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		return nil, "", fmt.Errorf("io read error: %v", err)
	}

	taskList := make([]domain.EsTrans, 0)
	contractAbiMap := make(map[common.Address][]byte)
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

		// 查询交易内部的trace tx信息
		traceTran, err := GetTraceTransaction(trans.Hash)
		if err != nil {
			log.Printf("Error getting trace transaction: %v", err)

			return
		}
		trans.InternalTx = traceTran

		db := db.GetDb()
		defer db.Close()

		// 获取交易内日志信息
		logs, err := e.ParseTransReceiptByHash(db, trans.Hash, &contractAbiMap)
		if err != nil {
			log.Printf("Error parsing transaction receipt: %v", err)
			return
		}
		trans.Logs = logs

		// 查询交易内部的ERC20转账交易信息
		taskList = append(taskList, trans)
		hashCount++
		//获得区块号
		if hashCount == constants.ETH_MAX_TRANS {
			lastBlockNumber = trans.BlockHeight
		}
		err = c.Insert(constants.ES_TRANSACTION, trans.Hash, trans)
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

// 处理交易基本信息，将信息存储与EsTrans中并返回
//func processTrans(value []byte, addr string) (domain.EsTrans, error) {
//	if len(value) == 0 {
//		return domain.EsTrans{}, errors.New("VALUE_IS_NIL")
//	}
//	var trans domain.EsTrans
//	var err error
//
//	toAddress, err := jsonparser.GetString(value, constants.TO_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s'_ADDRESS: %v", constants.TO_KEY, err)
//	}
//	transHash, err := jsonparser.GetString(value, constants.HASH_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.HASH_KEY, err)
//	}
//	transGas, err := jsonparser.GetString(value, constants.GAS_USED_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.GAS_USED_KEY, err)
//	}
//	transIsError, err := jsonparser.GetString(value, constants.IS_ERROR_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.IS_ERROR_KEY, err)
//	}
//	transContractAddress, err := jsonparser.GetString(value, constants.CONTRACT_ADDR_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s'_ADDRESS: %v", constants.CONTRACT_ADDR_KEY, err)
//	}
//	transFunName, err := jsonparser.GetString(value, constants.FUNCTION_NAME_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.FUNCTION_NAME_KEY, err)
//	}
//	transMethodId, err := jsonparser.GetString(value, constants.METHOD_ID_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.METHOD_ID_KEY, err)
//	}
//	transConfirm, err := jsonparser.GetString(value, constants.CONFIRMATIONS_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.CONFIRMATIONS_KEY, err)
//	}
//	transCumlGasUsed, err := jsonparser.GetString(value, constants.CUMULATIVE_GAS_USED_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.CUMULATIVE_GAS_USED_KEY, err)
//	}
//	transPrice, err := jsonparser.GetString(value, constants.GAS_PRICE_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.GAS_PRICE_KEY, err)
//	}
//	transTimeStr, err := jsonparser.GetString(value, constants.TIME_STAMP_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.TIME_STAMP_KEY, err)
//	}
//	transTime, err := strconv.Atoi(transTimeStr)
//	if err != nil {
//		return trans, fmt.Errorf("TIME_STAMP:FAIL_STRING_TO_INT: %v", err)
//	}
//	transBlockNumer, err := jsonparser.GetString(value, constants.BLOCK_NUMBER_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.BLOCK_NUMBER_KEY, err)
//	}
//	transBlockHash, err := jsonparser.GetString(value, constants.BLOCK_HASH_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.BLOCK_HASH_KEY, err)
//	}
//	transIndex, err := jsonparser.GetString(value, constants.TRANSACTION_INDEX_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.TRANSACTION_INDEX_KEY, err)
//	}
//	fromAddress, err := jsonparser.GetString(value, constants.FROM_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s'_ADDRESS: %v", constants.FROM_KEY, err)
//	}
//
//	var transValueVal big.Int
//	valueBytes, _, _, err := jsonparser.Get(value, constants.VALUE_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.VALUE_KEY, err)
//	}
//	if _, ok := transValueVal.SetString(string(valueBytes), 0); !ok {
//		return trans, errors.New("FAILED_TO_CONVERT_'value'_TO_BIG_INT")
//	}
//	transInput, err := jsonparser.GetString(value, constants.INPUT_KEY)
//	if err != nil {
//		return trans, fmt.Errorf("FAILED_TO_GET_'%s': %v", constants.INPUT_KEY, err)
//	}
//
//	return domain.EsTrans{
//		Hash:              transHash,
//		Address:           addr,
//		GasUsed:           transGas,
//		IsError:           transIsError,
//		ContractAddress:   transContractAddress,
//		FunctionName:      transFunName,
//		MethodId:          transMethodId,
//		Confirmations:     transConfirm,
//		CumulativeGasUsed: transCumlGasUsed,
//		GasPrice:          transPrice,
//		TxIndex:           transIndex,
//		Time:              int64(transTime),
//		BlockHeight:       transBlockNumer,
//		BlockHash:         transBlockHash,
//		Value:             transValueVal,
//		Inputs: []domain.InputsTrans{
//			{
//				Witness: transInput, // 判断转账地址是否为合约地址，如果input值为0x，则说明该转账对象为普通地址。
//				Addr:    fromAddress,
//			},
//		},
//		Out: []domain.OutTrans{
//			{
//				Addr: toAddress,
//			},
//		},
//	}, nil
//}

// 处理交易基本信息，将信息存储于 EsTrans 中并返回
func processTrans(value []byte, addr string) (domain.EsTrans, error) {
	if len(value) == 0 {
		return domain.EsTrans{}, errors.New("VALUE_IS_NIL")
	}

	var trans domain.EsTrans

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
		{constants.TIME_STAMP_KEY, &trans.Time},
		{constants.BLOCK_NUMBER_KEY, &trans.BlockHeight},
		{constants.BLOCK_HASH_KEY, &trans.BlockHash},
		{constants.TRANSACTION_INDEX_KEY, &trans.TxIndex},
		{constants.FROM_KEY, &trans.Inputs[0].Addr},
		{constants.VALUE_KEY, &trans.Value},
		{constants.INPUT_KEY, &trans.Inputs[0].Witness},
		{constants.TIME_STAMP_KEY, &trans.Time},
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
		case *big.Int:
			if _, ok := v.SetString(val, 0); !ok {
				return trans, fmt.Errorf("failed to convert '%s' to big int", key)
			}
		case *int64: // 处理 int64 类型字段
			intVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return trans, fmt.Errorf("failed to convert '%s' to int64", key)
			}
			*v = intVal
		}
	}
	// 设置额外字段
	trans.Address = addr
	return trans, nil
}

// GetTraceTransaction 获取指定交易hash的trace交易信息
func GetTraceTransaction(hash string) ([]domain.InternalTxn, error) {
	var iTxList []domain.InternalTxn
	url := getTraceTransactionUrl()
	payload := map[string]interface{}{
		"id":      constants.RPC_ID,
		"jsonrpc": constants.RPC_VERSION,
		"method":  constants.RPC_METHOD_TRACE,
		"params":  []string{hash},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("Marshal Error:", err.Error())
		return nil, err
	}
	////创建http请求
	req, err := http.NewRequest(constants.HTTP_POST, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	//res, _ := http.DefaultClient.Do(req)
	res, err := utils.DoRequest(req)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Io Read Error:", err)
		return nil, err
	}
	//遍历结果集，将查询到的信息存储于iTxList
	var continueLoop = true
	_, arrayErr := jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if !continueLoop {
			return // 如果已经发生错误，退出循环
		}
		iTx, err := processTraceTrans(value)
		if err != nil {
			log.Println("processTraceTrans error:", err.Error())
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

//// 解析交易内部trace交易信息
//func processTraceTrans(value []byte) (domain.InternalTxn, error) {
//	callType, err := jsonparser.GetString(value, "action", "callType")
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get callType")
//	}
//	utils.IsErrorFloat(err, "Fail get callType")
//	fromAddress, err := jsonparser.GetString(value, "action", "from")
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get from")
//	}
//	toAddress, err := jsonparser.GetString(value, "action", "to")
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get to")
//	}
//	valueVal, err := jsonparser.GetString(value, "action", "value")
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get value")
//	}
//	input, err := jsonparser.GetString(value, "action", "input")
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get input")
//	}
//	output, err := jsonparser.GetString(value, "result", "output")
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get output")
//	}
//	traceBytes, _, _, err := jsonparser.Get(value, "traceAddress")
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get traceAddress")
//	}
//	var traceInts []int64
//	err = json.Unmarshal(traceBytes, &traceInts)
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get traceAddress")
//	}
//	subtraces, err := jsonparser.GetInt(value, "subtraces")
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail get subtraces")
//	}
//	valueInt, err := utils.HexToBigInt(valueVal)
//	if err != nil {
//		return domain.InternalTxn{}, utils.HandleError(err, "Fail hex to big int")
//	}
//	//traceInts := utils.BytesToInt64Slice(traceBytes)
//	return domain.InternalTxn{
//		CallType:     callType,
//		FromAddr:     fromAddress,
//		ToAddr:       toAddress,
//		Value:        *valueInt,
//		InputTx:      input,
//		OutputTx:     output,
//		TraceAddress: traceInts,
//		SubTraces:    subtraces,
//		Id:           utils.GenerateTransactionID("call_", traceInts),
//	}, nil
//}

// 解析交易内部trace交易信息
func processTraceTrans(value []byte) (domain.InternalTxn, error) {
	type field struct {
		constantKeys []string
		field        interface{}
	}
	var internalTx domain.InternalTxn
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
		val, err := jsonparser.GetString(value, f.constantKeys...)
		if err != nil {
			return domain.InternalTxn{}, utils.HandleError(err, "fail get "+strings.Join(f.constantKeys, "."))
		}

		switch v := f.field.(type) {
		case *string:
			*v = val
		case *big.Int:
			valueInt, err := utils.HexToBigInt(val)
			if err != nil {
				return domain.InternalTxn{}, utils.HandleError(err, "fail hex to big int")
			}
			*v = *valueInt
		case *[]int64:
			traceBytes, _, _, err := jsonparser.Get(value, f.constantKeys...)
			if err != nil {
				return domain.InternalTxn{}, utils.HandleError(err, "fail get "+strings.Join(f.constantKeys, "."))
			}
			err = json.Unmarshal(traceBytes, v)
			if err != nil {
				return domain.InternalTxn{}, utils.HandleError(err, "fail unmarshal traceAddress")
			}
		}
	}

	internalTx.TraceAddress = utils.JoinInt64SliceToString(internalTx.TraceAddressInt, "_")
	internalTx.Id = fmt.Sprintf("call_%s", internalTx.TraceAddress)
	return internalTx, nil
}

// GetTransOnEthOrBnbOrArb 根据指定地址查询交易信息
//func GetTransOnEthOrBnbOrArb(chainType uint, addr string) ([]string, error) {
//	//time.Sleep(time.Millisecond * 200) //每请求一次，休眠0.2秒  -一分钟仅能请求五次
//	var err error
//	var subList []string //风险子名单信息
//	//用于去除重复数据
//	temp := map[string]struct{}{}
//	//获取url-根据地址查询以太坊普通交易信息
//	url := getNormalTransUrl(chainType, addr)
//	if url == "" || len(url) == 0 {
//		log.Println("Fail get url")
//		return nil, errors.New("Fail get url")
//	}
//	//发送HTTP请求
//	resp, err := MClient.Get(url)
//	if err != nil || resp.StatusCode != http.StatusOK {
//		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
//		return nil, err
//	}
//	defer resp.Body.Close()
//	//读取数据
//	body, err := io.ReadAll(resp.Body)
//	if err != nil || body == nil {
//		log.Fatal("Io Read Error:", err.Error())
//		return nil, err
//	}
//	//遍历result的每一条交易信息，将交易信息存储到es
//	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//		//首先判断该交易是否为转出交易:接收方地址是否为addr
//		toAddress, err := jsonparser.GetString(value, "to")
//		if err != nil {
//			log.Fatal("Fail get from")
//			return
//		}
//		//若该交易为转出交易，获取其余交易信息
//		if strings.EqualFold(toAddress, addr) {
//			transHash, err := jsonparser.GetString(value, "hash")
//			IsError(err, "Fail get hash")
//			transGas, err := jsonparser.GetString(value, "gasUsed")
//			IsError(err, "Fail get gasUsed")
//			transIsError, err := jsonparser.GetString(value, "isError")
//			IsError(err, "Fail get isError")
//			transContractAddress, err := jsonparser.GetString(value, "contractAddress")
//			IsError(err, "Fail get contractAddress")
//			transFunName, err := jsonparser.GetString(value, "functionName")
//			IsError(err, "Fail get functionName")
//			transMethodId, err := jsonparser.GetString(value, "methodId")
//			IsError(err, "Fail get methodId")
//			transConfirm, err := jsonparser.GetString(value, "confirmations")
//			IsError(err, "Fail get confirmations")
//			transCumlGasUsed, err := jsonparser.GetString(value, "cumulativeGasUsed")
//			IsError(err, "Fail get cumulativeGasUsed")
//			transPrice, err := jsonparser.GetString(value, "gasPrice")
//			IsError(err, "Fail get gasPrice")
//			transTimeStr, err := jsonparser.GetString(value, "timeStamp")
//			IsError(err, "Fail get timeStamp")
//			transTime, err := strconv.Atoi(transTimeStr)
//			IsError(err, "time string to int")
//			transBlockNumer, err := jsonparser.GetString(value, "blockNumber")
//			IsError(err, "Fail get blockNumber")
//			transBlockHash, err := jsonparser.GetString(value, "blockHash")
//			IsError(err, "Fail get blockHash")
//			transIndex, err := jsonparser.GetString(value, "transactionIndex")
//			IsError(err, "Fail get transactionIndex")
//			fromAddress, err := jsonparser.GetString(value, "from")
//			IsError(err, "Fail get from")
//			transValue, err := jsonparser.GetString(value, "value")
//			IsError(err, "Fail get value")
//			transInput, err := jsonparser.GetString(value, "input")
//			IsError(err, "Fail get input")
//			transInfo := domain.EsTrans{
//				Hash:              transHash,
//				Address:           addr,
//				TxType:            constants.TRANS_TYPE_NORMAL,
//				GasUsed:           transGas,
//				IsError:           transIsError,
//				ContractAddress:   transContractAddress,
//				FunctionName:      transFunName,
//				MethodId:          transMethodId,
//				Confirmations:     transConfirm,
//				CumulativeGasUsed: transCumlGasUsed,
//				GasPrice:          transPrice,
//				TxIndex:           transIndex,
//				Time:              int64(transTime),
//				BlockHeight:       transBlockNumer,
//				BlockHash:         transBlockHash,
//				Inputs: []domain.InputsTrans{
//					domain.InputsTrans{
//						Witness: transInput,
//						Addr:    fromAddress,
//						Value:   transValue,
//					},
//				},
//				Out: []domain.OutTrans{
//					domain.OutTrans{
//						Addr: toAddress,
//					},
//				},
//			}
//			//保存该交易信息
//			err = es.Insert(param.ADDRESS_TRANS_LIST, transHash, transInfo)
//			if err != nil {
//				return
//			}
//			/*
//				----------------------------------------------------------------------------------------
//				判断地址是否在白名单中，若在白名单中直接忽略该地址的处理操作，若不存在，再进行操作处理
//				----------------------------------------------------------------------------------------
//			*/
//
//			//判断转账地址是否为合约地址,如果input值为0x，则说明该转账对象为普通地址,否则为合约对象,（且未发生错误的交易），只处理有实际金额交易的账户，若交易金额为0暂不处理
//			//只对以下情况进行处理操作：
//			//该交易转出的地址为普通地址,且该交易未发生错误，并存在实际金额交易
//			if transInput == "0x" && transValue != "0" && transIsError == "0" {
//				//将地址存于子名单中，并将该地址存入到风险名单信息中
//				if _, ok := temp[toAddress]; !ok {
//					temp[toAddress] = struct{}{}
//					subList = append(subList, toAddress)
//					es.UpdateOrAddSourceOnTrans(toAddress, addr)
//				}
//				//否则为合约地址，且交易未发生错误，存在实际的金额交易
//			} else if transInput != "0x" && transIsError == "0" && transValue != "0" {
//				//获取该合约地址的交易总次数
//				//1.
//
//			}
//		}
//
//	}, "result")
//	if err != nil {
//		log.Println("ArrayEach result:", err.Error())
//		return nil, err
//	}
//	return subList, nil
//}
//
//// AddTranAndSourByAddrOnEth 根据指定地址添加（or 更新）指定层级的子风险名单以及地址信息
//func AddTranAndSourByAddrOnEth(addr string, level int) error {
//	arrList, err := GetTransOnEthOrBnbOrArb(constants.CHAIN_ETHEREUM, addr)
//	if err != nil {
//		log.Println("GetTransOnEthOrBnbOrArb:fail get sublist")
//		return err
//	}
//	resultList, _, _ := GetListByListOnEth(arrList, constants.INIT_LEVEL+1, level)
//	if resultList == nil || len(resultList) == 0 {
//		log.Println("GetListByListOnEth:fail get next addrList")
//		return errors.New("GetListByListOnEth:fail get next addrList")
//	}
//	return nil
//}
//
//// GetListByListOnEth 获取next层所有子名单信息
//func GetListByListOnEth(list []string, n, max int) ([]string, int, int) {
//	fmt.Println("第", n, "层")
//	//按层次存储子名单信息
//	if n <= max {
//		var allList []string
//		for _, v := range list {
//			nextList, err := GetTransOnEthOrBnbOrArb(constants.CHAIN_ETHEREUM, v)
//			if err != nil {
//				log.Println("GetTransOnEthOrBnbOrArb level ", n, " error:"+err.Error())
//				return GetListByListOnEth(nil, max+1, max)
//			}
//
//			//将查询到的子名单信息汇总
//			allList = append(allList, nextList...)
//		}
//		return GetListByListOnEth(allList, n+1, max)
//	}
//	return list, n, max
//}
//
//// GetTransInOnEthOrBnbOrArb 根据指定地址查询内部交易信息
//func GetTransInOnEthOrBnbOrArb(chainType uint, addr string) (domain.TransactionEthOrBnb, error) {
//	var err error
//	//内部交易信息
//	var trans domain.TransactionEthOrBnb
//	//获取url-根据地址查询以太坊内部交易信息
//	url := getInternalTransUrl(chainType, addr)
//	//若url为空，说明输入的chain type不正确
//	if url == "" || len(url) == 0 {
//		log.Println("incorrect chain type")
//		return domain.TransactionEthOrBnb{}, errors.New("incorrect chain type")
//	}
//	//发送HTTP请求
//	resp, err := MClient.Get(url)
//	if err != nil || resp.StatusCode != http.StatusOK {
//		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
//		return domain.TransactionEthOrBnb{}, err
//	}
//	defer resp.Body.Close()
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal("Io Read Error:", err.Error())
//		return domain.TransactionEthOrBnb{}, err
//	}
//	//将获取的数据反序列化为结构体信息
//	err = json.Unmarshal(body, &trans)
//	if err != nil {
//		log.Fatal(err.Error())
//		return domain.TransactionEthOrBnb{}, err
//	}
//	return trans, err
//}
//

// GetTransByContract 根据合约地址查询交易信息，追踪出合约中可疑交易数据，并将相关交易地址存于风险合约内部
func GetTransByContract(address string) {

}

// 根据chain类型获取地址的内部交易
func getInternalTransUrl(chainType uint, addr string) string {
	if chainType == constants.CHAIN_BNB {
		return getInterUrlBnb(addr)
	} else if chainType == constants.CHAIN_ETHEREUM {
		return getInternalEthUrl(addr)
	} else if chainType == constants.CHAIN_ARBITRUM {
		return getInterUrlArb(addr)
	}
	return ""
}

// 根据地址获取ethereum中usdt交易请求url
func getUsdtTransUrl(addr string) string {
	return constants.API_ETH_USDT + addr
}

// bnb根据地址获取内部交易请求url
func getNormalUrlBnb(addr string) string {
	return constants.API_BNB_TRANS + addr
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
func getNormalUrlEth(addr, startBlock, key string) string {
	return constants.ETH_ADDR_ETHSCAN + startBlock + "&address=" + addr + "&apikey=" + key
}

// arbitrum根据地址获取普通交易url
func getNormalUrlArb(addr string) string {
	return constants.API_ARB_TRANS + addr
}

// arbitrum根据地址获得内部交易url
func getInterUrlArb(addr string) string {
	return constants.API_ARB_INTRANS + addr
}
