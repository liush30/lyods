// Package bitcoin 查询bitcoin交易信息
// 及子层级名单，包含链：BTG、XBT、ETH、XMR、LTC、ZEC、DASH、ETC、BSV
package bitcoin

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"log"
	"lyods-adsTool/domain"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"math/big"
	"strconv"
	"strings"
)

//	func GetTxListByBlockChain(c *es.ElasticClient, addr string) ([]domain.EsTrans, error) {
//		//随机休眠几秒
//		//utils.RandomSleep()
//		var transList []domain.EsTrans
//		//获得url-获取到指定账户的所有交易信息url
//		url := getUrlByBlockChain(addr)
//		//发送http请求-根据url
//		resp, err := utils.SendHTTPRequestWithRateLimit(url, constants.HTTP_GET, nil)
//		if err != nil {
//			return nil, err
//		}
//		defer resp.Body.Close()
//		body, err := io.ReadAll(resp.Body)
//		if err != nil {
//			return nil, fmt.Errorf("read Body Error:%v", err)
//		}
//		//遍历该地址的每一条交易信息，并存储于transList中
//		_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//			trans := processTransaction(value, addr)
//			transList = append(transList, trans)
//			//将交易记录存储到es中
//			err = c.Insert(constants.ES_TRANSACTION, trans.Hash, trans)
//			if err != nil {
//				log.Println("Insert Transaction Error:", err.Error())
//				return
//			}
//		}, "txs")
//		if err != nil {
//			log.Println("ArrayEach Txs:", err.Error())
//			return nil, err
//		}
//		return transList, nil
//	}
func GetTxListByBtcCom(bitClient *BitClient, c *es.ElasticClient, addr, page string) ([]domain.EsTrans, int64, error) {
	log.Println("GetTxListByBtcCom:", addr)
	var transList []domain.EsTrans
	var pageTotal int64
	//获得url
	url := getUrlByBtcCom(addr, page)
	resp, err := bitClient.SendHTTPRequest(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read Body Error:%v", err)
	}
	//判断请求状态是否成功
	status, err := jsonparser.GetString(body, "status")
	if err != nil {
		return nil, 0, fmt.Errorf("get status Error:%v", err)
	} else if status != "success" {
		return nil, 0, fmt.Errorf("status is %s", status)
	}
	//判断地址所属交易是否为空
	totalCount, err := jsonparser.GetInt(body, "data", "total_count")
	if err != nil {
		return nil, 0, fmt.Errorf("get total_count Error:%v", err)
	}
	if totalCount == 0 {
		log.Println(addr, "total_count is 0")
		return nil, 0, nil
	}
	//page=1，说明是第一次请求
	if page == "1" {
		//获得页数总数
		pageTotal, err = jsonparser.GetInt(body, "data", "page_total")
		if err != nil {
			return nil, 0, fmt.Errorf("get page_total Error:%v", err)
		}
	}
	//遍历该地址的每一条交易信息，并存储于transList中
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		trans := processTransByBtcCom(value, addr)
		transList = append(transList, trans)
		//将交易记录存储到es中
		err = c.Insert(constants.ES_TRANSACTION, trans.Hash, trans)
		if err != nil {
			log.Println("Insert Transaction Error:", err.Error())
			return
		}
	}, "data", "list")
	if err != nil {
		return nil, 0, fmt.Errorf("ArrayEach %s Txs:%s", addr, err.Error())
	}
	return transList, pageTotal, nil
}

// GetAddressInfo 获取指定地址的余额
func GetAddressInfo(bitClient *BitClient, addr string) (float64, error) {
	url := getAddrInfoUrl(addr)
	resp, err := bitClient.SendHTTPRequest(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("GetAddressInfo read Body Error:%v", err)
	}
	//判断请求状态是否成功
	status, err := jsonparser.GetString(body, "status")
	if err != nil {
		return 0, fmt.Errorf("get status Error:%v", err)
	} else if status != "success" {
		return 0, fmt.Errorf("status is %s", status)
	}
	//判断地址所属交易是否为空
	balance, err := jsonparser.GetInt(body, "data", "balance")
	if err != nil {
		return 0, fmt.Errorf("get balance Error:%v", err)
	}
	balanceFloat, _ := ConvertSatoshiToBTC(big.NewInt(balance))
	return balanceFloat, nil

}
func processTransByBtcCom(value []byte, addr string) domain.EsTrans {
	var inputTransList []domain.InputsTrans //inputs信息
	var outTransList []domain.OutTrans      //out信息
	transHash, err := jsonparser.GetString(value, "hash")
	utils.IsError(err, "Fail get hash")
	transSize, err := jsonparser.GetInt(value, "size")
	utils.IsError(err, "Fail get size")
	inputCount, err := jsonparser.GetInt(value, "inputs_count")
	utils.IsError(err, "Fail get inputs_count")
	outputCount, err := jsonparser.GetInt(value, "outputs_count")
	utils.IsError(err, "Fail get outputs_count")
	inputsValue, err := jsonparser.GetInt(value, "inputs_value")
	utils.IsError(err, "Fail get inputs_value")
	outsValue, err := jsonparser.GetInt(value, "outputs_value")
	utils.IsError(err, "Fail get outputs_value")
	transWeight, err := jsonparser.GetInt(value, "weight")
	utils.IsError(err, "Fail get weight")
	transFee, _, _, err := jsonparser.Get(value, "fee")
	utils.IsError(err, "Fail get fee")
	transLockTime, err := jsonparser.GetInt(value, "lock_time")
	utils.IsError(err, "Fail get lock_time")
	//transTxIndex, _, _, err := jsonparser.Get(value, "tx_index")
	//utils.IsError(err, "Fail get tx_index")
	transConfirm, err := jsonparser.GetInt(value, "confirmations")
	utils.IsError(err, "Fail get confirmations")
	transDoubleSpend, err := jsonparser.GetBoolean(value, "is_double_spend")
	utils.IsError(err, "Fail get double_spend")
	transTime, err := jsonparser.GetInt(value, "block_time")
	utils.IsError(err, "Fail get time")
	transBlockIndex, err := jsonparser.GetInt(value, "block_height")
	utils.IsError(err, "Fail get block_height")
	//获得out信息
	jsonparser.ArrayEach(value, func(outValue []byte, dataType jsonparser.ValueType, offset int, err error) {
		//将单个out信息存储于out list中
		outTransList = append(outTransList, processOutByBtcCom(outValue))
	}, "outputs")
	//获取inputs信息
	jsonparser.ArrayEach(value, func(inputValue []byte, dataType jsonparser.ValueType, offset int, err error) {
		//将单个inputs信息存储于inputs list中
		inputTransList = append(inputTransList, processInputsByBtcCom(inputValue))
	}, "inputs")
	inputValueFloat, inputValueString := ConvertSatoshiToBTC(big.NewInt(inputsValue))
	outputValueFloat, _ := ConvertSatoshiToBTC(big.NewInt(outsValue))
	return domain.EsTrans{
		Hash:          transHash,
		Address:       addr,
		Size:          transSize,
		Weight:        transWeight,
		GasUsed:       string(transFee),
		LockTime:      transLockTime,
		Confirmations: strconv.Itoa(int(transConfirm)),
		DoubleSpend:   transDoubleSpend,
		Time:          utils.UnixToTime(transTime),
		BlockHeight:   strconv.Itoa(int(transBlockIndex)),
		Inputs:        inputTransList,
		Out:           outTransList,
		InputCount:    inputCount,
		OutputCount:   outputCount,
		InputValue:    inputValueFloat,
		OutputValue:   outputValueFloat,
		Chain:         constants.CHAIN_BTC,
		Value:         inputValueFloat,
		ValueText:     inputValueString,
	}
}

// processInputs 处理 inputs 部分的代码块
//
//	func processInputs(inputValue []byte) domain.InputsTrans {
//		inputSequence, err := jsonparser.GetInt(inputValue, "sequence")
//		utils.IsError(err, "Fail get input sequence")
//		inputWitness, err := jsonparser.GetString(inputValue, "witness")
//		utils.IsError(err, "Fail get input witness")
//		inputScript, err := jsonparser.GetString(inputValue, "script")
//		utils.IsError(err, "Fail get input script")
//		inputAddr, err := jsonparser.GetString(inputValue, "prev_out", "addr")
//		utils.IsErrorFloat(err, "Fail get input addr")
//		inputSpent, err := jsonparser.GetBoolean(inputValue, "prev_out", "spent")
//		utils.IsError(err, "Fail get input spent")
//		inputTxIndex, _, _, err := jsonparser.Get(inputValue, "prev_out", "tx_index")
//		utils.IsError(err, "Fail get input tx_index")
//		inputValueVal, _, _, err := jsonparser.Get(inputValue, "prev_out", "value")
//		utils.IsErrorFloat(err, "Fail get input value")
//		var result big.Int
//		return domain.InputsTrans{
//			Sequence: inputSequence,
//			Witness:  inputWitness,
//			Script:   inputScript,
//			Addr:     inputAddr,
//			Spent:    inputSpent,
//			TxIndex:  string(inputTxIndex),
//			Value:    result.SetBytes(inputValueVal).Int64(),
//		}
//		//将单个inputs信息存储于inputs list中
//		//inputTransList = append(inputTransList, inputsTrans)
//	}
func processInputsByBtcCom(inputValue []byte) domain.InputsTrans {
	inputSequence, err := jsonparser.GetInt(inputValue, "sequence")
	utils.IsError(err, "Fail get input sequence")
	inputWitness, err := jsonparser.GetString(inputValue, "prev_type")
	utils.IsError(err, "Fail get input witness")
	inputScript, err := jsonparser.GetString(inputValue, "prev_tx_hash")
	utils.IsError(err, "Fail get input script")
	inputTxIndex, err := jsonparser.GetInt(inputValue, "prev_position")
	utils.IsError(err, "Fail get input tx_index")
	inputValueVal, err := jsonparser.GetInt(inputValue, "prev_value")
	utils.IsErrorFloat(err, "Fail get input value")
	inputAddr, _, _, err := jsonparser.Get(inputValue, "prev_addresses")
	utils.IsErrorFloat(err, "Fail get input addr")
	var addrList []string
	err = json.Unmarshal(inputAddr, &addrList)
	if err != nil {
		log.Fatal("fail unmarshal:", err) // handle error
	}
	valueFloat, valueText := ConvertSatoshiToBTC(big.NewInt(inputValueVal))
	return domain.InputsTrans{
		Sequence:  inputSequence,
		Witness:   inputWitness,
		Script:    inputScript,
		Addr:      strings.Join(addrList, ","),
		TxIndex:   strconv.Itoa(int(inputTxIndex)),
		Value:     valueFloat,
		ValueText: valueText,
	}
	//将单个inputs信息存储于inputs list中
	//inputTransList = append(inputTransList, inputsTrans)
}

// processOut 处理 out 部分的代码块
//
//	func processOut(outValue []byte) domain.OutTrans {
//		outSpent, err := jsonparser.GetBoolean(outValue, "spent")
//		utils.IsError(err, "Fail get spent")
//		outValueVal, _, _, err := jsonparser.Get(outValue, "value")
//		utils.IsErrorFloat(err, "Fail get value")
//		outN, err := jsonparser.GetInt(outValue, "n")
//		utils.IsError(err, "Fail get n")
//		outTxIndex, _, _, err := jsonparser.Get(outValue, "tx_index")
//		utils.IsError(err, "Fail get tx_index")
//		outScript, err := jsonparser.GetString(outValue, "script")
//		utils.IsError(err, "Fail get script")
//		outAddr, err := jsonparser.GetString(outValue, "addr")
//		utils.IsErrorFloat(err, "Fail get addr")
//		var result big.Int
//		return domain.OutTrans{
//			Spent:   outSpent,
//			Value:   result.SetBytes(outValueVal).Int64(),
//			TxIndex: string(outTxIndex),
//			Script:  outScript,
//			Addr:    outAddr,
//			N:       outN,
//		}
//	}
func processOutByBtcCom(outValue []byte) domain.OutTrans {
	outValueVal, err := jsonparser.GetInt(outValue, "value")
	utils.IsErrorFloat(err, "Fail get value")
	outN, err := jsonparser.GetInt(outValue, "spent_by_tx_position")
	utils.IsError(err, "Fail get n")
	outScript, err := jsonparser.GetString(outValue, "spent_by_tx")
	utils.IsError(err, "Fail get spent_by_tx")
	outAddr, _, _, err := jsonparser.Get(outValue, "addresses")
	utils.IsErrorFloat(err, "Fail get addr")
	var addressList []string
	err = json.Unmarshal(outAddr, &addressList)
	if err != nil {
		log.Fatal("fail unmarshal:", err)
	}
	outValueFloat, outValueText := ConvertSatoshiToBTC(big.NewInt(outValueVal))
	return domain.OutTrans{
		Value:     outValueFloat,
		ValueText: outValueText,
		Addr:      strings.Join(addressList, ","),
		N:         outN,
		Script:    outScript,
	}
}

//func processTransaction(value []byte, addr string) domain.EsTrans {
//	var inputTransList []domain.InputsTrans //inputs信息
//	var outTransList []domain.OutTrans      //out信息
//	transHash, err := jsonparser.GetString(value, "hash")
//	utils.IsError(err, "Fail get hash")
//	transSize, err := jsonparser.GetInt(value, "size")
//	utils.IsError(err, "Fail get size")
//	transWeight, err := jsonparser.GetInt(value, "weight")
//	utils.IsError(err, "Fail get weight")
//	transFee, _, _, err := jsonparser.Get(value, "fee")
//	utils.IsError(err, "Fail get fee")
//	transLockTime, err := jsonparser.GetInt(value, "lock_time")
//	utils.IsError(err, "Fail get lock_time")
//	transTxIndex, _, _, err := jsonparser.Get(value, "tx_index")
//	utils.IsError(err, "Fail get tx_index")
//	transDoubleSpend, err := jsonparser.GetBoolean(value, "double_spend")
//	utils.IsError(err, "Fail get double_spend")
//	transTime, err := jsonparser.GetInt(value, "time")
//	utils.IsError(err, "Fail get time")
//	transBlockIndex, err := jsonparser.GetInt(value, "block_height")
//	utils.IsError(err, "Fail get block_height")
//	//获得out信息
//	jsonparser.ArrayEach(value, func(outValue []byte, dataType jsonparser.ValueType, offset int, err error) {
//		//将单个out信息存储于out list中
//		outTransList = append(outTransList, processOut(outValue))
//	}, "out")
//	//获取inputs信息
//	jsonparser.ArrayEach(value, func(inputValue []byte, dataType jsonparser.ValueType, offset int, err error) {
//		//将单个inputs信息存储于inputs list中
//		inputTransList = append(inputTransList, processInputs(inputValue))
//	}, "inputs")
//	return domain.EsTrans{
//		Hash:        transHash,
//		Address:     addr,
//		Size:        transSize,
//		Weight:      transWeight,
//		GasUsed:     string(transFee),
//		LockTime:    transLockTime,
//		TxIndex:     string(transTxIndex),
//		DoubleSpend: transDoubleSpend,
//		Time:        utils.UnixToTime(transTime),
//		BlockHeight: strconv.Itoa(int(transBlockIndex)),
//		Inputs:      inputTransList,
//		Out:         outTransList,
//	}
//}

// GetTxAndSublistByAddr 查询指定的地址的所有交易信息，并将该地址的交易信息以及地址信息存于es中
//func GetTxAndSublistByAddr(addr string) ([]string, error) {
//	var err error
//	var subList []string
//	//用于去除重复数据
//	temp := map[string]struct{}{}
//	//获得url
//	url := getUrlToBtcTrans(addr)
//	//发送http请求-根据url获取到指定账户的所有交易信息
//	resp, err := utils.SendHTTPRequest(url, constants.HTTP_GET, nil)
//	if err != nil {
//		log.Println("Request Error:", err.Error())
//		return nil, err
//	}
//	defer resp.Body.Close()
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal("Io Read Error:", err)
//		return nil, err
//	}
//	//遍历该地址的每一条交易信息，获取转出的交易信息，并将转出对象地址存储
//	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//		var inputTransList []domain.InputsTrans //inputs信息
//		var outTransList []domain.OutTrans      //out信息
//		var isTransOut bool                     //是否为转出交易
//		transHash, err := jsonparser.GetString(value, "hash")
//		IsErrorFloat(err, "Fail get hash")
//		transSize, err := jsonparser.GetInt(value, "size")
//		IsError(err, "Fail get size")
//		transWeight, err := jsonparser.GetInt(value, "weight")
//		IsError(err, "Fail get weight")
//		transFee, _, _, err := jsonparser.Get(value, "fee")
//		IsError(err, "Fail get fee")
//		transLockTime, err := jsonparser.GetInt(value, "lock_time")
//		IsError(err, "Fail get lock_time")
//		transTxIndex, _, _, err := jsonparser.Get(value, "tx_index")
//		IsError(err, "Fail get tx_index")
//		transDoubleSpend, err := jsonparser.GetBoolean(value, "double_spend")
//		IsError(err, "Fail get double_spend")
//		transTime, err := jsonparser.GetInt(value, "time")
//		IsError(err, "Fail get time")
//		transBlockIndex, err := jsonparser.GetInt(value, "block_height")
//		IsError(err, "Fail get block_height")
//		//获得out信息
//		jsonparser.ArrayEach(value, func(outValue []byte, dataType jsonparser.ValueType, offset int, err error) {
//			outSpent, err := jsonparser.GetBoolean(outValue, "spent")
//			IsError(err, "Fail get spent")
//			//outValueVal, _, _, err := jsonparser.Get(outValue, "value")
//			IsError(err, "Fail get value")
//			outN, err := jsonparser.GetInt(outValue, "n")
//			IsError(err, "Fail get n")
//			outTxIndex, _, _, err := jsonparser.Get(outValue, "tx_index")
//			IsError(err, "Fail get tx_index")
//			outScript, err := jsonparser.GetString(outValue, "script")
//			IsError(err, "Fail get script")
//			outAddr, err := jsonparser.GetString(outValue, "addr")
//			IsError(err, "Fail get addr")
//			outTransList = append(outTransList, domain.OutTrans{
//				Spent: outSpent,
//				//Value:   string(outValueVal),
//				TxIndex: string(outTxIndex),
//				Script:  outScript,
//				Addr:    outAddr,
//				N:       outN,
//			})
//			//若outAddr等于addr,说明该交易为转出交易
//			if strings.EqualFold(outAddr, addr) { //不区分大小写判断字符串是否相等
//				isTransOut = true
//			}
//		}, "out")
//		//获取inputs信息
//		jsonparser.ArrayEach(value, func(inputValue []byte, dataType jsonparser.ValueType, offset int, err error) {
//			inputSequence, err := jsonparser.GetInt(inputValue, "sequence")
//			IsError(err, "Fail get input sequence")
//			inputWitness, err := jsonparser.GetString(inputValue, "witness")
//			IsError(err, "Fail get input witness")
//			inputScript, err := jsonparser.GetString(inputValue, "script")
//			IsError(err, "Fail get input script")
//			inputAddr, err := jsonparser.GetString(inputValue, "prev_out", "addr")
//			IsError(err, "Fail get input addr")
//			inputSpent, err := jsonparser.GetBoolean(inputValue, "prev_out", "spent")
//			IsError(err, "Fail get input spent")
//			inputTxIndex, _, _, err := jsonparser.Get(inputValue, "prev_out", "tx_index")
//			IsError(err, "Fail get input tx_index")
//			//inputValueVal, _, _, err := jsonparser.Get(inputValue, "prev_out", "value")
//			IsError(err, "Fail get input value")
//			inputsTrans := domain.InputsTrans{
//				Sequence: inputSequence,
//				Witness:  inputWitness,
//				Script:   inputScript,
//				Addr:     inputAddr,
//				Spent:    inputSpent,
//				TxIndex:  string(inputTxIndex),
//				//Value:    string(inputValueVal),
//			}
//			//将单个inputs信息存储于inputs list中
//			inputTransList = append(inputTransList, inputsTrans)
//			//判断该交易是否为转出交易。若该交易为转出交易，将其转入地址名单存入风险名单信息中
//			if isTransOut {
//				//去除重复的子名单信息,更新地址来源
//				if _, ok := temp[inputAddr]; !ok {
//					temp[inputAddr] = struct{}{}
//					subList = append(subList, inputAddr)
//					es.UpdateOrAddSourceOnTrans(inputAddr, addr)
//				}
//			}
//		}, "inputs")
//		transInfo := domain.EsTrans{
//			Hash:    transHash,
//			Address: addr,
//			//TxType:      constants.TRANS_TYPE_NORMAL,
//			Size:        transSize,
//			Weight:      transWeight,
//			GasUsed:     string(transFee),
//			LockTime:    transLockTime,
//			TxIndex:     string(transTxIndex),
//			DoubleSpend: transDoubleSpend,
//			Time:        transTime,
//			BlockHeight: strconv.Itoa(int(transBlockIndex)),
//			Inputs:      inputTransList,
//			Out:         outTransList,
//		}
//		err = es.Insert(param.ADDRESS_TRANS_LIST, transHash, transInfo)
//		if err != nil {
//			log.Println("Insert trans_list error:", err.Error())
//			return
//		}
//	}, "txs")
//	if err != nil {
//		log.Println("ArrayEach Txs:", err.Error())
//		return nil, err
//	}
//	return subList, err
//}

// AddTransAndSourceByAddr 根据指定地址添加（or 更新）指定层级的子风险名单以及地址信息
//func AddTransAndSourceByAddr(addr string, level int) error {
//	arrList, err := GetSublistByAddr(addr)
//	if err != nil {
//		log.Println("GetSublistByAddr:fail get sublist")
//		return err
//	}
//	resultList, _, _ := GetListByListOnBtc(arrList, constants.INIT_LEVEL+1, level)
//	if resultList == nil || len(resultList) == 0 {
//		log.Println("GetListByList:fail get next addrList")
//		return errors.New("GetListByList:fail get next addrList")
//	}
//	return nil
//}

// GetListByListOnBtc 获取next层所有子名单信息
//func GetListByListOnBtc(list []string, n, max int) ([]string, int, int) {
//	fmt.Println("第", n, "层")
//	//按层次存储子名单信息
//	if n <= max {
//		var allList []string
//		for _, v := range list {
//			nextList, err := GetSublistByAddr(v)
//			if err != nil {
//				log.Println("GetSublistByAddr level ", n, " error:"+err.Error())
//				return GetListByListOnBtc(nil, max+1, max)
//			}
//
//			//将查询到的子名单信息汇总
//			allList = append(allList, nextList...)
//		}
//		return GetListByListOnBtc(allList, n+1, max)
//	}
//	return list, n, max
//}

// GetTransOnBtc 根据给定的地址，查询出该地址的所有交易信息,并将交易信息存储于es中
//func GetTransOnBtc(addr string) (domain.TransactionBtc, error) {
//	var err error
//	var trans domain.TransactionBtc
//	//根据指定地址查询交易信息
//	//获得url
//	url := getUrlToBtcTrans(addr)
//	//发送http请求
//	resp, err := utils.CreateClient().Get(url)
//	if err != nil || resp.StatusCode != http.StatusOK {
//		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
//		return domain.TransactionBtc{}, err
//	}
//	defer resp.Body.Close()
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal(err.Error())
//		return domain.TransactionBtc{}, err
//	}
//	//将查询到的数据绑定到结构体中
//	err = json.Unmarshal(
//		body,
//		&trans,
//	)
//	if err != nil {
//		log.Fatal("json unmarshal error:", err.Error())
//		return domain.TransactionBtc{}, err
//	}
//	return trans, nil
//}

// GetAssocAddr 根据给定地址，获取地址交易关联地址名单
//func GetAssocAddr(addr string) ([]string, error) {
//	var err error
//	//关联地址名单
//	var addrList []string
//	//用于去除重复数据
//	temp := map[string]struct{}{}
//	temp[addr] = struct{}{}
//	trans, err := GetTransOnBtc(addr)
//	if err == nil {
//		//遍历地址的交易信息，将地址关联存于addrList中
//		for _, tx := range trans.Txs {
//			for _, out := range tx.Out {
//				if _, ok := temp[out.Addr]; !ok {
//					addrList = append(addrList, out.Addr)
//				}
//			}
//		}
//	}
//	return addrList, err
//}

// 遍历交易的输入信息，判断该交易是否为转出
//
//	func isTransfer(data []byte, addr string) (bool, error) {
//		var errTrans error
//		var isTrue bool
//		//获取输入信息
//		_, errTrans = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//			//获得交易input（输入）信息的地址
//			val, err := jsonparser.GetString(value, "prev_out", "addr")
//			if err != nil {
//				log.Println("jsonparser GetString Error:", err)
//				errTrans = err
//				return
//			}
//			//若输入信息中的地址等于addr，说明该交易为地址的转出交易
//			if addr == val {
//				isTrue = true
//			}
//		}, "inputs")
//		return isTrue, errTrans
//	}
func getUrlByBlockChain(addr string) string {
	return constants.BTC_ADDR_BLOCKCHAIN + addr
}
func getUrlByBtcCom(addr, page string) string {
	return constants.BTC_ADDR + addr + "/tx?pagesize=" + constants.BTC_PAGRSIZE + "&page=" + page
}
func getAddrInfoUrl(addr string) string {
	return constants.BTC_ADDR + addr

}
