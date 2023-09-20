// Package tool 根据指定的源文件路径，查询地址名单，目前有：json、xml、csv格式
package list

import (
	"encoding/csv"
	"github.com/beevik/etree"
	"github.com/buger/jsonparser"
	"github.com/google/uuid"
	"io"
	"log"
	"lyods-adsTool/entities"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/tool"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// GetAddrListByJSONOnBitcoin 根据url获取json格式的地址名单-获取address & chain
func GetAddrListByJSONOnBitcoin(url string) error {
	var err error
	//用于去除重复数据
	temp := map[string]struct{}{}
	//创建http请求
	req, err := http.NewRequest(constants.HTTP_GET, url, nil)
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return err
	}
	//发送http请求
	resp, err := tool.MClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Do Error:", err.Error())
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
			log.Println("Fail get address", err.Error())
		}
		//获取地址
		chain, err := jsonparser.GetString(value, "blockchain")
		if err != nil {
			log.Println("Fail get blockchain", err.Error())
		}
		family, err := jsonparser.GetString(value, "family")
		if err != nil {
			log.Println("Fail get family", err.Error())
		}
		//将不重复数据存入到addrList中
		if _, ok := temp[addrStr]; !ok {
			temp[addrStr] = struct{}{}
			//根据address，查询该地址是否已经存在与
			isExist, err := es.IsExistById(constants.ES_ADDRESS, addrStr)
			if err != nil {
				log.Printf("IsExistById Error :%v\n", err.Error())
				return
			}
			//增加地址来源信息
			dsAddrInfo := entities.AdsDataSource{
				DsAddr:     url,
				DsType:     strconv.Itoa(constants.DS_TYPE_URL),
				Illustrate: "Suspected extortion through" + family,
				Time:       time.Now(),
			}
			//判断该风险地址是否已经存在，若该地址已经存在，则更新地址来源
			if isExist {
				log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
				err = es.AddDsAddrSource(addrStr, dsAddrInfo)
				if err != nil {
					log.Println("Fail add data source:", err.Error())
					return
				}
			} else {
				log.Printf("添加%s地址至风险名单中\n", addrStr)
				//地址不存在，则新建
				walletAddr := entities.WalletAddr{
					WaAddr:      addrStr,
					WaRiskLevel: constants.INIT_LEVEL,
					WaChain:     chain,
					DsAddr: []entities.AdsDataSource{
						dsAddrInfo,
					},
				}
				//将风险名单信息存储至风险名单中
				err = es.Insert(constants.ES_ADDRESS, addrStr, walletAddr)
				if err != nil {
					return
				}
				/************************************
				//增加该地址下子风险名单以及交易信息
				//err = AddTransAndSourceByAddr(addrStr, level)
				//if err != nil {
				//	log.Println("Fail get address transaction and subList:", err.Error())
				//	return
				//}
				*/
			}
		}
	}, "result")
	return nil
}

// GetAddrListOnCsv 根据url获得csv格式的地址名单-批量获取address,获取index列的数据
func GetAddrListOnCsv(url, ticker string, index uint) error {
	var err error
	//风险地址名单
	//var addrList []entities.WalletAddr
	//用于去除重复数据
	temp := map[string]struct{}{}
	//创建http请求
	req, err := http.NewRequest(constants.HTTP_GET, url, nil)
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return err
	}
	//发送http请求
	resp, err := tool.MClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
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
			log.Fatal(err)
			continue
		}
		addrStr := rec[0]  //获取风险地址
		category := rec[1] //获取风险类别
		riskType := rec[7] //获取风险类型
		//将不重复数据存入到addrList中
		if _, ok := temp[addrStr]; !ok {
			temp[addrStr] = struct{}{}
			//根据address，查询该地址是否已经存在
			isExist, err := es.IsExistById(constants.ES_ADDRESS, addrStr)
			if err != nil {
				log.Printf("IsExistById Error :%v\n", err.Error())
				continue
			}
			dsAddrInfo := entities.AdsDataSource{
				DsAddr:     url,
				DsType:     strconv.Itoa(constants.DS_TYPE_URL),
				Illustrate: category + "-" + riskType,
				Time:       time.Now(),
			}
			//若该地址已经存在，则更新地址来源
			if isExist {
				log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
				err = es.AddDsAddrSource(addrStr, dsAddrInfo)
				if err != nil {
					log.Println("Fail add data source:", err.Error())
					continue
				}
			} else {
				log.Printf("添加%s地址至风险名单信息\n", addrStr)
				//地址不存在，则新建
				walletAddr := entities.WalletAddr{
					WaAddr:      addrStr,
					WaRiskLevel: constants.INIT_LEVEL,
					WaChain:     "ethereum",
					DsAddr: []entities.AdsDataSource{
						dsAddrInfo,
					},
				}
				err = es.Insert(constants.ES_ADDRESS, addrStr, walletAddr)
				if err != nil {
					log.Println("Fail insert addr_list", err.Error())
					continue
				}
			}
			/*************************根据风险子名单信息，查询风险子名单信息
			err = AddTranAndSourByAddrOnEth(addrStr, level)
			if err != nil {
				log.Println("Fail get address transaction and subList:", err.Error())
				continue
			}
			*/
		}
	}
	return nil
}

// GetAddrListOnXmlByPath 根据url获取xml格式的地址名单-根据指定的路径查询,查询地址
//func GetAddrListOnXmlByPath(url, ticker, searchPath string, level uint) error {
//var err error
////用于去除重复数据
//temp := map[string]struct{}{}
////发送http请求
//resp, err := MClient.Get(url)
//if err != nil || resp.StatusCode != http.StatusOK {
//	log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
//	return err
//}
//defer resp.Body.Close()
////将文件读取至etree文档
//doc := etree.NewDocument()
//if i, err := doc.ReadFrom(resp.Body); err != nil || i <= 0 {
//	log.Fatal("error:", err, "or read file is nil")
//	return err
//}
////根据指定路径查询
//for _, e := range doc.FindElements(searchPath) {
//	//将不重复数据存入到addrList中
//	if _, ok := temp[e.Text()]; !ok {
//		temp[e.Text()] = struct{}{}
//		addrStr := e.Text() //地址
//		//根据address，查询该地址是否已经存在
//		isExist, err := IsExistById(param.WA_ADDR_INDEX, addrStr)
//		if err != nil {
//			log.Printf("IsExistById Error :%v\n", err.Error())
//			continue
//		}
//		dsAddrInfo := entities.AdsDataSource{
//			DsAddr: url,
//			DsType: param.DS_TYPE_URL,
//			Number: param.INIT_NUMBER,
//		}
//		//若该地址已经存在，则更新地址来源
//		if isExist {
//			log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
//			AddDsAddrSource(addrStr, dsAddrInfo)
//		} else {
//			log.Printf("添加%s地址信息\n", addrStr)
//			//地址不存在，则新建
//			walletAddr := entities.WalletAddr{
//				WaAddr:      addrStr,
//				WaRiskLevel: param.INIT_LEVEL,
//				WaTicker:    ticker,
//				WaChain:     GetChainByTicker(ticker),
//				DsAddr: []entities.AdsDataSource{
//					dsAddrInfo,
//				},
//			}
//			Insert(param.WA_ADDR_INDEX, addrStr, walletAddr)
//		}
//	}
//}
//return nil
//}
//func IsExistAddressTest() error {
//	var err error
//	//用于去除重复数据
//	hash := "81865933c27f00280d115264bcedf4731c614303af79453e524e7c1f88d26446"
//	address1 := "17TkPysx7GyMQTfLTFGq73d8se2Ewp7aQZ"
//	address2 := "3KX2XXxurkbEn3M1zYAW9Dn9QtjUWzGJbP"
//	//发送http请求
//	//resp, err := MClient.Get("https://www.treasury.gov/ofac/downloads/sdn.xml")
//	////fmt.Println(resp)
//	//if err != nil || resp.StatusCode != http.StatusOK {
//	//	//		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
//	//	return err
//	//}
//	bboy, err := os.ReadFile("./sdn.xml")
//	if err != nil {
//		fmt.Println("文件读取失败:", err.Error())
//		return err
//	}
//	//将文件读取至etree文档
//	doc := etree.NewDocument()
//	if i, err := doc.ReadFrom(bytes.NewReader(bboy)); err != nil || i <= 0 {
//		log.Fatal("error:", err, "or read file is nil")
//		return err
//	}
//	path := etree.MustCompilePath("./sdnList/sdnEntry[lastName='HYDRA MARKET']/idList")
//	//根据指定路径查询
//	idList := doc.FindElementPath(path)
//	for _, id := range idList.SelectElements("id") {
//		fmt.Println("1111111111")
//		//获取地址
//		addrStr := id.SelectElement("idNumber").Text()
//		//根据地址查询
//		url := getUrlToBtcTrans(addrStr)
//		sen := RandInt(1, 10)
//		time.Sleep(time.Millisecond * time.Duration(sen))
//		//发送http请求。根据url获取到指定账户的所有交易信息
//		clien := http.DefaultClient
//		respBit, err := clien.Get(url)
//		if err != nil {
//			log.Println("Do Error:", err.Error())
//			return err
//		}
//		body, err := io.ReadAll(respBit.Body)
//		if err != nil {
//			log.Fatal("Io Read Error:", err)
//			return err
//		}
//		var btcInfo entities.TransactionBtc
//		err = json.Unmarshal(body, &btcInfo)
//		if err != nil {
//			fmt.Println("body:", string(body), "\taddrStr:", addrStr)
//			log.Println("json unmarshal error:", err.Error())
//			continue
//		}
//		for _, tx := range btcInfo.Txs {
//			//如果哈希值相等
//			if strings.EqualFold(tx.Hash, hash) {
//				fmt.Println(hash, "是属于", addrStr)
//				return nil
//			}
//			for _, input := range tx.Inputs {
//				if strings.EqualFold(input.PrevOut.Addr, address1) {
//					fmt.Println(address1, "是属于", addrStr, tx.Hash, "中的input address")
//					return nil
//				} else if strings.EqualFold(input.PrevOut.Addr, address2) {
//					fmt.Println(address2, "是属于", addrStr, tx.Hash, "中的input address")
//					return nil
//				}
//			}
//			for _, out := range tx.Out {
//				if strings.EqualFold(out.Addr, address1) {
//					fmt.Println(address1, "是属于", addrStr, tx.Hash, "中的input address")
//					return nil
//				} else if strings.EqualFold(out.Addr, address2) {
//					fmt.Println(address2, "是属于", addrStr, tx.Hash, "中的input address")
//					return nil
//				}
//			}
//		}
//
//	}
//
//	return nil
//}

//func RandInt(min, max int) int {
//	if min >= max || min == 0 || max == 0 {
//		return max
//	}
//	return rand.Intn(max-min+1) + min
//}

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
func GetAddrListOnXmlByElement(url string) error {
	var err error
	//发送http请求
	resp, err := tool.MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return err
	}
	//将文件读取至etree文档
	doc := etree.NewDocument()
	if i, err := doc.ReadFrom(resp.Body); err != nil || i <= 0 {
		log.Fatal("error:", err, "or read file is nil")
	}
	defer resp.Body.Close()
	//访问根元素
	root := doc.SelectElement("sdnList")
	//访问元素 sdnList->sdnEntry->idList->id->idType
	for _, sdnEntry := range root.SelectElements("sdnEntry") {
		var entityInfo = new(entities.Entity)
		var emailList, websiteList, phoneList []string
		var otherList []entities.OtherInfo
		var idInfoList []entities.ID
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
					otherInfo := entities.OtherInfo{
						Type: idType,
						Info: idNumberValue,
					}
					otherList = append(otherList, otherInfo)
				} else {
					//匹配字符串,判断是否为数据地址信息
					isCurrency, err := regexp.MatchString(constants.CONDITION, idType)
					if err != nil {
						log.Fatal("字符串匹配失败，MatchString Error:", err.Error())
						return err
					}
					//若匹配失败，将id标签下的信息存储到ID列表中
					if !isCurrency {
						idInfo := entities.ID{
							IDType:         idType,
							IDNumber:       idNumberValue,
							IDCountry:      id.SelectElement(constants.ID_COUNTRY).Text(),
							ExpirationDate: id.SelectElement(constants.EXPRI_DATE).Text(),
						}
						idInfoList = append(idInfoList, idInfo)
						//数据匹配成功
					} else {
						//获取entity其他信息
						getEntityInfo(sdnEntry, entityInfo)
						//截取数据：所在链，以及地址,存储到地址风险名单中
						//根据货币名称获取所在链
						chain := getChainBySdn(idType)
						//根据address，查询该地址是否已经存在
						isExist, err := es.IsExistById(constants.ES_ADDRESS, idNumberValue)
						if err != nil {
							log.Printf("IsExistById Error :%s\n，id=%s", err.Error(), idNumberValue)
							continue
						}
						//记录来源信息
						dsAddrInfo := entities.AdsDataSource{
							DsAddr:     url,
							DsType:     strconv.Itoa(constants.DS_OFAC),
							Illustrate: "Derived from the OFAC Sanctions List.",
							Time:       time.Now(),
						}
						//若该地址已经存在，则增加地址来源
						if isExist {
							log.Printf("%s信息已存在,添加该数据来源\n", idNumberValue)
							es.AddDsAddrSource(idNumberValue, dsAddrInfo)
						} else {
							log.Printf("添加%s地址信息\n", idNumberValue)
							//地址不存在，则新建风险名单信息,并存储到es中
							walletAddr := entities.WalletAddr{
								WaAddr:      idNumberValue,
								EntityId:    entityId,
								WaRiskLevel: constants.INIT_LEVEL,
								WaChain:     chain,
								DsAddr: []entities.AdsDataSource{
									dsAddrInfo,
								},
							}
							err = es.Insert(constants.ES_ADDRESS, idNumberValue, walletAddr)
							if err != nil {
								log.Fatal("Fail insert risk address:", err.Error())
							}
						}
					}
				}
			}
			//get name info
		}
		//存储entity信息
		entityInfo.Email = emailList
		entityInfo.Website = websiteList
		entityInfo.PhoneNumber = phoneList
		entityInfo.IDList = idInfoList
		entityInfo.OtherInfo = otherList
		//存入es中
		err = es.Insert(constants.ES_ENTITY, entityId, &entityInfo)
		if err != nil {
			log.Fatal("Fail insert risk entity:", err.Error())
		}
	}
	return nil
}

// 根据指定的element存储sdn名单上的entity信息
func getEntityInfo(sdnEntry *etree.Element, entityInfo *entities.Entity) {
	var akaName, otherAddress []string
	var addrInfoList []entities.AddressList   //地址列表
	var daBiList []entities.DateOfBirth       //出生日期列表
	var plBiList []entities.PlaceOfBirth      //出生地址列表
	var ntList, czList []entities.Nationality //国籍列表
	//获取风险名单实体信息
	sdnType := sdnEntry.SelectElement("sdnType").Text()
	//获取实体名字
	if sdnType == "Individual" {
		entityInfo.IsIndividual = true
		entityInfo.Name = sdnEntry.SelectElement("lastName").Text() + sdnEntry.SelectElement("firstName").Text()
	} else {
		entityInfo.IsIndividual = false
		entityInfo.Name = sdnEntry.SelectElement("lastName").Text()
	}
	//获取实体别名列表
	if akaList := sdnEntry.SelectElement("akaList"); akaList != nil {
		for _, aka := range akaList.SelectElements("aka") {
			akaName = append(akaName, aka.SelectElement("lastName").Text())
		}
	}
	entityInfo.AkaList = akaName
	//获取实体地址列表
	if addressList := sdnEntry.SelectElement("addressList"); addressList != nil {
		for _, address := range addressList.SelectElements("address") {
			country := address.SelectElement("country").Text()
			stateOrProvince := address.SelectElement("stateOrProvince").Text()
			city := address.SelectElement("city").Text()
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
			addressList := entities.AddressList{
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
			daBi := entities.DateOfBirth{
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
			plBi := entities.PlaceOfBirth{
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
			nt := entities.Nationality{
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
			cz := entities.Nationality{
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
