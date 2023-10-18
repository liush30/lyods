package ethereum

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"lyods-adsTool/db"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/constants"
	"math/big"
	"strings"
)

// EthClient ethereum查询风险账号相关操作
type EthClient struct {
	*ethclient.Client
}

func CreateEthClient() *ethclient.Client {
	client, err := ethclient.Dial(constants.URL_INFRUA)
	if err != nil {
		log.Fatal("连接失败:", err)
		return nil
	}
	return client
}

// getInternalTxn 根据事件信息进行解析，存储到InternalTxn结构体中
// eventName 事件名称
// eventNameToValueByAddress 事件相关参数信息
func (e *EthClient) getInternalTxn(interParam *InternalTxnParam, internalTxn *domain.InternalTxn) {
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
func (e *EthClient) ParseTransReceiptByHash(dbClient *sql.DB, hash string, contractAbiMap *map[common.Address][]byte) ([]domain.Logs, error) {
	//根据交易哈希查询交易的receipt信息
	receipt, err := e.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt info for %v: %v", hash, err)
	}
	//若log为空，直接返回空
	if len(receipt.Logs) == 0 {
		return nil, nil
	}
	var resultList []domain.Logs //log信息列表
	//查询log信息
	if receipt.Status == 1 { //receipt status 1表示执行成功，0表示执行失败
		resultList, err = getLogsInReceipt(dbClient, receipt.Logs, contractAbiMap)
		if err != nil {
			return nil, fmt.Errorf("fail get logs info:%v", err)
		}
	}
	return resultList, nil
}
func getLogsInReceipt(dbClient *sql.DB, logs []*types.Log, contractAbiMap *map[common.Address][]byte) ([]domain.Logs, error) {
	var resultList []domain.Logs
	//遍历receipt log
	for _, logInfo := range logs {
		abiByte, err := getOrLoadABI(dbClient, logInfo.Address, contractAbiMap)
		if err != nil {
			return nil, fmt.Errorf("failed to get ABI info: %v", err)
		}
		var resultLog domain.Logs
		//若合约未被验证，记录未被解析的log数据，并退出循环，解析下一个log
		if string(abiByte) == constants.ABI_NO {
			resultLog = parseNoVerifyLog(logInfo)
		} else {
			resultLog, err = parseVerifyLog(abiByte, logInfo)
			if err != nil {
				return nil, fmt.Errorf("fail parse log:%v", err.Error())
			}
		}
		resultList = append(resultList, resultLog)
	}
	return resultList, nil
}
func getOrLoadABI(dbClient *sql.DB, logAddr common.Address, contractAbiMap *map[common.Address][]byte) ([]byte, error) {
	if abiByte, ok := (*contractAbiMap)[logAddr]; ok {
		return abiByte, nil
	}

	abiByte, err := loadABIFromDatabase(dbClient, logAddr)
	if err != nil {
		return nil, err
	}

	(*contractAbiMap)[logAddr] = abiByte
	return abiByte, nil
}

// 根据合约地址查询abi信息：1.先从数据库中根据合约地址查询abi信息2.若未查询到则通过api查询abi信息
func loadABIFromDatabase(dbClient *sql.DB, logAddr common.Address) ([]byte, error) {
	// 查询数据库中的合约ABI信息
	abiByte, err := db.GetAbi(dbClient, logAddr.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI info from the database: %v", err)
	}

	if abiByte == nil {
		// 通过地址查询合约的ABI信息
		contractAbiStr, err := GetContractAbiOnEth(logAddr.String())
		if err != nil {
			return nil, fmt.Errorf("failed to get ABI info from the API: %v", err)
		}
		abiByte = []byte(contractAbiStr)
	}

	return abiByte, nil
}

// 获取未验证的合约的log信息
func parseNoVerifyLog(logInfo *types.Log) domain.Logs {
	var topicsInfo []domain.TopicsValStruct
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
	}
}

// 获取已验证合约的log信息
func parseVerifyLog(abiByte []byte, logInfo *types.Log) (domain.Logs, error) {
	//根据获得的abi信息，生成合约abi对象
	contractABI, err := abi.JSON(bytes.NewReader(abiByte))
	if err != nil {
		return domain.Logs{}, fmt.Errorf("ParseTransReceiptByHash: Fail to create abi->%s", err.Error())
	}
	evenSig := logInfo.Topics[0] //默认topics的第一个参数为事件的签名
	//获取事件签名，并根据事件签名和abi解析出事件
	eventInfo, err := contractABI.EventByID(evenSig)
	if err != nil {
		return domain.Logs{}, fmt.Errorf("ParseTransReceiptByHash:Fail to get event info by %s->%v", logInfo.Topics[0].String(), err.Error())
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
			return domain.Logs{}, fmt.Errorf("parseData: %s", err.Error())
		}
	}
	return domain.Logs{
		Address:   logInfo.Address.String(),
		EventInfo: eventInfoStr,
		Topics:    mapToTopicsValStruct(eventNameToValue),
	}, nil
}

// GetRiskListOnContract 查询指定合约地址中的风险地址信息- 合约必须被验证版本
//func (e *EthClient) GetRiskListOnContract(addr string) ([]domain.EsTrans, error) {
//	var err error
//	var transList []domain.EsTrans
//	//发送htt请求,获取合约的abi信息
//	//abiStr, err := e.GetContractAbi(addr)
//	if err != nil {
//		return nil, fmt.Errorf("GetRiskListOnContract:Fail get contract abi info -> %s", err.Error())
//	}
//	//发送http请求，查询到合约的交易列表
//	resp, err := e.MClient.Get(getNormalTransUrl(addr))
//	if err != nil || resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("GetRiskListOnContract:Fail request,http status is not ok or happen error:%s", err.Error())
//	}
//	defer resp.Body.Close()
//	//读取数据
//	body, err := io.ReadAll(resp.Body)
//	if err != nil || body == nil {
//		return nil, fmt.Errorf("GetRiskListOnContract:Io Read error ->%s", err.Error())
//	}
//	//遍历result的每一条交易信息
//	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//		//获取地址的交易哈希
//		transHash, err := jsonparser.GetString(value, "hash")
//		if err != nil {
//			log.Println("Fail get translation hash!")
//			return //跳过本元素处理
//		}
//		//判断交易是否是成功状态，返回0为成功，1为失败。不处理失败的交易（暂定）
//		transIsError, err := jsonparser.GetString(value, "isError")
//		if err != nil {
//			log.Println("Fail get translation isError status in", transHash)
//			return
//		}
//		if transIsError == "1" {
//			log.Println(transHash, "is a fail translation")
//			return
//		}
//		//首先判断该交易是否为转出交易:接收方地址是否为addr
//		toAddress, err := jsonparser.GetString(value, "to")
//		if err != nil {
//			log.Println("Fail get translation 'to' address in", transHash)
//			return
//		}
//
//		if err != nil {
//			log.Println("GetRiskListOnContract:Fail parse receipt by hash", transHash, "error is ", err.Error())
//			return
//		}
//		log.Println("Get translation info in", transHash)
//		transGas, err := jsonparser.GetString(value, "gasUsed")
//		utils.ErrorLog(err, "Fail get gasUsed")
//		transContractAddress, err := jsonparser.GetString(value, "contractAddress")
//		utils.ErrorLog(err, "Fail get contractAddress")
//		transFunName, err := jsonparser.GetString(value, "functionName")
//		utils.ErrorLog(err, "Fail get functionName")
//		transMethodId, err := jsonparser.GetString(value, "methodId")
//		utils.ErrorLog(err, "Fail get methodId")
//		transConfirm, err := jsonparser.GetString(value, "confirmations")
//		utils.ErrorLog(err, "Fail get confirmations")
//		transCumlGasUsed, err := jsonparser.GetString(value, "cumulativeGasUsed")
//		utils.ErrorLog(err, "Fail get cumulativeGasUsed")
//		transPrice, err := jsonparser.GetString(value, "gasPrice")
//		utils.ErrorLog(err, "Fail get gasPrice")
//		transTimeStr, err := jsonparser.GetString(value, "timeStamp")
//		utils.ErrorLog(err, "Fail get timeStamp")
//		transTime, err := strconv.Atoi(transTimeStr)
//		utils.ErrorLog(err, "time string to int")
//		transBlockNumer, err := jsonparser.GetString(value, "blockNumber")
//		utils.ErrorLog(err, "Fail get blockNumber")
//		transBlockHash, err := jsonparser.GetString(value, "blockHash")
//		utils.ErrorLog(err, "Fail get blockHash")
//		transIndex, err := jsonparser.GetString(value, "transactionIndex")
//		utils.ErrorLog(err, "Fail get transactionIndex")
//		fromAddress, err := jsonparser.GetString(value, "from")
//		utils.ErrorLog(err, "Fail get from")
//		transValue, err := jsonparser.GetString(value, "value")
//		utils.ErrorLog(err, "Fail get value")
//		transInput, err := jsonparser.GetString(value, "input")
//		utils.ErrorLog(err, "Fail get input")
//		transInfo := domain.EsTrans{
//			Hash:              transHash,
//			Address:           addr,
//			TxType:            constants.TRANS_TYPE_NORMAL,
//			GasUsed:           transGas,
//			IsError:           transIsError,
//			ContractAddress:   transContractAddress,
//			FunctionName:      transFunName,
//			MethodId:          transMethodId,
//			Confirmations:     transConfirm,
//			CumulativeGasUsed: transCumlGasUsed,
//			GasPrice:          transPrice,
//			TxIndex:           transIndex,
//			Time:              int64(transTime),
//			BlockHeight:       transBlockNumer,
//			BlockHash:         transBlockHash,
//			Inputs: []domain.InputsTrans{
//				{
//					Witness: transInput,
//					Addr:    fromAddress,
//					Value:   transValue,
//				},
//			},
//			Out: []domain.OutTrans{
//				{
//					Addr: toAddress,
//				},
//			},
//		}
//		//InternalTx: receiptInfoList,
//
//		//判断该交易是否为合约创建交易，若toAddress为空，说明该交易为合约创建交易，直接将合约创建者列为风险名单
//		if toAddress == "" {
//			log.Println("this translation is creat contract,hash is", transHash)
//			es.UpdateOrAddSourceOnTrans(fromAddress, addr)
//			err = es.Insert(param.ADDRESS_TRANS_LIST, transHash, transInfo)
//			if err != nil {
//				log.Fatal("Fail insert translation info:", err.Error())
//				return
//			}
//		} else {
//			txHash := common.HexToHash(transHash)
//			//db, err := sql.Open("mysql", "root:lyods@123@tcp(192.168.1.212:3306)/sit_nf_vaw")
//			//if err != nil {
//			//	log.Fatal(err)
//			//}
//			//defer db.Close()
//			receiptInfoList, err := e.ParseTransReceiptByHash(db, txHash)
//			transInfo.InternalTx = receiptInfoList
//			transList = append(transList, transInfo)
//		}
//		//判断该交易类型为out or in
//		//交易类型为in，获取其余交易信息
//		//if strings.EqualFold(toAddress, addr) {
//
//		//判断转账地址是否为合约地址,如果input值为0x，则说明该转账对象为普通地址,否则为合约对象,（且未发生错误的交易），只处理有实际金额交易的账户，若交易金额为0暂不处理
//		//只对以下情况进行处理操作：
//		//该交易转出的地址为普通地址,且该交易未发生错误，并存在实际金额交易
//		//if transInput == "0x" && transValue != "0" && transIsError == "0" {
//		//	//将地址存于子名单中，并将该地址存入到风险名单信息中
//		//	if _, ok := temp[toAddress]; !ok {
//		//		temp[toAddress] = struct{}{}
//		//		subList = append(subList, toAddress)
//		//		UpdateOrAddSourceOnTrans(toAddress, addr)
//		//	}
//		//	//否则为合约地址，且交易未发生错误，存在实际的金额交易
//		//} else if transInput != "0x" && transIsError == "0" && transValue != "0" {
//		//	//获取该合约地址的交易总次数
//		//	//1.
//		//
//		//}
//		//}
//	}, "result")
//	if err != nil {
//		log.Println("ArrayEach result:", err.Error())
//		return nil, err
//	}
//	return transList, nil
//}

// GetNonceForTransaction 获得指定交易时，调用者的nonce值
func (e *EthClient) GetNonceForTransaction(txHash string) (uint64, error) {
	// 查询交易
	tx, _, err := e.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return 0, err
	}
	return tx.Nonce(), nil
}
