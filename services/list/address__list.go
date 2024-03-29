// Package list 根据指定的源文件路径，查询地址名单，目前有：json、xml、csv格式
package list

import (
	"encoding/csv"
	"fmt"
	"github.com/beevik/etree"
	"github.com/buger/jsonparser"
	"github.com/google/uuid"
	"io"
	"log"
	"lyods-adsTool/domain"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"lyods-adsTool/services/bitcoin"
	"lyods-adsTool/services/evm"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetAddrListByJSONOnBitcoin 根据url获取json格式的地址名单-获取address & chain
func GetAddrListByJSONOnBitcoin(url string, bitClient *bitcoin.BitClient, c *es.ElasticClient) error {
	var err error
	//用于去除重复数据
	temp := map[string]struct{}{}
	client := utils.CreateClient()
	//创建http请求
	resp, err := utils.SendHTTPRequest(client, url)
	if err != nil {
		log.Println("GetAddrListByJSONOnBitcoin Request Error:", err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Io Read Error:", err)
		return err
	}
	//初始化计数器和时间戳
	//requestCount := 0
	//lastRequestTime := time.Now()
	//创建一个定时器，每秒增加计数器
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	//解析读取的数据，根据指定路径road，获取指定的json字段fields
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//获取风险地址
		addrStr, err := jsonparser.GetString(value, "address")
		if err != nil {
			log.Fatal("Fail get address", err.Error())
		}
		//获取所在链
		//chain, err := jsonparser.GetString(value, "blockchain")
		//if err != nil {
		//	log.Fatal("Fail get blockchain", err.Error())
		//}
		family, err := jsonparser.GetString(value, "family")
		if err != nil {
			log.Fatal("Fail get family", err.Error())
		}
		id := addrStr + "_" + constants.CHAIN_BTC
		//将不重复数据存入到addrList中
		if _, ok := temp[addrStr]; !ok {
			temp[addrStr] = struct{}{}
			//根据address，查询该地址是否已经存在
			isExist, err := c.IsExistById(constants.ES_ADDRESS, id)
			if err != nil {
				log.Fatalf("IsExistById Error :%v\n", err.Error())
				return
			}
			//增加地址来源信息
			dsAddrInfo := domain.AdsDataSource{
				DsAddr:     url,
				DsType:     constants.DS_OPENSANCTIONS,
				Illustrate: "Suspected extortion through" + family,
				Time:       time.Now().Format(time.DateTime),
			}
			//判断该风险地址是否已经存在，若该地址已经存在，则更新地址来源
			if isExist {
				log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
				err = c.AddDsAddrSource(id, dsAddrInfo)
				if err != nil {
					log.Fatal("Fail add data source:", err.Error())
					return
				}
			} else {
				addrBalance, err := bitcoin.GetAddressInfo(bitClient, addrStr)
				if err != nil {
					log.Println("Fail get address info:", err.Error())
					return
				}
				//地址不存在，则新建
				walletAddr := domain.WalletAddr{
					WaAddr:      addrStr,
					WaRiskLevel: constants.INIT_LEVEL,
					WaChain:     constants.CHAIN_BTC,
					DsAddr: []domain.AdsDataSource{
						dsAddrInfo,
					},
					IsNeedTrace: true,
					Balance:     addrBalance,
				}
				//将风险名单信息存储至风险名单中，id=地址+所在链
				err = c.Insert(constants.ES_ADDRESS, strings.ToUpper(id), walletAddr)
				if err != nil {
					log.Println("Fail inset risk address to es", err.Error())
					return
				}
				//如果请求计数超过限制，等待1分钟
				//checkRequestStatus(&requestCount, &lastRequestTime)
				//查询该地址的交易信息
				_, pageTotal, err := bitcoin.GetTxListByBtcCom(bitClient, c, addrStr, constants.BTC_INIT_PAGE)
				if err != nil {
					log.Println("Fail get tx list on btc:", err.Error())
					return
				}
				//bitClient.AddReqCount()
				pageNum := 1
				//若pageTotal>1,则继续查询后续交易信息
				if pageTotal > 1 {
					for int64(pageNum) < pageTotal {
						//checkRequestStatus(&requestCount, &lastRequestTime)
						pageNum++
						_, _, err = bitcoin.GetTxListByBtcCom(bitClient, c, addrStr, strconv.Itoa(pageNum))
						if err != nil {
							log.Fatal("Fail get tx list on btc:", err.Error())
							return
						}
					}
				}
			}
		}
	}, "result")
	if err != nil {
		return fmt.Errorf("fail parse json:%v", err)
	}
	return nil
}

// 检查请求状态

// GetAddrListOnCsv 根据url获得csv格式的地址名单-批量获取address,获取index列的数据
func GetAddrListOnCsv(url string, c *es.ElasticClient, e *evm.EVMClient) error {
	var err error
	//风险地址名单
	//var addrList []domain.WalletAddr
	//用于去除重复数据
	temp := map[string]struct{}{}
	client := utils.CreateClient()
	//创建http请求
	resp, err := utils.SendHTTPRequest(client, url)
	if err != nil {
		log.Println("Request Error:", err.Error())
		return err
	}
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	//跳过文件头
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return err
	}
	for {
		//逐行读取数据
		rec, err := reader.Read()
		//fmt.Println(rec)
		//文件读取完成
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		addrStr := rec[0]  //获取风险地址
		category := rec[1] //获取风险类别
		riskType := rec[7] //获取风险类型
		//将不重复数据存入到addrList中
		if _, ok := temp[addrStr]; !ok {
			temp[addrStr] = struct{}{}
			//根据address，查询该地址是否已经存在
			isExist, err := c.IsExistById(constants.ES_ADDRESS, addrStr)
			if err != nil {
				log.Printf("IsExistById Error :%v\n", err.Error())
				return fmt.Errorf("IsExistById %s Error :%v\n", addrStr, err)
			}
			dsAddrInfo := domain.AdsDataSource{
				DsAddr:     url,
				DsType:     constants.DS_UNISWAP,
				Illustrate: category + "-" + riskType,
				Time:       time.Now().Format(time.DateTime),
			}
			//若该地址已经存在，则更新地址来源
			if isExist {
				//log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
				err = c.AddDsAddrSource(addrStr, dsAddrInfo)
				if err != nil {
					return fmt.Errorf("fail add %s data source:%v", addrStr, err.Error())
				}
			} else {
				//log.Printf("添加%s地址至风险名单信息\n", addrStr)
				//地址不存在，则新建
				walletAddr := domain.WalletAddr{
					WaAddr:      addrStr,
					WaRiskLevel: constants.INIT_LEVEL,
					WaChain:     constants.CHAIN_ETH,
					DsAddr: []domain.AdsDataSource{
						dsAddrInfo,
					},
					IsNeedTrace: true,
				}
				err = c.Insert(constants.ES_ADDRESS, strings.ToUpper(addrStr+"_"+walletAddr.WaChain), walletAddr)
				//if err != nil {
				//	return fmt.Errorf("fail insert %s to addr_list:%v", addrStr, err)
				//}
				//_, block, err := e.GetTxListOnEth(addrStr, constants.ETH_START_BLOCK)
				//if err != nil {
				//	return fmt.Errorf("fail get %s tx list on evm:%v", addrStr, err)
				//}
				//for block == "" {
				//	_, block, err = e.GetTxListOnEth(addrStr, constants.ETH_START_BLOCK)
				//	if err != nil {
				//		return fmt.Errorf("fail get %s tx list on evm:%v", addrStr, err)
				//	}
				//}
			}
		}
	}
	return nil
}

type RClient struct {
	EsClient  *es.ElasticClient
	BtcClient *bitcoin.BitClient
	EthClient *evm.EVMClient
	CbClient  *evm.ChainBaseClient
	BscClient *evm.EVMClient
}

// GetAddrListOnXmlByElement 根据url获取xml格式的地址名单-根据访问元素，查询所在链以及地址
// 目前链[XBT(BTC), ETH,BSC(bnb),
// ,ARB,
// LTC
// XMR
// ETC
// ZEC
// DASH、BTG
// BSV
// BCH(BitcoinCash)
// XVG
// USDT
// XRP
// USDC
// -----("https://www.treasury.gov/ofac/downloads/sdn.xml", `^Digital Currency Address - ([\D]{3,16}$)
func (r *RClient) GetAddrListOnXmlByElement(path string) error {
	var err error
	////发送http请求
	//resp, err := utils.SendHTTPRequest(url, constants.HTTP_GET, nil)
	//if err != nil {
	//	log.Println("Request Error:", err.Error())
	//	return err
	//}
	//defer resp.Body.Close()
	//将文件读取至etree文档
	//doc := etree.NewDocument()
	//if i, err := doc.ReadFrom(resp.Body); err != nil || i <= 0 {
	//	log.Fatal("error:", err, "or read file is nil")
	//}
	// 读取本地 XML 文件
	xmlData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading XML file: %s", err)
	}
	// 创建 etree 文档
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		log.Fatalf("Error parsing XML data: %s", err)
	}
	//访问根元素
	root := doc.SelectElement("sdnList")
	//访问元素 sdnList->sdnEntry->idList->id->idType
	for _, sdnEntry := range root.SelectElements("sdnEntry") {
		var entityInfo domain.Entity
		var isCurrency bool
		var emailList, websiteList, phoneList []string
		var otherList []domain.OtherInfo
		var idInfoList []domain.ID
		//为每一个sdnEntry创建一个id
		entityId := uuid.New().String()
		if idList := sdnEntry.SelectElement("idList"); idList != nil {
			for _, id := range idList.SelectElements("id") {
				//获取idType值
				idType := id.SelectElement("idType").Text()
				//获得idNumber值
				idNumberValue := id.SelectElement("idNumber").Text()
				if idType == constants.GENDER { //获取性别信息
					entityInfo.Gender = idNumberValue
				} else if idType == constants.EMAIL { //获取邮箱信息
					emailList = append(emailList, idNumberValue)
				} else if idType == constants.WEBSITE { //获取网站信息
					websiteList = append(websiteList, idNumberValue)
				} else if idType == constants.PHONE_NUMBER { //获取电话号码信息
					phoneList = append(phoneList, idNumberValue)
				} else if idType == constants.ORGAN_TYPE {
					entityInfo.OrganizationType = idNumberValue //获取组织类型
				} else if idType == constants.ORGAN_DATE { //获取机构成立时间
					dateStr, err := utils.DateChange(idNumberValue)
					if err != nil {
						log.Println("Organization Established Date:invalid date format", err.Error())
						continue
					}
					entityInfo.OrgEstDate = dateStr
				} else if idType == constants.OTHER_INFO1 || idType == constants.OTHER_INFO2 || idType == constants.OTHER_INFO3 { //获取其他备注信息
					otherInfo := domain.OtherInfo{
						Type: idType,
						Info: idNumberValue,
					}
					otherList = append(otherList, otherInfo)
				} else {
					//匹配字符串,判断是否为数据地址信息
					isCurrency, err = regexp.MatchString(constants.CONDITION, idType)
					if err != nil {
						log.Println("MatchString Error:", err.Error())
						continue
					}
					expirationDateStr := elementBySdn(constants.EXPRI_DATE, id)
					issueDateStr := elementBySdn(constants.ISSUE_DATE, id)
					//若匹配失败，将id标签下的信息存储到ID列表中
					if !isCurrency {
						idInfo := domain.ID{
							IDType:    idType,
							IDNumber:  idNumberValue,
							IDCountry: elementBySdn(constants.ID_COUNTRY, id),
						}
						if expirationDateStr != "" {
							expirationDateStr, err = utils.DateChange(expirationDateStr)
							if err != nil {
								log.Println("Expiration Date:invalid date format", err.Error())
								continue
							}
							idInfo.ExpirationDate = expirationDateStr
						}
						if issueDateStr != "" {
							issueDateStr, err = utils.DateChange(issueDateStr)
							if err != nil {
								log.Println("Issue Date:invalid date format", err.Error())
								continue
							}
							idInfo.IssueDate = issueDateStr
						}
						idInfoList = append(idInfoList, idInfo)
						//数据匹配成功
					} else {
						//获取entity其他信息
						err = getEntityInfo(sdnEntry, &entityInfo)
						if err != nil {
							log.Println("getEntityInfo Error:", err.Error())
							continue
						}
						//截取数据：所在链，以及地址,存储到地址风险名单中
						//根据货币名称获取所在链
						chain := getChainBySdn(idType)
						if chain == constants.SDN_CHAIN_BTC {
							chain = constants.CHAIN_BTC
						}
						//根据address，查询该地址是否已经存在
						isExist, err := r.EsClient.IsExistById(constants.ES_ADDRESS, strings.ToUpper(idNumberValue)+"_"+chain)
						if err != nil {
							log.Printf("IsExistById Error :%s\n，id=%s", err.Error(), idNumberValue)
							continue
						}
						//记录来源信息
						dsAddrInfo := domain.AdsDataSource{
							DsAddr:     constants.DSADDR_SDN,
							DsType:     constants.DS_OFAC,
							Illustrate: "Derived from the OFAC Sanctions List.",
							Time:       time.Now().Format(time.DateTime),
						}
						//若该地址已经存在，则增加地址来源
						if isExist {
							log.Printf("%s信息已存在,添加该数据来源\n", idNumberValue)
							r.EsClient.AddDsAddrSource(idNumberValue+"_"+chain, dsAddrInfo)
						} else {
							//地址不存在，则新建风险名单信息,并存储到es中
							walletAddr := domain.WalletAddr{
								WaAddr:      idNumberValue,
								EntityId:    entityId,
								WaRiskLevel: constants.INIT_LEVEL,
								WaChain:     chain,
								DsAddr: []domain.AdsDataSource{
									dsAddrInfo,
								},
								IsNeedTrace: true,
								//Balance:     addrBalance,
							}
							walletAddr.AddressId = strings.ToUpper(idNumberValue) + "_" + chain
							if chain == constants.CHAIN_BTC {
								//查询该地址的交易信息
								_, addrBalance, err := r.BtcClient.GetTxListByBlockChain(r.EsClient, idNumberValue)
								if err != nil {
									return fmt.Errorf("fail get tx list on btc:%v", err)
								}
								//-------------------------------暂时不处理-------------------------------------------
								//pageNum := 1
								////若pageTotal>1,则继续查询后续交易信息
								//if pageTotal > 1 {
								//	for int64(pageNum) < pageTotal {
								//		//checkRequestStatus(&requestCount, &lastRequestTime)
								//		pageNum++
								//		_, _, err = bitcoin.GetTxListByBtcCom(bitClient, c, idNumberValue, strconv.Itoa(pageNum))
								//		if err != nil {
								//			log.Fatal("Fail get tx list on btc:", err.Error())
								//			return fmt.Errorf("fail get tx list on btc:%v and page is %d", err, pageNum)
								//		}
								//	}
								//}
								//--------------------------------------------------------------------------
								//获取账户的余额
								//addrBalance, err := bitcoin.GetAddressInfo(r.BtcClient, idNumberValue)
								//if err != nil {
								//	return fmt.Errorf("fail get bitcoin address balance info:%v", err)
								//}
								walletAddr.Balance = addrBalance
							} else if chain == constants.SDN_CHAIN_ETH {
								_, _, err := r.EthClient.GetTxList(r.EsClient, r.CbClient, idNumberValue, "0")
								if err != nil {
									return fmt.Errorf("fail get tx list on evm:%v", err)
								}
								//获取账户余额
								addrBalance, err := r.EthClient.GetBalance(idNumberValue, nil)
								if err != nil {
									return fmt.Errorf("fail get ethereum address balance info:%v", err)
								}
								balanceFloat, _, err := evm.WeiToEth(addrBalance)
								if err != nil {
									return fmt.Errorf("fail convert wei to evm:%v", err)
								}
								walletAddr.Balance = balanceFloat
							}
							walletAddr.WaRiskLevel = constants.INIT_LEVEL
							walletAddr.LevelTime = time.Now().Format(time.DateTime)
							walletAddr.RiskChgHistory = []domain.RiskChangeRecord{
								{
									RiskLevel:    constants.INIT_LEVEL,
									Description:  "OFAC SDN LIST INIT",
									DateOfChange: time.Now().Format(time.DateTime),
								},
							}
							err = r.EsClient.Insert(constants.ES_ADDRESS, strings.ToUpper(walletAddr.AddressId), walletAddr)
							if err != nil {
								return fmt.Errorf("fail insert %s to addr_list:%v", idNumberValue+"_"+walletAddr.WaChain, err)
							}
							time.Sleep(time.Second * 2)
						}
					}
				}
			}
			//get name info
		}
		if isCurrency {
			//存储entity信息
			entityInfo.EntityId = entityId
			entityInfo.Email = emailList
			entityInfo.Website = websiteList
			entityInfo.PhoneNumber = phoneList
			entityInfo.IDList = idInfoList
			entityInfo.OtherInfo = otherList
			entityInfo.RiskLevel = constants.INIT_LEVEL
			entityInfo.LevelTime = time.Now().Format(time.DateTime)
			//存储风险 变更记录
			entityInfo.RiskChgHistory = []domain.RiskChangeRecord{
				{
					DateOfChange: time.Now().Format(time.DateTime),
					RiskLevel:    constants.INIT_LEVEL,
					Description:  "OFAC SDN LIST INIT",
				},
			}
			//存入es中
			err = r.EsClient.Insert(constants.ES_ENTITY, entityId, &entityInfo)
			if err != nil {
				log.Fatal("Fail insert risk entity:", err.Error())
			}
		}
	}
	return nil
}

// 根据指定的element存储sdn名单上的entity信息
func getEntityInfo(sdnEntry *etree.Element, entityInfo *domain.Entity) error {
	var akaName, otherAddress []string
	var addrInfoList []domain.AddressList   //地址列表
	var daBiList []domain.DateOfBirth       //出生日期列表
	var plBiList []domain.PlaceOfBirth      //出生地址列表
	var ntList, czList []domain.Nationality //国籍列表
	//获取风险名单实体信息
	sdnType := sdnEntry.SelectElement("sdnType").Text()
	//获取实体名字
	if sdnType == "Individual" {
		entityInfo.IsIndividual = true
		entityInfo.Name = elementBySdn("lastName", sdnEntry) + elementBySdn("firstName", sdnEntry)
	} else {
		entityInfo.IsIndividual = false
		entityInfo.Name = elementBySdn("lastName", sdnEntry)
	}
	//获取实体别名列表
	if akaList := sdnEntry.SelectElement("akaList"); akaList != nil {
		for _, aka := range akaList.SelectElements("aka") {
			akaName = append(akaName, elementBySdn("lastName", aka))
		}
	}
	entityInfo.AkaList = akaName
	//获取实体地址列表
	if addressList := sdnEntry.SelectElement("addressList"); addressList != nil {
		for _, address := range addressList.SelectElements("address") {
			country := elementBySdn("country", address)
			stateOrProvince := elementBySdn("stateOrProvince", address)
			city := elementBySdn("city", address)
			//获取address下的子节点信息
			children := address.ChildElements()
			// 遍历子节点并提取其他文本内容
			for _, child := range children {
				// 忽略 uid、city 、stateOrProvince、country节点
				if child.Tag == "uid" || child.Tag == "city" || child.Tag == "stateOrProvince" || child.Tag == "country" {
					continue
				}
				// 存储子节点的文本内容
				otherAddress = append(otherAddress, child.Text())
			}
			//将地址信息存储到地址信息列表中
			addressList := domain.AddressList{
				Country:         country,
				StateOrProvince: stateOrProvince,
				City:            city,
				Other:           otherAddress,
			}
			addrInfoList = append(addrInfoList, addressList)
		}
	}
	entityInfo.AddressList = addrInfoList
	//获取出生列表信息
	if dateOfBirthList := sdnEntry.SelectElement("dateOfBirthList"); dateOfBirthList != nil {
		for _, dateOfBirthItem := range dateOfBirthList.SelectElements("dateOfBirthItem") {
			dateOfBirth := dateOfBirthItem.SelectElement("dateOfBirth").Text()
			formatDate, err := utils.DateChange(dateOfBirth)
			if err != nil {
				return fmt.Errorf("dateOfBirth invalid date format: %s", dateOfBirth)
			}
			mainEntryStr := dateOfBirthItem.SelectElement("mainEntry").Text()
			daBi := domain.DateOfBirth{
				DateOfBirth: formatDate,
				MainEntry:   mainEntryStr == "true",
			}
			daBiList = append(daBiList, daBi)
		}
	}
	entityInfo.DateOfBirthList = daBiList
	//获取出生地列表信息
	if placeOfBirthList := sdnEntry.SelectElement("placeOfBirthList"); placeOfBirthList != nil {
		for _, placeOfBirthItem := range placeOfBirthList.SelectElements("placeOfBirthItem") {
			placeOfBirth := placeOfBirthItem.SelectElement("placeOfBirth").Text()
			mainEntryStr := placeOfBirthItem.SelectElement("mainEntry").Text()
			plBi := domain.PlaceOfBirth{
				PlaceOfBirth: placeOfBirth,
				MainEntry:    mainEntryStr == "true",
			}
			plBiList = append(plBiList, plBi)
		}
	}
	entityInfo.PlaceOfBirthList = plBiList
	// 获取国籍列表信息
	if nationalityList := sdnEntry.SelectElement("nationalityList"); nationalityList != nil {
		for _, nationality := range nationalityList.SelectElements("nationality") {
			country := nationality.SelectElement("country").Text()
			mainEntryStr := nationality.SelectElement("mainEntry").Text()
			nt := domain.Nationality{
				Country:   country,
				MainEntry: mainEntryStr == "true",
			}
			ntList = append(ntList, nt)
		}
	}
	entityInfo.NationalityList = ntList
	//获取公民列表信息
	if citizenshipList := sdnEntry.SelectElement("citizenshipList"); citizenshipList != nil {
		for _, citizenship := range citizenshipList.SelectElements("citizenship") {
			country := citizenship.SelectElement("country").Text()
			mainEntryStr := citizenship.SelectElement("mainEntry").Text()
			cz := domain.Nationality{
				Country:   country,
				MainEntry: mainEntryStr == "true",
			}
			czList = append(czList, cz)
		}
	}
	entityInfo.CitizenshipList = czList
	return nil
}

// 解析sdn名单中地址对应的链
func getChainBySdn(idType string) string {
	flysnowRegexp, err := regexp.Compile(constants.CONDITION)
	if err != nil {
		log.Fatal("Fail compile address currency:", err.Error())
	}
	params := flysnowRegexp.FindStringSubmatch(idType)
	return params[len(params)-1]
}
func elementBySdn(idType string, id *etree.Element) string {
	element := id.SelectElement(idType)
	if element != nil {
		return element.Text()
	}
	return ""
}
