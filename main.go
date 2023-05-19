package main

import (
	"fmt"
	"lyods-adsTool/tool"
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
	//list := tool.GetAddrListOnXmlByElement("https://www.treasury.gov/ofac/downloads/sdn.xml", `^Digital Currency Address - ([\D]{3,16}$)`, 1)
	//temp := map[string]struct{}{}
	//var tickers []string
	//for _, v := range list {
	//	if _, ok := temp[v.WaTicker]; !ok {
	//		temp[v.WaTicker] = struct{}{}
	//		tickers = append(tickers, v.WaTicker)
	//	}
	//}
	//for _, v := range tickers {
	//	fmt.Printf("%v\t", v)
	//}
	//fmt.Println()
	//fmt.Println(len(tickers))
	//fmt.Println(len(list))
	//trans := tool.GetTransOnBtc("17TMc2UkVRSga2yYvuxSD9Q1XyB2EPRjTF")
	//tran, err := tool.GetTransOnEth("0x03Cf40B900971561AC6bd997ef1Fe939DcbC95e2")
	//if err != nil {
	//	err.Error()
	//}
	//for _, v := range tran.Result {
	//	fmt.Println(v)
	//}
	//trans, err := tool.GetTransOnDash("Xs3vzQmNvAxRa3Xo8XzQqUb3BMgb9EogF4", 0)
	//url:entity.API_DASH_TRANS + addr + "&pageNum=" + strconv.Itoa(int(pageNum))
	//trans, err := tool.GetTransOnBtg("GPwg6991XoHqQPNmAucFACuQ5H9sGCDv9TpS", 0)
	//if err == nil {
	//	fmt.Println(trans.Txs)
	//}
	//trans, err := tool.GetTransRecodeOnVerge("DFFJhnQNZf8rf67tYnesPu7MuGUpYtzv7Z")
	//trans, err := tool.GetTransInfoOnVerge("1cb9cbf4cc26d6687c4f27dd29be4f3e6baf13a0692d84800d746c35140de89b")
	//if err == nil {
	//	fmt.Println(trans)
	//
	//err := tool.GetSublistByAddr("17TMc2UkVRSga2yYvuxSD9Q1XyB2EPRjTF")
	//if err != nil {
	//	fmt.Println("Error:", err.Error())
	//}
	list := []string{"17TMc2UkVRSga2yYvuxSD9Q1XyB2EPRjTF"}
	aaa := tool.GetSublistByLevel(3, list)
	for i, v := range aaa {
		fmt.Printf("第%d层地址\n", i)
		for _, v2 := range v {
			fmt.Printf("%s \t", v2)
		}
		fmt.Println()
	}
	//list, _ := tool.GetAssocAddr("17TMc2UkVRSga2yYvuxSD9Q1XyB2EPRjTF")
	//for _, v := range list {
	//	fmt.Println(v)
	//}
}
