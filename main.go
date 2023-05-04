package main

import (
	"encoding/csv"
	"fmt"
	"github.com/beevik/etree"
	"github.com/buger/jsonparser"
	"io"
	"log"
	t "lyods-adsTool/entity"
	"net/http"
	"regexp"
)

func main() {
	//解析数据来源参数：目前支持格式：json、csv、xml
	//getDataSourceByJSON("https://api.ransomwhe.re/export")·路劲
	// getDataSourceByCsv("https://gist.githubusercontent.com/banteg/1657d4778eb86c460e03bc58b99970c0/raw/2b8e0b2c1074b995b992397f34ab2843cf6bdf72/uniswap-trm.csv", "ETH", 1)
	//list := getAddrListByJSON("http://api.ransomwhe.re/export", 1, []string{"address", "blockchain"}, "result")
	//for _, v := range list {
	//	fmt.Println(v)
	//}
	//fmt.Println(len(list))
	// getAddrListOnXml("https://www.treasury.gov/ofac/downloads/sdn.xml", `/sdnList/sdnEntry/idList/id[idType='Digital Currency Address - XBT']/idNumber`)
	getAddrListOnXmlByElement("https://www.treasury.gov/ofac/downloads/sdn.xml")
}

// 根据url获取json格式的地址名单-获取address & chain
func getAddrListOnJSON(url string, level int, fields []string, road ...string) []t.WalletAddr {
	var err error
	//风险地址名单
	var addrList []t.WalletAddr
	//获取数据
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, err.Error())
		return addrList
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return addrList
	}
	//解析读取的数据，根据指定路径road，获取指定的json字段fields
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		address, _, _, err := jsonparser.Get(value, fields[0])
		if err != nil {
			log.Fatal(err)
		}
		chain, _, _, err := jsonparser.Get(value, fields[1])
		if err != nil {
			log.Fatal(err)
		}
		//将获取到的信息存储到addrList
		addrList = append(addrList, t.WalletAddr{
			WaAddr:      string(address),
			WaRiskLevel: level,
			WaChain:     string(chain),
		})
	}, road...)
	return addrList
}

// 根据url获得csv格式的地址名单-批量获取address
func getAddrListOnCsv(url, chain string, level, index int) []t.WalletAddr {
	var err error
	//风险地址名单
	var addrList []t.WalletAddr
	//用于去除重复数据
	temp := map[string]struct{}{}
	//获取数据
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, err.Error())
	}
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	//跳过文件头
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return []t.WalletAddr{}
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
			addrList = append(addrList, t.WalletAddr{WaAddr: rec[0], WaRiskLevel: level, WaChain: chain})
		}
	}
	return addrList
}

// 根据url获取xml格式的地址名单-根据指定的路径查询
func getAddrListOnXmlByPath(url, chain, searchPath string, level int) {
	var err error
	//用于去除重复数据
	temp := map[string]struct{}{}
	//风险地址名单
	var addrList []t.WalletAddr
	resp, err := http.Get(url)
	if err != nil {
		log.Println("http status is :", resp.StatusCode, err.Error())
	}
	defer resp.Body.Close()
	//将文件读取至etree文档
	doc := etree.NewDocument()
	if i, err := doc.ReadFrom(resp.Body); err != nil || i <= 0 {
		log.Fatal("error:", err, "or read file is nil")
	}
	//根据指定路径查询
	for _, e := range doc.FindElements(searchPath) {
		//将不重复数据存入到addrList中
		if _, ok := temp[e.Text()]; !ok {
			temp[e.Text()] = struct{}{}
			addrList = append(addrList, t.WalletAddr{WaAddr: e.Text(), WaRiskLevel: level, WaChain: chain})
		}
		fmt.Println(e.Text())
	}
}

// 根据url获取xml格式的地址名单-根据访问元素
func getAddrListOnXmlByElement(url string) {
	var err error
	resp, err := http.Get(url)
	if err != nil {
		log.Println("http status is :", resp.StatusCode, err.Error())
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
				isCurrency, err := regexp.MatchString(`^Digital Currency Address - `, idType)
				if err != nil {
					log.Fatal(err.Error())
				}
				//数据匹配成功，截取数据：所在链，以及地址
				if isCurrency {
					flysnowRegexp, err := regexp.Compile(`^Digital Currency Address - ([\D]{3,16}$)`)
					if err != nil {
						log.Fatal(err.Error())
					}
					params := flysnowRegexp.FindStringSubmatch(idType)

					address := id.SelectElement("idNumber").Text()
					fmt.Printf("%v \t %s \n", params[len(params)-1], address)
				}
			}
		}
	}
}
