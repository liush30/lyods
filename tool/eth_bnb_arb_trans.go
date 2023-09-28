// Package tool 查询ethereum交易信息及层级信息
package tool

import (
	"lyods-adsTool/pkg/constants"
)

// GetTransOnEthOrBnbOrArb 根据指定地址查询普通交易信息
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
//// GetUsdtTransOnEth 根据指定地址查询erc20转账信息
//func GetUsdtTransOnEth(addr string) (domain.TransactionEthOrBnb, error) {
//	var err error
//	var trans domain.TransactionEthOrBnb
//	//获取url-根据地址查询以太坊中的USDT交易信息
//	url := getUsdtTransUrl(addr)
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
//		log.Fatal(err.Error())
//		return domain.TransactionEthOrBnb{}, err
//	}
//	//将获取的数据反序列化为结构体信息
//	err = json.Unmarshal(body, &trans)
//	if err != nil {
//		return domain.TransactionEthOrBnb{}, err
//	}
//	return trans, nil
//}

// GetTransByContract 根据合约地址查询交易信息，追踪出合约中可疑交易数据，并将相关交易地址存于风险合约内部
func GetTransByContract(address string) {

}

// 根据chain 类型获取指定地址的普通交易
func getNormalTransUrl(chainType uint, addr string) string {
	if chainType == constants.CHAIN_BNB {
		return getNormalUrlBnb(addr)
	} else if chainType == constants.CHAIN_ARBITRUM {
		return getNormalUrlArb(addr)
	} else if chainType == constants.CHAIN_ETHEREUM {
		return getNormalUrlEth(addr)
	}
	return ""
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
func getNormalUrlEth(addr string) string {
	return constants.API_ETH_TRANS + addr
}

// arbitrum根据地址获取普通交易url
func getNormalUrlArb(addr string) string {
	return constants.API_ARB_TRANS + addr
}

// arbitrum根据地址获得内部交易url
func getInterUrlArb(addr string) string {
	return constants.API_ARB_INTRANS + addr
}
