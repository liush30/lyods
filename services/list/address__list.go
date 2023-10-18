// Package list 根据指定的源文件路径，查询地址名单，目前有：json、xml、csv格式
package list

import (
	"encoding/csv"
	"github.com/beevik/etree"
	"github.com/buger/jsonparser"
	"github.com/google/uuid"
	"io"
	"log"
	"lyods-adsTool/domain"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"lyods-adsTool/services/bitcoin"
	"os"
	"regexp"
	"time"
)

// GetAddrListByJSONOnBitcoin 根据url获取json格式的地址名单-获取address & chain
func GetAddrListByJSONOnBitcoin(url string, c *es.ElasticClient) error {
	var err error
	//用于去除重复数据
	temp := map[string]struct{}{}
	//创建http请求
	resp, err := utils.SendHTTPRequest(url, constants.HTTP_GET, nil)
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
	//解析读取的数据，根据指定路径road，获取指定的json字段fields
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//获取风险地址
		addrStr, err := jsonparser.GetString(value, "address")
		if err != nil {
			log.Fatal("Fail get address", err.Error())
		}
		//获取所在链
		chain, err := jsonparser.GetString(value, "blockchain")
		if err != nil {
			log.Fatal("Fail get blockchain", err.Error())
		}
		family, err := jsonparser.GetString(value, "family")
		if err != nil {
			log.Fatal("Fail get family", err.Error())
		}
		id := addrStr + chain
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
				Time:       time.Now(),
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
				//log.Printf("添加%s地址至风险名单中\n", addrStr)
				//地址不存在，则新建
				walletAddr := domain.WalletAddr{
					WaAddr:      addrStr,
					WaRiskLevel: constants.INIT_LEVEL,
					WaChain:     chain,
					DsAddr: []domain.AdsDataSource{
						dsAddrInfo,
					},
					IsNeedTrace: true,
				}
				//将风险名单信息存储至风险名单中，id=地址+所在链
				err = c.Insert(constants.ES_ADDRESS, id, walletAddr)
				if err != nil {
					log.Fatal("Fail inset risk address to es", err.Error())
					return
				}
				//查询该地址的交易信息
				_, err := bitcoin.GetTxListOnBTC(c, addrStr)
				if err != nil {
					log.Fatal("Fail get tx list on btc", err.Error())
					return
				}
			}
		}
	}, "result")
	return nil
}

// GetAddrListOnCsv 根据url获得csv格式的地址名单-批量获取address,获取index列的数据
func GetAddrListOnCsv(url string, c *es.ElasticClient) error {
	var err error
	//风险地址名单
	//var addrList []domain.WalletAddr
	//用于去除重复数据
	temp := map[string]struct{}{}
	//创建http请求
	resp, err := utils.SendHTTPRequest(url, constants.HTTP_GET, nil)
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
				continue
			}
			dsAddrInfo := domain.AdsDataSource{
				DsAddr:     url,
				DsType:     constants.DS_UNISWAP,
				Illustrate: category + "-" + riskType,
				Time:       time.Now(),
			}
			//若该地址已经存在，则更新地址来源
			if isExist {
				//log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
				err = c.AddDsAddrSource(addrStr, dsAddrInfo)
				if err != nil {
					log.Println("Fail add data source:", err.Error())
					continue
				}
			} else {
				log.Printf("添加%s地址至风险名单信息\n", addrStr)
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
				err = c.Insert(constants.ES_ADDRESS, addrStr, walletAddr)
				if err != nil {
					log.Println("Fail insert addr_list", err.Error())
					continue
				}
			}
		}
	}
	return nil
}

// GetAddrListOnXmlByElement 根据url获取xml格式的地址名单-根据访问元素，查询所在链以及地址
// 目前链[XBT(BTC), ETH,BSC(bnb),ARB,
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
func GetAddrListOnXmlByElement(path string, c *es.ElasticClient) error {
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
					entityInfo.OrgEstDate = idNumberValue
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
						log.Fatal("字符串匹配失败，MatchString Error:", err.Error())
						return err
					}
					//若匹配失败，将id标签下的信息存储到ID列表中
					if !isCurrency {
						idInfo := domain.ID{
							IDType:         idType,
							IDNumber:       idNumberValue,
							IDCountry:      elementBySdn(constants.ID_COUNTRY, id),
							ExpirationDate: elementBySdn(constants.EXPRI_DATE, id),
							IssueDate:      elementBySdn(constants.ISSUE_DATE, id),
						}
						idInfoList = append(idInfoList, idInfo)
						//数据匹配成功
					} else {
						//获取entity其他信息
						getEntityInfo(sdnEntry, &entityInfo)
						//截取数据：所在链，以及地址,存储到地址风险名单中
						//根据货币名称获取所在链
						chain := getChainBySdn(idType)
						//根据address，查询该地址是否已经存在
						isExist, err := c.IsExistById(constants.ES_ADDRESS, idNumberValue)
						if err != nil {
							log.Printf("IsExistById Error :%s\n，id=%s", err.Error(), idNumberValue)
							continue
						}
						//记录来源信息
						dsAddrInfo := domain.AdsDataSource{
							DsAddr:     constants.DSADDR_SDN,
							DsType:     constants.DS_OFAC,
							Illustrate: "Derived from the OFAC Sanctions List.",
							Time:       time.Now(),
						}
						//若该地址已经存在，则增加地址来源
						if isExist {
							//log.Printf("%s信息已存在,添加该数据来源\n", idNumberValue)
							c.AddDsAddrSource(idNumberValue, dsAddrInfo)
						} else {
							//log.Printf("添加%s地址信息\n", idNumberValue)
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
							}
							err = c.Insert(constants.ES_ADDRESS, idNumberValue, walletAddr)
							if err != nil {
								log.Fatal("Fail insert risk address:", err.Error())
							}
						}
					}
				}
			}
			//get name info
		}
		if isCurrency {
			//存储entity信息
			entityInfo.Email = emailList
			entityInfo.Website = websiteList
			entityInfo.PhoneNumber = phoneList
			entityInfo.IDList = idInfoList
			entityInfo.OtherInfo = otherList
			//存入es中
			err = c.Insert(constants.ES_ENTITY, entityId, &entityInfo)
			if err != nil {
				log.Fatal("Fail insert risk entity:", err.Error())
			}
		}
	}
	return nil
}

// 根据指定的element存储sdn名单上的entity信息
func getEntityInfo(sdnEntry *etree.Element, entityInfo *domain.Entity) {
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
			mainEntryStr := dateOfBirthItem.SelectElement("mainEntry").Text()
			daBi := domain.DateOfBirth{
				DateOfBirth: dateOfBirth,
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
	entityInfo.PlaceOfBirth = plBiList
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
}

// 解析sdn名单中地址对应的链
func getChainBySdn(idType string) string {
	flysnowRegexp, err := regexp.Compile(constants.CONDITION)
	if err != nil {
		log.Fatal("Fail compile address currency:", err.Error())
	}
	params := flysnowRegexp.FindStringSubmatch(idType)
	return pkg.CurrencyToChain[params[len(params)-1]]
}
func elementBySdn(idType string, id *etree.Element) string {
	element := id.SelectElement(idType)
	if element != nil {
		return element.Text()
	}
	return ""
}
