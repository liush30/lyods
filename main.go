package main

import (
	"fmt"
	tool "lyods-adsTool/tool"
)

func main() {
	//解析数据来源参数：目前支持格式：json、csv、xml
	//getDataSourceByJSON("https://api.ransomwhe.re/export")
	// getDataSourceByCsv("https://gist.githubusercontent.com/banteg/1657d4778eb86c460e03bc58b99970c0/raw/2b8e0b2c1074b995b992397f34ab2843cf6bdf72/uniswap-trm.csv", "ETH", 1)
	//list := getAddrListByJSON("http://api.ransomwhe.re/export", 1, []string{"address", "blockchain"}, "result")
	//for _, v := range list {
	//	fmt.Println(v)
	//}
	//fmt.Println(len(list))
	// getAddrListOnXml("https://www.treasury.gov/ofac/downloads/sdn.xml", `/sdnList/sdnEntry/idList/id[idType='Digital Currency Address - XBT']/idNumber`)
	//list := getAddrListOnXmlByElement("https://www.treasury.gov/ofac/downloads/sdn.xml", `^Digital Currency Address - ([\D]{3,16}$)`, 1)
	//fmt.Println(len(list))
	//trans := tool.GetTransOnBtc("17TMc2UkVRSga2yYvuxSD9Q1XyB2EPRjTF")
	list, err := tool.GetAssocAddr("17TMc2UkVRSga2yYvuxSD9Q1XyB2EPRjTF")
	if err != nil {
		err.Error()
	}
	for _, v := range list {
		fmt.Println(v)
	}
}
