package tool

import (
	"encoding/json"
	"io"
	"log"
	"lyods-adsTool/config"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"net/http"
)

// GetTransRecodeOnVerge 根据指定地址获取地址的交易记录
//func GetTransRecodeOnVerge(addr string) ([]string, error) {
//var err error
//var trans domain.TransRecordByAddrOnVerge
////根据addr去查询地址的交易记录
//resp, err := MClient.Get(getTransRecodeUrl(addr))
//if err != nil || resp.StatusCode != http.StatusOK {
//	log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
//	return nil, err
//}
//defer resp.Body.Close()
//body, err := io.ReadAll(resp.Body)
//if err != nil {
//	log.Fatal(err.Error())
//	return nil, err
//}
////获得并遍历该地址的交易列表
//_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//	transHash, err := jsonparser.GetString(value, "txid")
//	IsError(err, "Fail get tx id")
//	transType, err := jsonparser.GetString(value, "type") //理解为发送解决。vin为交易的发送方，vout为交易的接收方
//	IsError(err, "Fail get type")
//}, "data")
//if err != nil {
//	log.Println("ArrayEach Txs:", err.Error())
//	return nil, err
//}
//////将查询到的数据绑定到结构体中
////err = json.Unmarshal(
////	body,
////	&trans,
////)
////if err != nil {
////	log.Fatal("json unmarshal error:", err.Error())
////	return domain.TransRecordByAddrOnVerge{}, err
////}
//return trans, nil
//}

// GetTransInfoOnVerge 根据交易id获得具体交易信息
func GetTransInfoOnVerge(txId string) (domain.TransactionOnVerge, error) {
	//根据地址查询出交易记录，根据交易类型，查询转出交易记录
	var err error
	//根据addr去查询地址的交易记录
	resp, err := utils.CreateClient().Get(getTransInfoUrl(txId))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return domain.TransactionOnVerge{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return domain.TransactionOnVerge{}, err
	}
	var trans domain.TransactionOnVerge
	//将查询到的数据绑定到结构体中
	err = json.Unmarshal(
		body,
		&trans,
	)
	if err != nil {
		log.Fatal("json unmarshal error:", err.Error())
		return domain.TransactionOnVerge{}, err
	}
	return trans, nil
}

// 获取查询指定地址交易记录的url
func getTransRecodeUrl(addr string) string {
	return constants.API_VERGE_ADDRTRANS + addr + "/0/" + config.MAX_RECODE
}

// 获取查询指定交易id的交易信息的url
func getTransInfoUrl(txId string) string {
	return constants.API_VERGE_TRANS + txId
}
