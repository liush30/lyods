// Package tool 根据指定的源文件路径，查询地址名单，目前有：json、xml、csv格式
package tool

import (
	"encoding/csv"
	"github.com/beevik/etree"
	"github.com/buger/jsonparser"
	"io"
	"log"
	"lyods-adsTool/entity"
	"lyods-adsTool/param"
	"net/http"
	"regexp"
)

// GetAddrListOnJSON 根据url获取json格式的地址名单-获取address & chain
func GetAddrListOnJSON(url string, level uint, fields []string, road ...string) error {
	var err error
	//用于去除重复数据
	temp := map[string]struct{}{}
	//创建http请求
	req, err := http.NewRequest(param.HTTP_GET, url, nil)
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return err
	}
	//发送http请求
	resp, err := MClient.Do(req)
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
		address, _, _, err := jsonparser.Get(value, fields[0])
		if err != nil {
			log.Fatal(err.Error())
		}
		addrStr := string(address)
		chain, _, _, err := jsonparser.Get(value, fields[1])
		if err != nil {
			log.Fatal(err.Error())
		}
		//将不重复数据存入到addrList中
		if _, ok := temp[addrStr]; !ok {
			temp[addrStr] = struct{}{}
			//根据address，查询该地址是否已经存在与
			isExist, err := IsExistById(param.WA_ADDR_INDEX, addrStr)
			if err != nil {
				log.Printf("IsExistById Error :%v\n", err.Error())
			}
			dsAddrInfo := entity.AdsDataSource{
				DsAddr: url,
				DsType: param.DS_TYPE_URL,
				Number: param.INIT_NUMBER,
			}
			//若该地址已经存在，则更新地址来源
			if isExist {
				log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
				UpdateDsAddrSource(addrStr, dsAddrInfo)
			} else {
				log.Printf("添加%s地址信息\n", addrStr)
				//地址不存在，则新建
				walletAddr := entity.WalletAddr{
					WaAddr:      addrStr,
					WaRiskLevel: level,
					//WaTicker:    string(ticker),
					WaChain: string(chain),
					DsAddr: []entity.AdsDataSource{
						dsAddrInfo,
					},
				}
				Insert(param.WA_ADDR_INDEX, addrStr, walletAddr)
			}
		}

	}, road...)
	return nil
}

// GetAddrListOnCsv 根据url获得csv格式的地址名单-批量获取address
func GetAddrListOnCsv(url, ticker string, level, index uint) error {
	var err error
	//风险地址名单
	//var addrList []entity.WalletAddr
	//用于去除重复数据
	temp := map[string]struct{}{}
	//创建http请求
	req, err := http.NewRequest(param.HTTP_GET, url, nil)
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return err
	}
	//发送http请求
	resp, err := MClient.Do(req)
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
		//文件读取完成
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			continue
		}
		//将不重复数据存入到addrList中
		if _, ok := temp[rec[index]]; !ok {
			temp[rec[index]] = struct{}{}
			addrStr := rec[0] //地址
			//根据address，查询该地址是否已经存在
			isExist, err := IsExistById(param.WA_ADDR_INDEX, addrStr)
			if err != nil {
				log.Printf("IsExistById Error :%v\n", err.Error())
				continue
			}
			dsAddrInfo := entity.AdsDataSource{
				DsAddr: url,
				DsType: param.DS_TYPE_URL,
				Number: param.INIT_NUMBER,
			}
			//若该地址已经存在，则更新地址来源
			if isExist {
				log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
				UpdateDsAddrSource(addrStr, dsAddrInfo)
			} else {
				log.Printf("添加%s地址信息\n", addrStr)
				//地址不存在，则新建
				walletAddr := entity.WalletAddr{
					WaAddr:      addrStr,
					WaRiskLevel: level,
					WaTicker:    ticker,
					WaChain:     GetChainByTicker(ticker),
					DsAddr: []entity.AdsDataSource{
						dsAddrInfo,
					},
				}
				Insert(param.WA_ADDR_INDEX, addrStr, walletAddr)
			}
		}
	}
	return nil
}

// GetAddrListOnXmlByPath 根据url获取xml格式的地址名单-根据指定的路径查询,查询地址
func GetAddrListOnXmlByPath(url, ticker, searchPath string, level uint) error {
	var err error
	//用于去除重复数据
	temp := map[string]struct{}{}
	//发送http请求
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return err
	}
	defer resp.Body.Close()
	//将文件读取至etree文档
	doc := etree.NewDocument()
	if i, err := doc.ReadFrom(resp.Body); err != nil || i <= 0 {
		log.Fatal("error:", err, "or read file is nil")
		return err
	}
	//根据指定路径查询
	for _, e := range doc.FindElements(searchPath) {
		//将不重复数据存入到addrList中
		if _, ok := temp[e.Text()]; !ok {
			temp[e.Text()] = struct{}{}
			addrStr := e.Text() //地址
			//根据address，查询该地址是否已经存在
			isExist, err := IsExistById(param.WA_ADDR_INDEX, addrStr)
			if err != nil {
				log.Printf("IsExistById Error :%v\n", err.Error())
				continue
			}
			dsAddrInfo := entity.AdsDataSource{
				DsAddr: url,
				DsType: param.DS_TYPE_URL,
				Number: param.INIT_NUMBER,
			}
			//若该地址已经存在，则更新地址来源
			if isExist {
				log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
				UpdateDsAddrSource(addrStr, dsAddrInfo)
			} else {
				log.Printf("添加%s地址信息\n", addrStr)
				//地址不存在，则新建
				walletAddr := entity.WalletAddr{
					WaAddr:      addrStr,
					WaRiskLevel: level,
					WaTicker:    ticker,
					WaChain:     GetChainByTicker(ticker),
					DsAddr: []entity.AdsDataSource{
						dsAddrInfo,
					},
				}
				Insert(param.WA_ADDR_INDEX, addrStr, walletAddr)
			}
		}
	}
	return nil
}

// GetAddrListOnXmlByElement 根据url获取xml格式的地址名单-根据访问元素，查询所在链以及地址
func GetAddrListOnXmlByElement(url, condition string, level uint) error {
	var err error
	//发送http请求
	resp, err := MClient.Get(url)
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
		if idList := sdnEntry.SelectElement("idList"); idList != nil {
			for _, id := range idList.SelectElements("id") {
				//获取idType
				idType := id.SelectElement("idType").Text()
				//根据获得到的idType匹配字符串
				isCurrency, err := regexp.MatchString(condition, idType)
				if err != nil {
					log.Fatal("字符串匹配失败，MatchString Error:", err.Error())
					return err
				}
				//数据匹配成功，截取数据：所在链，以及地址,存储到地址风险名单中
				if isCurrency {
					flysnowRegexp, err := regexp.Compile(condition)
					if err != nil {
						log.Fatal(err.Error())
						return err
					}
					params := flysnowRegexp.FindStringSubmatch(idType)
					addrStr := id.SelectElement("idNumber").Text()
					//根据address，查询该地址是否已经存在
					isExist, err := IsExistById(param.WA_ADDR_INDEX, addrStr)
					if err != nil {
						log.Printf("IsExistById Error :%v\n", err.Error())
						continue
					}
					dsAddrInfo := entity.AdsDataSource{
						DsAddr: url,
						DsType: param.DS_TYPE_URL,
						Number: param.INIT_NUMBER,
					}
					//若该地址已经存在，则更新地址来源
					if isExist {
						log.Printf("%s信息已存在,添加该数据来源\n", addrStr)
						UpdateDsAddrSource(addrStr, dsAddrInfo)
					} else {
						log.Printf("添加%s地址信息\n", addrStr)
						ticker := params[len(params)-1]
						//地址不存在，则新建
						walletAddr := entity.WalletAddr{
							WaAddr:      addrStr,
							WaRiskLevel: level,
							WaTicker:    ticker,
							WaChain:     GetChainByTicker(ticker),
							DsAddr: []entity.AdsDataSource{
								dsAddrInfo,
							},
						}
						Insert(param.WA_ADDR_INDEX, addrStr, walletAddr)
					}
				}
			}
		}
	}
	return nil
}
