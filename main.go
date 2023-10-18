//package main
//
//import (
//	"context"
//	"fmt"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/ethclient"
//	"log"
//	"regexp"
//)
//func main() {

// url := "https://rest.cryptoapis.io/blockchain-data/xrp-specific/testnet/addresses/rA9bXGJcXvZKaWofrRphdJsBWzhyCfH3z/transactions?context=yourExampleString&limit=50&offset=0&transactionType=payment"
//
// req, _ := http.NewRequest("GET", url, nil)
//
// req.Header.Add("Content-Type", "application/json")
// req.Header.Add("X-API-Key", "my-api-key")
//
// res, _ := http.DefaultClient.Do(req)
//
// defer res.Body.Close()
// body, _ := ioutil.ReadAll(res.Body)
//
// fmt.Println(res)
//
//	cfg := elasticsearch8.Config{
//		Addresses: []string{"https://localhost:9200"},
//		Username:  "elastic",
//		Password:  "123123",
//		Transport: &http.Transport{
//			MaxIdleConnsPerHost:   10,
//			ResponseHeaderTimeout: time.Second,
//			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
//			TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
//		},
//	}
//
// client, err := elasticsearch8.NewClient(cfg)
//
//	if err != nil {
//		// Handle error
//		log.Printf("Elastic 连接失败: %v\n", err.Error())
//	} else {
//
//		log.Println("Elastic 连接成功")
//	}
//
// log.Print(client.Transport.(*elastictransport.Client).URLs())
// insert
//
//	err := tool.Insert("walletdata", "9999id", domain.WalletAddr{WaAddr: "ox999999999999", WaTicker: "T", DsAddr: []domain.AdsDataSource{
//		{
//			DsAddr: "ooooooooooooooooooo",
//		},
//	}})
//
//	if err != nil {
//		fmt.Println("出错啦")
//	}
//
// is, err := tool.GetWalletAddrByte("walletdata", "9999id")
//
//	if err != nil {
//		fmt.Println("出错", err.Error())
//	}
//
// fmt.Println(string(is))
//
//	err = tool.AddDsAddrSource("9999id", domain.AdsDataSource{
//		DsAddr: "bbbbbbbb",
//		DsType: "1",
//		Number: 2,
//	})
//
//	if err != nil {
//		fmt.Println("新增出错啦")
//	}
//
// 更新标记次数
// err := tool.UpdateDsAddrNumber("9999id", "ooooooooooooooooooo", 3)
//
//	if err != nil {
//		fmt.Println("更新标记次数出错啦")
//	}
//
// //更新地址风险层次
// err = tool.UpdateAddrLevel("9999id", 9)
//
//	if err != nil {
//		fmt.Println("更新地址风险层次出错啦")
//	}
//
// is, err := tool.GetWalletAddrById("walletdata", "9999id")
//
//	if err != nil {
//		fmt.Println("出错", err.Error())
//	}
//
// fmt.Println(is)
// level, err := tool.GetAddrLevel("3337878733id")
//
//	if err != nil {
//		log.Println("报错啦")
//	}
//
// fmt.Println(level)
// list, _ := tool.GetAddrListOnXmlByElement("https://www.treasury.gov/ofac/downloads/sdn.xml", `^Digital Currency Address - ([\D]{3,16}$)`, 1)
// for i, _ := range list {
//
// }
// err := tool.IsExistAddressTest()
//
//	if err != nil {
//		fmt.Println("出错啦", err.Error())
//	}
//
// }
package main

import (
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"lyods-adsTool/domain"
	"lyods-adsTool/pre"

	_ "github.com/go-sql-driver/mysql"
	"log"
	"lyods-adsTool/config"
	"lyods-adsTool/pkg/constants"
	"net/http"
	"net/url"
)

// 获取client对象
func createClient() *http.Client {
	//创建代理URL
	proxyURL, err := url.Parse(config.URL_PROXY)
	if err != nil {
		log.Fatal("Error parsing proxy URL: ", err)
		return nil
	}
	// 创建 HTTP Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: transport,
	}
	return client
}
func createEthClient() *ethclient.Client {
	client, err := ethclient.Dial(constants.URL_INFRUA)
	if err != nil {
		log.Fatal("连接失败:", err)
		return nil
	}
	return client
}
func main() {
	//err := list.GetAddrListByJSONOnBitcoin("https://data.opensanctions.org/datasets/20230927/ransomwhere/source.json")
	//if err != nil {
	//	fmt.Println("出错啦", err.Error())
	//}
	//list, err := ethereum.GetTraceTransaction("0x2d25001f57fe2c695771bb3a52a3904a153d0265ec7691bb7fe01e1c748a36a2")
	//if err != nil {
	//	log.Println("出错了")
	//	return
	//}
	//log.Println(list)
	//hash := common.HexToHash("0xe2505d6c07482d284fbbd5a25b346f6dbec03e598ae3b676deb26e7db6e9ed4a")
	//receipt, err := createEthClient().TransactionReceipt(context.Background(), hash)
	//if err != nil {
	//	log.Println("出错啦", err.Error())
	//	return
	//}
	//for _, v := range receipt.Logs[0].Topics {
	//	log.Println(v.Hex())
	//	log.Println("0x" + v.Hex()[len(v.Hex())-40:])
	//	log.Println(strings.EqualFold("0xa42303EE9B2eC1DB7E2a86Ed6C24AF7E49E9e8B9", "0x"+v.Hex()[len(v.Hex())-40:]))
	//	log.Println(utils.ConvertAddressEth(v), "0x"+v.Hex()[len(v.Hex())-40:])
	//}
	err := pre.GetContractAddressToDb(`D:\Code\GoProjec\lyods-adsTool\pre\file\erc20.csv`)
	if err != nil {
		fmt.Println("出错了:", err)
	}
	//                        0x4d29360c2F7Cc54b8d8A28CB4f29343df867748b
	//0x000000000000000000000000a42303ee9b2ec1db7e2a86ed6c24af7e49e9e8b9
	//contractAddress := common.HexToAddress("0x4Fabb145d64652a948d72533023f6E7A623C7C53")
	//code, err := createEthClient().CodeAt(context.Background(), contractAddress, nil)
	//if err != nil {
	//	fmt.Println("Failed to fetch contract code:", err)
	//	return
	//}
	//fmt.Println(hex.EncodeToString(code))
	//fmt.Println("=================================================================================")
	//co1 := common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376")
	//code1, err := createEthClient().CodeAt(context.Background(), co1, nil)
	//if err != nil {
	//	fmt.Println("Failed to fetch contract code:", err)
	//	return
	//}
	//fmt.Println(hex.EncodeToString(code1))
	// 获取交易所在的块号
	//myInt := new(big.Int)
	//myInt.SetString("8720906", 10)
	//increment := new(big.Int)
	//increment.SetInt64(1)
	//// 执行加法操作
	//result := new(big.Int)
	//result.Sub(myInt, increment)
	//// 查询地址在交易前的余额
	//balanceBefore, err := tool.EthClient.BalanceAt(context.Background(), common.HexToAddress("0x1e63Dbc9F73A03900De27a41e57fa0ee5ee5d905"), result)
	//if err != nil {
	//	fmt.Println("get before info error:", err.Error())
	//}
	//
	//fmt.Println(balanceBefore)
	//
	//// 查询地址在交易后的余额
	//balanceAfter, err := tool.EthClient.BalanceAt(context.Background(), common.HexToAddress("0x1e63Dbc9F73A03900De27a41e57fa0ee5ee5d905"), myInt)
	//if err != nil {
	//	fmt.Println("get after info error:", err.Error())
	//}
	//fmt.Println(balanceAfter)

	// 计算余额变化
	//balanceChange := new(big.Int).Sub(balanceAfter, balanceBefore)
	//fmt.Println("变化为：", balanceChange)
	//client := ethereum.EthClient{MClient: createClient(), Client: createEthClient()}
	//iii, err := client.GetBalanceChange(big.NewInt(9162124), "0xefAB18983029d2BA840E34698eFb67fDF8120711")
	//if err != nil {
	//	fmt.Println("Fail GetBalanceChange:", err.Error())
	//}
	//fmt.Println("调用GetBalanceChange得到的结果为：", iii)
	//gg, _ := client.GetERC20TokenBalanceChange("0x5d3a536e4d6dbd6114cc1ead35777bab948e3643", "0x1BdB1783505F626A55a5e3dd3e366df1cd69c055", big.NewInt(9229710))
	//gga, _ := client.GetERC20TokenBalanceChange("0x5d3a536e4d6dbd6114cc1ead35777bab948e3643", "0x22aaA7720ddd5388A3c0A3333430953C68f1849b", big.NewInt(9229710))
	//fmt.Println(gg)
	//fmt.Println(gga)
	//defer client.Close()
	//db, err := sql.Open("mysql", "root:lyods@123@tcp(192.168.1.212:3306)/sit_nf_vaw")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()
	//err = services.AddABIRecords(db, &client)
	//if err != nil {
	//	fmt.Println("出错啦：", err.Error())
	//}
	//re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	//fmt.Printf("is valid: %v\n", re.MatchString("0x323b5d4c32345ced77393b3530b1eed0f346429d")) // is valid: true
	//fmt.Printf("is valid: %v\n", re.MatchString("0xZYXb5d4c32345ced77393b3530b1eed0f346429d")) // is valid: false

	//client, err := ethclient.Dial("https://fabled-radial-isle.discover.quiknode.pro/406bcd9a632fae87343b9baf3ef47d664161810b/")
	//if err != nil {
	//	log.Fatal(err)
	//}
	// 0x Protocol Token (ZRX) smart contract address
	//address := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	//bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is latest block
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//isContract := len(bytecode) > 0
	//
	//fmt.Printf("is contract: %v\n", isContract) // is contract: true

	// a random user account address
	//address := common.HexToAddress("0x901bb9583b24D97e995513C6778dc6888AB6870e")
	//bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is latest block
	//if err != nil {
	//	log.Fatal(err)
	//}
	//client := etherscan.New(etherscan.Mainnet, "[your API key]")

	//isContract, _ := tool.IsContractAddress("0x0836222F2B2B24A3F36f98668Ed8F0B38D1a872f")
	//fmt.Printf("is contract: %v\n", isContract) // is contract: true

	//number, err := tool.EthClient.NonceAt(context.Background(), address, nil)
	//if err != nil {
	//	fmt.Println("c")
	//}

	//	/*------------------------------------------------------------------------------------------------
	//addr := "0x7Db418b5D567A4e0E8c59Ad71BE1FcE48f3E6107" //合约地址 -0x1559FA1b8F28238FD5D76D9f434ad86FD20D1559
	//resp, err := tool.MClient.Get(constants.API_ETH_ABI + "0xa3c1e324ca1ce40db73ed6026c4a177f099b5770")
	//if err != nil || resp.StatusCode != http.StatusOK {
	//	log.Fatal("Http status is :", resp.StatusCode, "Do Error:", err.Error())
	//
	//}
	//defer resp.Body.Close()
	////读取数据
	//body, err := io.ReadAll(resp.Body)
	//if err != nil || body == nil {
	//	log.Fatal("Io Read Error:", err.Error())
	//}
	//abiStr, err := jsonparser.GetString(body, "result")
	//if err != nil {
	//	log.Fatal("Fail get result:", err.Error())
	//}
	////fmt.Println(abiStr)
	//txHash := common.HexToHash("0x36f58d93d566181bd0cb4f256d788a12333b3fb441282020914d49e695838cbb")
	////获取abi
	//receipt, err := tool.EthClient.TransactionReceipt(context.Background(), txHash)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////创建合约abi对象
	//contractABI, err := abi.JSON(bytes.NewReader([]byte(abiStr)))
	//fmt.Println("Transaction Hash:", receipt.TxHash)
	////fmt.Println("Block Hash:", receipt.BlockHash.Hex())
	////fmt.Println("Block Number:", receipt.BlockNumk'k'k'kber)
	////fmt.Println("Gas Used:", receipt.GasUsed)
	////fmt.Println("Status:", receipt.Status) //1表示执行成功，0表示执行失败
	//// 其他交易收据字段...
	//if len(receipt.Logs) == 0 {
	//	fmt.Println("===============为空=================")
	//}
	//event, err := contractABI.EventByID(receipt.Logs[0].Topics[0])
	//fmt.Println(receipt.Logs[0].Topics[0])
	//if err != nil {
	//	fmt.Println("fail get even:", err.Error())
	//}
	//// 获取事件名称
	//var paramInfo = make(map[string]string)
	//var paramName, paramType []string
	//topicIndex := 1
	////fmt.Println("even name = ", eventName)
	////fmt.Println("RawName:", event.RawName)
	////fmt.Println("is Anonymous:", event.Anonymous)
	//fmt.Println("event sig:", event.Sig) //包含根据ABI规范的字符串签名-Withdrawal(address,bytes32,address,uint256)
	////var expectedType []abi.Type
	//////var eventParamInfo
	////
	////// 遍历参数名称和数据类型
	//str := event.RawName + "("
	//for i, input := range event.Inputs {
	//	if i == len(event.Inputs)-1 {
	//		str += input.Type.String() + " " + input.Name
	//	} else {
	//		str += input.Type.String() + " " + input.Name + ","
	//	}
	//	inputType := input.Type.String()
	//	//若事件参数被标记为indexed，参数值存在于topics数组中
	//	if input.Indexed {
	//		if input.Type.String() == "address" {
	//			// 使用 TrimPrefix 去除前导 "0x"
	//			cleanedHex := strings.TrimPrefix(receipt.Logs[0].Topics[topicIndex].Hex(), "0x")
	//			// 使用 TrimLeft 去除开头的连续零
	//			cleanedHex = strings.TrimLeft(cleanedHex, "0")
	//			// 重新加上 "0x"
	//			cleanedHex = "0x" + cleanedHex
	//			paramInfo[input.Name] = cleanedHex
	//		} else if strings.HasPrefix(inputType, "uint") || strings.HasPrefix(inputType, "int") {
	//			paramInfo[input.Name] = receipt.Logs[0].Topics[topicIndex].Big().String()
	//		} else {
	//			paramInfo[input.Name] = receipt.Logs[0].Topics[topicIndex].Hex()
	//		}
	//		topicIndex++
	//		//若参数未被标记为indexed，数据将会存在data中,先将事件参数的名字和类型进行存储
	//	} else {
	//		paramName = append(paramName, input.Name)
	//		paramType = append(paramType, input.Type.String())
	//	}
	//}
	//str += ")"
	//fmt.Println(str)
	//dataInter, err := contractABI.Unpack(event.Name, receipt.Logs[0].Data)
	//if err != nil {
	//	log.Fatal("出错了：", err.Error())
	//}
	//
	////解析data数据，若数据是bytes32,则先将数据转换为hex格式数据
	//for i, dataItem := range dataInter {
	//	if bytes32V, ok := dataItem.([32]byte); ok {
	//		paramInfo[paramName[i]] = common.BytesToHash(bytes32V[:]).Hex()
	//	} else if addressV, addressOk := dataItem.(common.Address); addressOk {
	//		paramInfo[paramName[i]] = addressV.Hex()
	//	} else if uint256V, uint256Ok := dataItem.(*big.Int); uint256Ok {
	//		paramInfo[paramName[i]] = uint256V.String()
	//	} else if boolValue, boolOk := dataItem.(bool); boolOk {
	//		paramInfo[paramName[i]] = strconv.FormatBool(boolValue)
	//	} else if stringValue, stringOk := dataItem.(string); stringOk {
	//		paramInfo[paramName[i]] = stringValue
	//	} else if bytesValue, bytesOK := dataItem.([]byte); bytesOK {
	//		paramInfo[paramName[i]] = hex.EncodeToString(bytesValue)
	//	} else {
	//		log.Fatal("ParseTransReceiptByHash: Type assertion failed dataItem ,----->", i, paramName[i])
	//	}
	//}
	//fmt.Println(paramInfo)
	//fmt.Println("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	//fmt.Println("len=", len(dataInter))
	//for _, va := range dataInter {
	//	v := reflect.ValueOf(va)
	//	t := reflect.TypeOf([32]byte{})
	//	if t == reflect.TypeOf(va) {
	//		vvv, ok := va.([32]byte)
	//		if ok {
	//			fmt.Println(common.BytesToHash(vvv[:]).Hex())
	//			fmt.Println(common.BytesToHash(vvv[:]).String())
	//		}
	//
	//	}
	//	//if v.Type().ConvertibleTo(expectedType[i].GetType()) {
	//	//
	//	//	fmt.Println(v.Convert(expectedType[i].GetType()).String())
	//	//} else {
	//	//	fmt.Printf("Cannot convert value at index %d to the specified type\n", i)
	//	//}
	//	fmt.Println(v)
	//v,ok:=va.(expectedType[i].GetType())
	//fmt.Printf("%v", va)
	//}
	//// 如果交易成功执行，您还可以获取合约事件日志
	//for _, log := range receipt.Logs {
	//	//fmt.Println("***************************************************")
	//	//fmt.Println(log.Address)
	//	//// 解析log的Topics[0]
	//	//fmt.Println(log.Topics[0])
	//
	//	//fmt.Println("***************************************************")
	//
	//	// 解析事件名称
	//	//eventName, err := tool.EthClient.HeaderByNumber(context.Background(), log.BlockNumber)
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//}
	//	// 解析事件参数
	//	//eventArgs := make(map[string]interface{})
	//	//err = abi.UnpackIntoMap(eventArgs, "EventName", log.Data)
	//	//if err != nil {
	//	//	//log.Fatal(err)
	//	//}
	//	//fmt.Println("Event Args:", eventArgs)
	//	//fmt.Println("Event Name:", eventName)
	//	fmt.Println("--------------------------------------------------")
	//	fmt.Println("Log Address:", log.Address.String())
	//	fmt.Println("Log Hash :", log.TxHash)
	//	//fmt.Println("Log Topics:", log.Topics[0].TerminalString())
	//	fmt.Println("--------------------------------------------------")
	//
	//	for _, topice := range log.Topics {
	//		fmt.Println(topice.Hex())
	//	}
	//	fmt.Println("--------------------------------------------------")
	//
	//	fmt.Println("Log Data:", log.Data)
	//	//fmt.Println("Log TxHash:", log.TxHash)
	//	fmt.Println("Log Index:", log.Index)
	//	//fmt.Println("Log BlockHash:", log.BlockHash)
	//	fmt.Println("Log Removed:", log.Removed)
	//	//fmt.Println("Log TxIndex:", log.TxIndex)
	//	//fmt.Println("Log BlockNumber:", log.BlockNumber)
	//	fmt.Println("--------------------------------------------------")
	//	//fmt.Println("--------------------------------------------")
	//	////fmt.Println("Log Data:", log.Data)
	//	//dataInt := new(big.Int)
	//	//dataInt.SetBytes(log.Data)
	//	//dataUint := dataInt.Uint64()
	//	//fmt.Println("Log Data (Uint):", dataUint)
	//
	//	// 其他事件日志字段...
	//}
	//fmt.Println("number:", number)
	//------------------------------------------------------------------------------------------------------------------ */
	//address := common.HexToAddress("0x901bb9583b24D97e995513C6778dc6888AB6870e")
	//sourceCode, err := tool.EthClient.CodeAt(context.Background(), address, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 判断合约是否为空
	//if len(sourceCode) == 0 {
	//	fmt.Println("合约源代码为空")
	//} else {
	//	fmt.Println("合约源代码：", string(sourceCode))
	//}
	//*********************************************MYSQL*********************************************************************
	//db, err := sql.Open("mysql", "root:lyods@123@tcp(192.168.1.212:3306)/sit_nf_vaw")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()
	//addrs := []domain.WhitelistAddr{
	//	{
	//		TWARKey:    uuid.New().String(),
	//		CID:        1,
	//		TWAddr:     "0x4Fabb145d64652a948d72533023f6E7A623C7C53",
	//		TWChain:    "ethereum",
	//		TWType:     "1",
	//		AddType:    "1",
	//		AddrIll:    "BUSD",
	//		TagKey:     "5e1f2747-ba70-4ccf-8fa1-71fd81345b58",
	//		TokenName:  "BUSD",
	//		ProxyAddr:  "0x5864c777697bf9881220328bf2f16908c9afcd7e",
	//		Website:    "http://www.paxos.com/busd",
	//		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
	//		ModifierID: time.Now().Format("2006-01-02 15:04:05"),
	//		Version:    1,
	//	}, {
	//		TWARKey:    uuid.New().String(),
	//		CID:        1,
	//		TWAddr:     "0x0000000000085d4780B73119b644AE5ecd22b376",
	//		TWChain:    "ethereum",
	//		TWType:     "1",
	//		AddType:    "1",
	//		AddrIll:    "TrueUSD: TUSD Token",
	//		TagKey:     "5e1f2747-ba70-4ccf-8fa1-71fd81345b58",
	//		TokenName:  "TUSD",
	//		ProxyAddr:  "0xB650eb28d35691dd1BD481325D40E65273844F9b",
	//		Website:    "https://trueusd.com/",
	//		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
	//		ModifierID: time.Now().Format("2006-01-02 15:04:05"),
	//		Version:    1,
	//	}, {
	//		TWARKey:    uuid.New().String(),
	//		CID:        1,
	//		TWAddr:     "0x853d955aCEf822Db058eb8505911ED77F175b99e",
	//		TWChain:    "ethereum",
	//		TWType:     "1",
	//		AddType:    "1",
	//		AddrIll:    "FRAX",
	//		TagKey:     "5e1f2747-ba70-4ccf-8fa1-71fd81345b58",
	//		TokenName:  "FRAX",
	//		ProxyAddr:  "",
	//		Website:    "https://frax.finance",
	//		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
	//		ModifierID: time.Now().Format("2006-01-02 15:04:05"),
	//		Version:    1,
	//	}, {
	//		TWARKey:    uuid.New().String(),
	//		CID:        1,
	//		TWAddr:     "0x0C10bF8FcB7Bf5412187A595ab97a3609160b5c6",
	//		TWChain:    "ethereum",
	//		TWType:     "1",
	//		AddType:    "1",
	//		AddrIll:    "USDD",
	//		TagKey:     "5e1f2747-ba70-4ccf-8fa1-71fd81345b58",
	//		TokenName:  "USDD",
	//		ProxyAddr:  "",
	//		Website:    "https://usdd.io/",
	//		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
	//		ModifierID: time.Now().Format("2006-01-02 15:04:05"),
	//		Version:    1,
	//	},
	//}
	//// 开始事务
	//tx, err := db.Begin()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// 准备插入语句
	//stmt, err := db.Prepare("INSERT INTO T_WHITELIST_ADDR (TWAR_KEY, CID, TW_ADDR, TW_CHAIN, TW_TYPE, ADD_TYPE, ADDR_ILL, ADDR_SOURCE, TAG_KEY,TOKEN_NAME, ABI,PROXY_ADDR,WEBSITE,CREATOR_ID, CREATE_DATE, MODIFIER_ID, LAST_MODIFY_DATE, VERSION) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?)")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer stmt.Close()
	//
	//// 执行批量插入
	//for _, addr := range addrs {
	//	_, err = stmt.Exec(addr.TWARKey, addr.CID, addr.TWAddr, addr.TWChain, addr.TWType, addr.AddType, addr.AddrIll, addr.AddrSource, addr.TagKey, addr.TokenName, addr.Abi, addr.ProxyAddr, addr.Website, addr.CreatorID, addr.CreateDate, addr.ModifierID, addr.LastModifyDate, addr.Version)
	//	if err != nil {
	//		tx.Rollback()
	//		log.Fatal(err)
	//	}
	//}
	//// 提交事务
	//err = tx.Commit()
	//if err != nil {
	//	log.Fatal(err)
	//}

}

// 添加信息
func addWhitelistAddr(db *sql.DB, addr domain.WhitelistAddr) error {
	stmt, err := db.Prepare("INSERT INTO T_WHITELIST_ADDR (TWAR_KEY, CID, TW_ADDR, TW_CHAIN, TW_TYPE, ADD_TYPE, ADDR_ILL, ADDR_SOURCE, TAG_KEY,TOKEN_NAME, ABI,PROXY_ADDR,WEBSITE,CREATOR_ID, CREATE_DATE, MODIFIER_ID, LAST_MODIFY_DATE, VERSION) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(addr.TWARKey, addr.CID, addr.TWAddr, addr.TWChain, addr.TWType, addr.AddType, addr.AddrIll, addr.AddrSource, addr.TagKey, addr.TokenName, addr.Abi, addr.ProxyAddr, addr.Website, addr.CreatorID, addr.CreateDate, addr.ModifierID, addr.LastModifyDate, addr.Version)
	if err != nil {
		return err
	}
	return nil
}
