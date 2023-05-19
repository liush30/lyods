// Package tool 查询dash交易信息及层级信息
package tool

import (
	"encoding/json"
	"io"
	"log"
	"lyods-adsTool/entity"
	"lyods-adsTool/param"
	"net/http"
	"strconv"
)

// GetTransOnDashOrBitGold 根据指定url获取指定账户Dash或Bitcoin Gold上的交易信息,pageNum,页数(默认为零)
func GetTransOnDashOrBitGold(chainType, pageNum uint, addr string) (entity.TransactionDashOrBGold, error) {
	var err error
	var trans entity.TransactionDashOrBGold
	//根据指定地址以及page Num查询交易信息
	//获取url
	url := getUrl(chainType, pageNum, addr)
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return entity.TransactionDashOrBGold{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("IO Read Error:", err.Error())
		return entity.TransactionDashOrBGold{}, err
	}
	//将数据反序列化为json数组，存储到结构体中
	err = json.Unmarshal(body, &trans)
	if err != nil {
		log.Println("Json Unmarshal:", err.Error())
		return entity.TransactionDashOrBGold{}, err
	}
	return trans, nil
}

func getUrl(chainType, pageNum uint, addr string) string {
	if chainType == param.CHAIN_DASH {
		return getUrlDash(addr, pageNum)
	} else {
		return getUrlBitGold(addr, pageNum)
	}
}

// 获得dash链的根据指定地址查询交易信息
func getUrlDash(addr string, pageNum uint) string {
	return param.API_DASH_TRANS + addr + "&pageNum=" + strconv.Itoa(int(pageNum))
}

// 获得bitgold链的根据指定地址查询交易信息
func getUrlBitGold(addr string, pageNum uint) string {
	return param.API_BTG_TRANS + addr + "&pageNum=" + strconv.Itoa(int(pageNum))
}
