// Package tool 根据指定的源文件路径，查询地址名单，目前有：json、xml、csv格式
package tool

import (
	"encoding/csv"
	"encoding/json"
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
func GetAddrListOnJSON(url, chain string, level uint, fields []string, road ...string) ([]string, error) {
	//fmt.Printf("*************************clinet:%v\n", MClient)
	var err error
	//用于去除重复数据
	temp := map[string]struct{}{}
	//风险地址名单
	//var addrList []entity.WalletAddr
	var addrListStr []string
	//创建http请求
	req, err := http.NewRequest(param.HTTP_GET, url, nil)
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return nil, err
	}
	//发送http请求
	resp, err := MClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Do Error:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Io Read Error:", err)
		return nil, err
	}
	//解析读取的数据，根据指定路径road，获取指定的json字段fields
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		address, _, _, err := jsonparser.Get(value, fields[0])
		if err != nil {
			log.Fatal(err)
		}
		ticker, _, _, err := jsonparser.Get(value, fields[1])
		if err != nil {
			log.Fatal(err)
		}
		//将不重复数据存入到addrList中
		if _, ok := temp[string(address)]; !ok {
			//将获取到的信息存储到addrList
			walletAddr := entity.WalletAddr{
				WaAddr:      string(address),
				WaRiskLevel: level,
				WaTicker:    string(ticker),
				WaChain:     chain,
				DsAddr: []entity.AdsDataSource{
					{
						DsAddr: url,
						DsType: param.DS_TYPE_URL,
						Number: param.INIT_NUMBER,
					},
				},
			}
			addrByter, err := json.Marshal(walletAddr)
			if err != nil {
				log.Println("json marshal error:", err)
			}
			addrListStr = append(addrListStr, string(addrByter))
		}

	}, road...)
	return addrListStr, nil

}

// GetAddrListOnCsv 根据url获得csv格式的地址名单-批量获取address
func GetAddrListOnCsv(url, chain string, level, index uint) ([]entity.WalletAddr, error) {
	var err error
	//风险地址名单
	var addrList []entity.WalletAddr
	//用于去除重复数据
	temp := map[string]struct{}{}
	//创建http请求
	req, err := http.NewRequest(param.HTTP_GET, url, nil)
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return nil, err
	}
	//发送http请求
	resp, err := MClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	//跳过文件头
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return nil, err
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
		}
		//将不重复数据存入到addrList中
		if _, ok := temp[rec[index]]; !ok {
			temp[rec[index]] = struct{}{}
			addrList = append(addrList, entity.WalletAddr{WaAddr: rec[0], WaRiskLevel: level, WaTicker: chain})
		}
	}
	return addrList, nil
}

// GetAddrListOnXmlByPath 根据url获取xml格式的地址名单-根据指定的路径查询,查询地址
func GetAddrListOnXmlByPath(url, chain, searchPath string, level uint) ([]entity.WalletAddr, error) {
	var err error
	//用于去除重复数据
	temp := map[string]struct{}{}
	//风险地址名单
	var addrList []entity.WalletAddr
	//发送http请求
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	//将文件读取至etree文档
	doc := etree.NewDocument()
	if i, err := doc.ReadFrom(resp.Body); err != nil || i <= 0 {
		log.Fatal("error:", err, "or read file is nil")
		return nil, err
	}
	//根据指定路径查询
	for _, e := range doc.FindElements(searchPath) {
		//将不重复数据存入到addrList中
		if _, ok := temp[e.Text()]; !ok {
			temp[e.Text()] = struct{}{}
			addrList = append(addrList, entity.WalletAddr{WaAddr: e.Text(), WaRiskLevel: level, WaTicker: chain})
		}
	}
	return addrList, nil
}

// GetAddrListOnXmlByElement 根据url获取xml格式的地址名单-根据访问元素，查询所在链以及地址
func GetAddrListOnXmlByElement(url, condition string, level uint) ([]entity.WalletAddr, error) {
	var err error
	//风险地址名单那
	var addrList []entity.WalletAddr
	//发送http请求
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return nil, err
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
				idType := id.SelectElement("idType").Text()
				//匹配字符串
				isCurrency, err := regexp.MatchString(condition, idType)
				if err != nil {
					log.Fatal(err.Error())
					return nil, err
				}
				//数据匹配成功，截取数据：所在链，以及地址,存储到地址风险名单中
				if isCurrency {
					flysnowRegexp, err := regexp.Compile(condition)
					if err != nil {
						log.Fatal(err.Error())
						return nil, err
					}
					params := flysnowRegexp.FindStringSubmatch(idType)
					address := id.SelectElement("idNumber").Text()
					addrList = append(addrList, entity.WalletAddr{WaAddr: address, WaTicker: params[len(params)-1], WaRiskLevel: level})
				}
			}
		}
	}
	return addrList, nil
}
