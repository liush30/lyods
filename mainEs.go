package main

import (
	"fmt"
	"lyods-adsTool/tool"
)

func main() {
	//// 创建ES client用于后续操作ES
	//client, err := elastic.NewClient(
	//	// 设置ES服务地址，支持多个地址
	//	elastic.SetURL("http://127.0.0.1:9200"),
	//	// 设置基于http base auth验证的账号和密码
	//	elastic.SetBasicAuth("user", "secret"))
	//if err != nil {
	//	// Handle error
	//	fmt.Printf("连接失败: %v\n", err)
	//} else {
	//	fmt.Println("连接成功")
	//}
	//// 执行ES请求需要提供一个上下文对象
	//ctx := context.Background()
	//// 首先检测下addrList索引是否存在
	//exists, err := client.IndexExists("addrList").Do(ctx)
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//if !exists {
	//	// addrList索引不存在，则创建一个
	//	_, err := client.CreateIndex("addrList").BodyString(param.ADDR_MAPPING).Do(ctx)
	//	if err != nil {
	//		// Handle error
	//		panic(err)
	//	}
	//}
	//walletAddr := entity.WalletAddr{
	//	WaAddr:      string("0xooooooo88888888"),
	//	WaRiskLevel: 1,
	//	WaTicker:    string("BTC"),
	//	WaChain:     "bitcoin",
	//	DsAddr: []entity.AdsDataSource{
	//		{
	//			DsAddr: "xxxxxxx.test0000.com",
	//			DsType: "1",
	//			Number: 1,
	//		},
	//	},
	//}
	//err := tool.Insert("addr_list", walletAddr.WaAddr, walletAddr)
	//if err != nil {
	//	fmt.Printf("Insert =========Error=========:\v\n", err.Error())
	//}
	//bol, err := tool.IsExistById("addr_list", "13echkf999999999999999999999")
	////err, waAddr := tool.GetWalletAddr("addr_list", "0xooooooo88888888")
	//if err != nil {
	//	fmt.Printf("GetWalletAddr =========Error=========:\v\n", err.Error())
	//}
	//fmt.Println(bol)
	count, err := tool.GetIndexDocNum("addr_list")
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Println(count)
	//dsAddr := entity.AdsDataSource{
	//	DsAddr: "888888.test0000.com",
	//	DsType: "2",
	//	Number: 1,
	//}
	//err = tool.UpdateDsAddrSource("0xooooooo88888888", dsAddr)
	//if err != nil {
	//	fmt.Printf("UpdateDsAddrSource =========Error=========:\v\n", err.Error())
	//}
	//err, waAddr := tool.GetWalletAddr("addr_list", "13echkfV54pw7fxN1gkQmg5yV8FC4xx16T")
	//if err != nil {
	//	fmt.Printf("GetWalletAddr =========Error=========:\v\n", err.Error())
	//}
	//fmt.Println(waAddr)
	//tool.GetAddrListOnJSON("http://api.ransomwhe.re/export", 1, []string{"address", "blockchain"}, "result")
}
