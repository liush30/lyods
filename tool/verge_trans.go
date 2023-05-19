package tool

import (
	"encoding/json"
	"io"
	"log"
	"lyods-adsTool/entity"
	"lyods-adsTool/param"
	"net/http"
)

// GetTransRecodeOnVerge 根据指定地址获取地址的交易记录
func GetTransRecodeOnVerge(addr string) (entity.TransRecordByAddrOnVerge, error) {
	var err error
	var trans entity.TransRecordByAddrOnVerge
	//根据addr去查询地址的交易记录
	resp, err := MClient.Get(getTransRecodeUrl(addr))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return entity.TransRecordByAddrOnVerge{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return entity.TransRecordByAddrOnVerge{}, err
	}
	//将查询到的数据绑定到结构体中
	err = json.Unmarshal(
		body,
		&trans,
	)
	if err != nil {
		log.Fatal("json unmarshal error:", err.Error())
		return entity.TransRecordByAddrOnVerge{}, err
	}
	return trans, nil
}

// GetTransInfoOnVerge 根据交易id获得具体交易信息
func GetTransInfoOnVerge(txId string) (entity.TransactionOnVerge, error) {
	//根据地址查询出交易记录，根据交易类型，查询转出交易记录
	var err error
	//根据addr去查询地址的交易记录
	resp, err := MClient.Get(getTransInfoUrl(txId))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return entity.TransactionOnVerge{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return entity.TransactionOnVerge{}, err
	}
	var trans entity.TransactionOnVerge
	//将查询到的数据绑定到结构体中
	err = json.Unmarshal(
		body,
		&trans,
	)
	if err != nil {
		log.Fatal("json unmarshal error:", err.Error())
		return entity.TransactionOnVerge{}, err
	}
	return trans, nil
}

// 获取查询指定地址交易记录的url
func getTransRecodeUrl(addr string) string {
	return param.API_VERGE_ADDRTRANS + addr + "/0/" + param.MAX_RECODE
}

// 获取查询指定交易id的交易信息的url
func getTransInfoUrl(txId string) string {
	return param.API_VERGE_TRANS + txId
}
