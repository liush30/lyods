// Package tool 查询bitcoin交易信息及子层级名单，包含链：BTG、XBT、ETH、XMR、LTC、ZEC、DASH、ETC、BSV
package tool

import (
	"encoding/json"
	jsonparser "github.com/buger/jsonparser"
	"io"
	"log"
	"lyods-adsTool/entity"
	"lyods-adsTool/param"
	"net/http"
)

// GetSublistByAddr 根据指定地址查询该地址的子名单信息
func GetSublistByAddr(addr string, road ...string) ([]string, error) {
	var err error
	var subList []string
	//用于去除重复数据
	temp := map[string]struct{}{}
	//获得url
	url := getUrlToBtcTrans(addr)
	//发送http请求。根据url获取到指定账户的所有交易信息
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Io Read Error:", err)
		return nil, err
	}
	//遍历该地址的每一条交易信息，获取转出的交易信息，并将转出对象地址存储
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		isTrue, err := isTransfer(value, addr)
		if err != nil {
			log.Println("isTransfer Error:", err.Error())
			return
		}
		//isTrue==true,该交易为转出交易，获得子名单信息，并获得其交易信息
		if isTrue {
			jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				//获取转出对象地址
				val, err := jsonparser.GetString(value, "addr")
				if err != nil {
					log.Println("jsonparser GetString Error:", err.Error())
					return
				}
				//去除重复的子名单信息，并将其存入到
				if _, ok := temp[val]; !ok {
					temp[val] = struct{}{}
					subList = append(subList, val)
				}
			}, "out")
		}
		return
	}, "txs")
	if err != nil {
		log.Println("ArrayEach Txs:", err.Error())
		return nil, err
	}
	return subList, nil
}

// GetSublistByLevel 根据监控层次获取n层的子名单信息
// 修改：错误信息！& 去除重复信息
func GetSublistByLevel(n int, addr []string) [][]string {
	var levelList [][]string
	//addr为第一层，获取n层
	levelList, _, _, _ = getListByList(levelList, addr, 2, n)
	return levelList
}

// 获取next层所有子名单信息
func getListByList(levelList [][]string, list []string, n, max int) ([][]string, []string, int, int) {
	//按层次存储子名单信息
	if n <= max {
		var allList []string
		for _, v := range list {
			nextList, _ := GetSublistByAddr(v)
			//将查询到的next 名单信息汇总
			allList = append(allList, nextList...)
		}
		levelList = append(levelList, allList)
		return getListByList(levelList, allList, n+1, max)
	}
	return levelList, list, n, max
}

// GetTransOnBtc 根据给定的地址，查询出该地址的所有交易信息
func GetTransOnBtc(addr string) (entity.TransactionBtc, error) {
	var err error
	var trans entity.TransactionBtc
	//根据指定地址查询交易信息
	//获得url
	url := getUrlToBtcTrans(addr)
	//发送http请求
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return entity.TransactionBtc{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return entity.TransactionBtc{}, err
	}
	//将查询到的数据绑定到结构体中
	err = json.Unmarshal(
		body,
		&trans,
	)
	if err != nil {
		log.Fatal("json unmarshal error:", err.Error())
		return entity.TransactionBtc{}, err
	}
	return trans, nil
}

// GetAssocAddr 根据给定地址，获取地址交易关联地址名单
//func GetAssocAddr(addr string) ([]string, error) {
//	var err error
//	//关联地址名单
//	var addrList []string
//	//用于去除重复数据
//	temp := map[string]struct{}{}
//	temp[addr] = struct{}{}
//	trans, err := GetTransOnBtc(addr)
//	if err == nil {
//		//遍历地址的交易信息，将地址关联存于addrList中
//		for _, tx := range trans.Txs {
//			for _, out := range tx.Out {
//				if _, ok := temp[out.Addr]; !ok {
//					addrList = append(addrList, out.Addr)
//				}
//			}
//		}
//	}
//	return addrList, err
//}

// 遍历交易的输入信息，判断该交易是否为转出
func isTransfer(data []byte, addr string) (bool, error) {
	var errTrans error
	var isTrue bool
	//获取输入信息
	_, errTrans = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//获得交易input（输入）信息的地址
		val, err := jsonparser.GetString(value, "prev_out", "addr")
		if err != nil {
			log.Println("jsonparser GetString Error:", err)
			errTrans = err
			return
		}
		//若输入信息中的地址等于addr，说明该交易为地址的转出交易
		if addr == val {
			isTrue = true
		}
	}, "inputs")
	return isTrue, errTrans
}
func getUrlToBtcTrans(addr string) string {
	return param.API_BTC_TRANS + addr
}
