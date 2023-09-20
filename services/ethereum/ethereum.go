package ethereum

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nanmu42/etherscan-api"
	"io"
	"log"
	"lyods-adsTool/db"
	"lyods-adsTool/entities"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

// EthClient ethereum查询风险账号相关操作
type EthClient struct {
	*ethclient.Client
	MClient *http.Client
}
type AddrMap struct {
	abi   []byte //abi信息
	token string //token name
}

// GetNormalTransUrl 返回etherScan中查询指定地址的普通交易信息列表Url
func getNormalTransUrl(addr string) string {
	return constants.API_ETH_TRANS + addr
}

// GetContractAbiUrl 返回etherScan中查询指定合约地址的abiUrl
func getContractAbiUrl(addr string) string {
	return constants.API_ETH_ABI + addr
}

// IsContractAddress 判断地址是否为合约地址-以太坊
func (e *EthClient) IsContractAddress(addressStr string) (bool, error) {
	var address common.Address
	//获取字节码信息
	bytecode, err := e.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Println("Fail get byte code:", err)
		return false, err
	}
	//bytecode>0，说明是合约地址
	return len(bytecode) > 0, nil
}

// getInternalTxn 根据事件信息进行解析，存储到InternalTxn结构体中
// eventName 事件名称
// eventNameToValueByAddress 事件相关参数信息
func (e *EthClient) getInternalTxn(eventName, contractAddress string, eventNameToValueByAddress map[string]map[string]string, blockNumber *big.Int, isErc20 bool, internalTxn *entities.InternalTxn, length int) {
	///如果不存在address类型的参数,则跳过参数解析
	if _, ok := eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]; !ok {
		return

	}
	//仅存在address类型的参数
	if length == 1 && len(eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]) == 1 {
		var changeMoney *big.Int
		var v string
		var err error
		//判断address是from还是to
		//获取交易前后金额差
		if isErc20 {
			for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
				changeMoney, err = e.GetERC20TokenBalanceChange(contractAddress, v, blockNumber)
				if err != nil {
					log.Fatal("Fail GetERC20TokenBalanceChange:", err)
				}
			}
		} else {
			for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
				changeMoney, err = e.GetBalanceChange(blockNumber, v)
				if err != nil {
					log.Fatal("Fail GetERC20TokenBalanceChange:", err)
				}
			}
		}
		//判断changeMoney是否大于0
		if changeMoney.Cmp(big.NewInt(0)) > 0 {
			internalTxn.ToAddr = v
			internalTxn.FromAddr = contractAddress
			internalTxn.ActualAmount = changeMoney
		} else {
			internalTxn.FromAddr = v
			internalTxn.ToAddr = contractAddress
		}
		return
	} else if length == 2 && len(eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]) == 2 {

	}
	//首先对常规相关转账的行为进行判断
	//1.Transfer 事件解析-Transfer(address indexed from, address indexed to, uint256 value)
	if strings.EqualFold(eventName, constants.EVENT_NAME_TRANSFER) {
		fromValue, fromOk := eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS][constants.EVENT_PARAM_FROM]
		toValue, toOk := eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS][constants.EVENT_PARAM_TO]
		//判断是否存在(address indexed from, address indexed to, uint256 value)结构
		if len(eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]) == 2 && len(eventNameToValueByAddress[constants.EVENT_TYPR_UINT]) == 1 {
			//from,to都存在,根据字段名称存储
			if fromOk && toOk {
				internalTxn.FromAddr = fromValue
				internalTxn.ToAddr = toValue
				//to字段存在，from字段不存在
			} else if fromOk && !toOk {
				internalTxn.FromAddr = fromValue
				delete(eventNameToValueByAddress, constants.EVENT_PARAM_FROM)
				for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
					internalTxn.ToAddr = v
				}
				//from字段存在，to字段存在
			} else if !fromOk && toOk {
				internalTxn.ToAddr = toValue
				delete(eventNameToValueByAddress, constants.EVENT_PARAM_TO)
				for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
					internalTxn.FromAddr = v
				}
				//from,to都不存在
			} else {
				//判断地址中哪个是from,哪个是to
				if isErc20 {
					for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
						changeMoney, err := e.GetERC20TokenBalanceChange(contractAddress, v, blockNumber)
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
					for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
						changeMoney, err := e.GetBalanceChange(blockNumber, v)
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
			for _, v := range eventNameToValueByAddress[constants.EVENT_TYPR_UINT] {
				//将v转成big.int类型
				value, _ := new(big.Int).SetString(v, 10)
				internalTxn.Amount = (etherscan.BigInt)(value)
			}
			return
		} else {

		}
		// Deposit (index_topic_1 address sender, uint256 value)
	} else if strings.EqualFold(eventName, constants.EVENT_NAME_DEPOSIT) {

	}

}

// ParseTransReceiptByHash 根据交易hash解析交易中的Receipt,封装成内部交易信息
// 0x0000000000000000000000000000000000000000
func (e *EthClient) ParseTransReceiptByHash(dbClient *sql.DB, hash common.Hash) ([]entities.InternalTxn, error) {
	var err error
	var abiInfo map[common.Address]AddrMap //存储abi信息与token信息
	var resultList []entities.InternalTxn
	//var paramInfo map[string]string
	//根据交易哈希查询交易的receipt信息
	receipt, err := e.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("fail get %v receipt info:%v", hash.String(), err.Error())
	}
	//receipt status 1表示执行成功，0表示执行失败
	if receipt.Status == 1 && len(receipt.Logs) != 0 {
		//交易为普通交易，遍历receipt log
		for _, logInfo := range receipt.Logs {
			var abiByte []byte
			var name, tokenName string
			var paramName []string
			var paramType []string
			var tokenDecimal int
			var internalTxn entities.InternalTxn
			//根据事件参数类型分类，以参数名与参数值的键值对进行存储
			var eventNameToValueByAddress map[string]map[string]string //type-> (name->value)参数类型[参数名]=参数值
			logAddr := logInfo.Address                                 //获取该log所属的合约地址
			//根据logAddr获取abi信息
			//判断该地址信息是否已经存在，若未被记录，则先从数据库中根据合约地址查询abi信息，若未查询到则通过api查询abi信息
			if _, ok := abiInfo[logAddr]; !ok {
				//从数据库中查询合约的abi信息以及token name
				abiByte, tokenName, tokenDecimal, err = db.GetAbiAndTokenByAddr(dbClient, logAddr.String())
				if err != nil {
					return nil, fmt.Errorf("ParseTransReceiptByHash:Fail get info by database->%s", err.Error())
					//若该addr不是白名单地址或是abi信息不存在，通过api查询abi与token name信息并存储
				} else if abiByte == nil {
					contractAbiStr, err := e.GetContractAbi(logAddr.String())
					if err != nil {
						return nil, fmt.Errorf("ParseTransReceiptByHash:Fail get  info by api ->%s", err.Error())
					}
					abiByte = []byte(contractAbiStr)
					abiInfo[logAddr] = AddrMap{
						abi:   abiByte,
						token: tokenName,
					}
					//若该地址存在于白名单中，且存在abi信息,将abi存储至abiInfo
				} else {
					abiInfo[logAddr] = AddrMap{
						abi:   abiByte,
						token: tokenName,
					}
				}
			} else {
				//若abi信息已经存在，则直接根据logAddr获得abi和token name信息
				abiByte = abiInfo[logAddr].abi
				tokenName = abiInfo[logAddr].token
			}
			//判断该交易是否为erc20交易-若token name不为空，且tokenDecimal不为0，则说明该地址为erc20合约地址，则将该交易记为erc20交易
			if tokenName != "" && tokenDecimal != 0 {
				internalTxn.IsErc20 = true
				internalTxn.TokenDecimal = tokenDecimal
				internalTxn.Token = tokenName
				internalTxn.TokenAddress = logAddr.Hex()
			}
			//根据获得的abi信息，生成合约abi对象
			contractABI, err := abi.JSON(bytes.NewReader(abiByte))
			if err != nil {
				return nil, fmt.Errorf("ParseTransReceiptByHash: Fail to create abi->%s", err.Error())
			}
			evenSig := logInfo.Topics[0] //默认topics的第一个参数为事件的签名
			//获取事件签名，并根据事件签名和abi解析出事件
			eventInfo, err := contractABI.EventByID(evenSig)
			if err != nil {
				log.Printf("ParseTransReceiptByHash:Fail to get event info by %s->%v", logInfo.Topics[0].String(), err.Error())
				break
			}
			//var paramInfo = make(map[string]string)
			topicIndex := 1
			//遍历事件信息，解析事件参数名称和类型，以事件名（参数类型 参数名称,....）格式返回
			name = eventInfo.RawName + "("
			for i, input := range eventInfo.Inputs {
				inputType := input.Type.String()
				inputName := input.Name
				if i == len(eventInfo.Inputs)-1 {
					name += inputType + " " + input.Name
				} else {
					name += inputType + " " + input.Name + ","
				}
				//存储事件信息，以参数名称和类型键值对形式存储

				//若事件参数被标记为indexed，参数值存在于topics数组中
				if input.Indexed {
					if inputType == "address" {
						// 格式化address信息
						cleanedHex := strings.TrimPrefix(logInfo.Topics[topicIndex].Hex(), "0x")
						cleanedHex = strings.TrimLeft(cleanedHex, "0")
						cleanedHex = "0x" + cleanedHex
						//paramInfo[input.Name] = cleanedHex
						eventNameToValueByAddress[inputType][inputName] = cleanedHex
					} else if strings.HasPrefix(inputType, "uint") || strings.HasPrefix(inputType, "int") {
						//paramInfo[input.Name] = logInfo.Topics[topicIndex].Big().String()
						eventNameToValueByAddress[inputType][inputName] = logInfo.Topics[topicIndex].Big().String()
					} else {
						//paramInfo[input.Name] = logInfo.Topics[topicIndex].Hex()
						eventNameToValueByAddress[inputType][inputName] = logInfo.Topics[topicIndex].Hex()
					}
					topicIndex++
					//若参数未被标记为indexed，数据将会存在data中,先将事件参数的名字和类型进行存储
				} else {
					paramName = append(paramName, inputName)
					paramType = append(paramType, inputType)
				}
			}
			name += ")"
			//没有被标记为indexed的参数的值则存储于log.data中，需解析出值后存储于键值对中
			//解析data数据
			dataInter, err := contractABI.Unpack(eventInfo.Name, logInfo.Data)
			if err != nil {
				return nil, fmt.Errorf("ParseTransReceiptByHash:Fail unpack data ->%s", err.Error())
			}
			//解析data数据，将参数名与值以键值对的形式存储到paramInfo中
			for i, dataItem := range dataInter {
				dataType := paramType[i]
				dataName := paramName[i]
				// 尝试转换数据项并检查是否成功
				if convertedValue, ok := convertDataItem(dataItem, dataType); ok {
					eventNameToValueByAddress[dataType][dataName] = convertedValue
				} else {
					log.Fatalf("Failed to convert %s type. param name is %s", dataType, dataName)
				}
			}

			//resultList = append(resultList, entities.ReceiptInfo{
			//	Event:     str,
			//	ParamInfo: interReceiptInfo,
			//})
		}
	}
	return resultList, nil
}

// GetRiskListOnContract 查询指定合约地址中的风险地址信息- 合约必须被验证版本
func (e *EthClient) GetRiskListOnContract(addr string) ([]entities.EsTrans, error) {
	var err error
	var transList []entities.EsTrans
	//发送htt请求,获取合约的abi信息
	//abiStr, err := e.GetContractAbi(addr)
	if err != nil {
		return nil, fmt.Errorf("GetRiskListOnContract:Fail get contract abi info -> %s", err.Error())
	}
	//发送http请求，查询到合约的交易列表
	resp, err := e.MClient.Get(getNormalTransUrl(addr))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetRiskListOnContract:Fail request,http status is not ok or happen error:%s", err.Error())
	}
	defer resp.Body.Close()
	//读取数据
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		return nil, fmt.Errorf("GetRiskListOnContract:Io Read error ->%s", err.Error())
	}
	//遍历result的每一条交易信息
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//获取地址的交易哈希
		transHash, err := jsonparser.GetString(value, "hash")
		if err != nil {
			log.Println("Fail get translation hash!")
			return //跳过本元素处理
		}
		//判断交易是否是成功状态，返回0为成功，1为失败。不处理失败的交易（暂定）
		transIsError, err := jsonparser.GetString(value, "isError")
		if err != nil {
			log.Println("Fail get translation isError status in", transHash)
			return
		}
		if transIsError == "1" {
			log.Println(transHash, "is a fail translation")
			return
		}
		//首先判断该交易是否为转出交易:接收方地址是否为addr
		toAddress, err := jsonparser.GetString(value, "to")
		if err != nil {
			log.Println("Fail get translation 'to' address in", transHash)
			return
		}

		if err != nil {
			log.Println("GetRiskListOnContract:Fail parse receipt by hash", transHash, "error is ", err.Error())
			return
		}
		log.Println("Get translation info in", transHash)
		transGas, err := jsonparser.GetString(value, "gasUsed")
		utils.ErrorLog(err, "Fail get gasUsed")
		transContractAddress, err := jsonparser.GetString(value, "contractAddress")
		utils.ErrorLog(err, "Fail get contractAddress")
		transFunName, err := jsonparser.GetString(value, "functionName")
		utils.ErrorLog(err, "Fail get functionName")
		transMethodId, err := jsonparser.GetString(value, "methodId")
		utils.ErrorLog(err, "Fail get methodId")
		transConfirm, err := jsonparser.GetString(value, "confirmations")
		utils.ErrorLog(err, "Fail get confirmations")
		transCumlGasUsed, err := jsonparser.GetString(value, "cumulativeGasUsed")
		utils.ErrorLog(err, "Fail get cumulativeGasUsed")
		transPrice, err := jsonparser.GetString(value, "gasPrice")
		utils.ErrorLog(err, "Fail get gasPrice")
		transTimeStr, err := jsonparser.GetString(value, "timeStamp")
		utils.ErrorLog(err, "Fail get timeStamp")
		transTime, err := strconv.Atoi(transTimeStr)
		utils.ErrorLog(err, "time string to int")
		transBlockNumer, err := jsonparser.GetString(value, "blockNumber")
		utils.ErrorLog(err, "Fail get blockNumber")
		transBlockHash, err := jsonparser.GetString(value, "blockHash")
		utils.ErrorLog(err, "Fail get blockHash")
		transIndex, err := jsonparser.GetString(value, "transactionIndex")
		utils.ErrorLog(err, "Fail get transactionIndex")
		fromAddress, err := jsonparser.GetString(value, "from")
		utils.ErrorLog(err, "Fail get from")
		transValue, err := jsonparser.GetString(value, "value")
		utils.ErrorLog(err, "Fail get value")
		transInput, err := jsonparser.GetString(value, "input")
		utils.ErrorLog(err, "Fail get input")
		transInfo := entities.EsTrans{
			Hash:              transHash,
			Address:           addr,
			TxType:            constants.TRANS_TYPE_NORMAL,
			GasUsed:           transGas,
			IsError:           transIsError,
			ContractAddress:   transContractAddress,
			FunctionName:      transFunName,
			MethodId:          transMethodId,
			Confirmations:     transConfirm,
			CumulativeGasUsed: transCumlGasUsed,
			GasPrice:          transPrice,
			TxIndex:           transIndex,
			Time:              int64(transTime),
			BlockHeight:       transBlockNumer,
			BlockHash:         transBlockHash,
			Inputs: []entities.InputsTrans{
				{
					Witness: transInput,
					Addr:    fromAddress,
					Value:   transValue,
				},
			},
			Out: []entities.OutTrans{
				{
					Addr: toAddress,
				},
			},
		}
		//InternalTx: receiptInfoList,

		//判断该交易是否为合约创建交易，若toAddress为空，说明该交易为合约创建交易，直接将合约创建者列为风险名单
		if toAddress == "" {
			log.Println("this translation is creat contract,hash is", transHash)
			es.UpdateOrAddSourceOnTrans(fromAddress, addr)
			err = es.Insert(param.ADDRESS_TRANS_LIST, transHash, transInfo)
			if err != nil {
				log.Fatal("Fail insert translation info:", err.Error())
				return
			}
		} else {
			txHash := common.HexToHash(transHash)
			//db, err := sql.Open("mysql", "root:lyods@123@tcp(192.168.1.212:3306)/sit_nf_vaw")
			//if err != nil {
			//	log.Fatal(err)
			//}
			//defer db.Close()
			receiptInfoList, err := e.ParseTransReceiptByHash(db, txHash)
			transInfo.InternalTx = receiptInfoList
			transList = append(transList, transInfo)
		}
		//判断该交易类型为out or in
		//交易类型为in，获取其余交易信息
		//if strings.EqualFold(toAddress, addr) {

		//判断转账地址是否为合约地址,如果input值为0x，则说明该转账对象为普通地址,否则为合约对象,（且未发生错误的交易），只处理有实际金额交易的账户，若交易金额为0暂不处理
		//只对以下情况进行处理操作：
		//该交易转出的地址为普通地址,且该交易未发生错误，并存在实际金额交易
		//if transInput == "0x" && transValue != "0" && transIsError == "0" {
		//	//将地址存于子名单中，并将该地址存入到风险名单信息中
		//	if _, ok := temp[toAddress]; !ok {
		//		temp[toAddress] = struct{}{}
		//		subList = append(subList, toAddress)
		//		UpdateOrAddSourceOnTrans(toAddress, addr)
		//	}
		//	//否则为合约地址，且交易未发生错误，存在实际的金额交易
		//} else if transInput != "0x" && transIsError == "0" && transValue != "0" {
		//	//获取该合约地址的交易总次数
		//	//1.
		//
		//}
		//}
	}, "result")
	if err != nil {
		log.Println("ArrayEach result:", err.Error())
		return nil, err
	}
	return transList, nil
}

// GetNonceForTransaction 获得指定交易时，调用者的nonce值
func (e *EthClient) GetNonceForTransaction(txHash string) (uint64, error) {
	// 查询交易
	tx, _, err := e.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return 0, err
	}
	return tx.Nonce(), nil
}
