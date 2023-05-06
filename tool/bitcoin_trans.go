// Package tool 查询bitcoin交易信息及子层级名单，包含链：BTG、XBT、ETH、XMR、LTC、ZEC、DASH、ETC、BSV
package tool

import (
	"encoding/json"
	"io"
	"log"
	entity "lyods-adsTool/entity"
	"net/http"
)

// GetTransAndAddrListOnBtc 获得指定地址的BTC交易以及第n层的交易数据
func GetTransAndAddrListOnBtc(address string, level int) {

}

// GetTransOnBtc 根据给定的地址，查询出该地址的所有交易信息
func GetTransOnBtc(addr string) (entity.TransactionBtc, error) {
	var err error
	var trans entity.TransactionBtc
	//根据指定地址查询交易信息
	resp, err := http.Get(entity.ApiBitcoinTrans + addr)
	if err != nil {
		log.Println(err.Error())
		return trans, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return trans, err
	}
	//将查询到的数据绑定到结构体中
	json.Unmarshal(
		body,
		&trans,
	)
	return trans, nil
}

// GetAssocAddr 根据给定地址，获取地址交易关联地址名单
func GetAssocAddr(addr string) ([]string, error) {
	var err error
	//关联地址名单
	var addrList []string
	//用于去除重复数据
	temp := map[string]struct{}{}
	temp[addr] = struct{}{}
	trans, err := GetTransOnBtc(addr)
	if err == nil {
		//遍历地址的交易信息，将地址关联存于addrList中
		for _, tx := range trans.Txs {
			for _, out := range tx.Out {
				if _, ok := temp[out.Addr]; !ok {
					addrList = append(addrList, out.Addr)
				}
			}
		}
	}
	return addrList, err

}
