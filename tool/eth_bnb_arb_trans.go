// Package tool 查询ethereum交易信息及层级信息
package tool

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"lyods-adsTool/entity"
	"lyods-adsTool/param"
	"net/http"
)

// GetTransOnEthOrBnbOrArb 根据指定地址查询普通交易信息
func GetTransOnEthOrBnbOrArb(chainType uint, addr string) (entity.TransactionEthOrBnb, error) {
	var err error
	//普通交易信息
	var trans entity.TransactionEthOrBnb

	//获取url-根据地址查询以太坊普通交易信息
	url := getNormalTransUrl(chainType, addr)
	//发送HTTP请求
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return entity.TransactionEthOrBnb{}, err
	}
	defer resp.Body.Close()
	////读取数据
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		log.Fatal("Io Read Error:", err.Error())
		return entity.TransactionEthOrBnb{}, err
	}
	//将获取的数据反序列化为结构体信息
	err = json.Unmarshal(body, &trans)
	if err != nil {
		log.Fatal("Json Unmarshal Error:", err.Error())
		return entity.TransactionEthOrBnb{}, err
	}
	return trans, nil
}

// GetTransInOnEthOrBnbOrArb 根据指定地址查询内部交易信息
func GetTransInOnEthOrBnbOrArb(chainType uint, addr string) (entity.TransactionEthOrBnb, error) {
	var err error
	//内部交易信息
	var trans entity.TransactionEthOrBnb
	//获取url-根据地址查询以太坊内部交易信息
	url := getInternalTransUrl(chainType, addr)
	//若url为空，说明输入的chain type不正确
	if url == "" || len(url) == 0 {
		log.Println("incorrect chain type")
		return entity.TransactionEthOrBnb{}, errors.New("incorrect chain type")
	}
	//发送HTTP请求
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return entity.TransactionEthOrBnb{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Io Read Error:", err.Error())
		return entity.TransactionEthOrBnb{}, err
	}
	//将获取的数据反序列化为结构体信息
	err = json.Unmarshal(body, &trans)
	if err != nil {
		log.Fatal(err.Error())
		return entity.TransactionEthOrBnb{}, err
	}
	return trans, err
}

// GetUsdtTransOnEth 根据指定地址查询erc20转账信息
func GetUsdtTransOnEth(addr string) (entity.TransactionEthOrBnb, error) {
	var err error
	var trans entity.TransactionEthOrBnb
	//获取url-根据地址查询以太坊中的USDT交易信息
	url := getUsdtTransUrl(addr)
	//若url为空，说明输入的chain type不正确
	if url == "" || len(url) == 0 {
		log.Println("incorrect chain type")
		return entity.TransactionEthOrBnb{}, errors.New("incorrect chain type")
	}
	//发送HTTP请求
	resp, err := MClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("http status is :", resp.StatusCode, "Do Error:", err.Error())
		return entity.TransactionEthOrBnb{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return entity.TransactionEthOrBnb{}, err
	}
	//将获取的数据反序列化为结构体信息
	err = json.Unmarshal(body, &trans)
	if err != nil {
		return entity.TransactionEthOrBnb{}, err
	}
	return trans, nil

}

// 根据chain 类型获取指定地址的普通交易
func getNormalTransUrl(chainType uint, addr string) string {
	if chainType == param.CHAIN_BNB {
		return getNormalUrlBnb(addr)
	} else if chainType == param.CHAIN_ARBITRUM {
		return getNormalUrlArb(addr)
	} else if chainType == param.CHAIN_ETHEREUM {
		return getNormalUrlEth(addr)
	}
	return ""
}

// 根据chain类型获取地址的内部交易
func getInternalTransUrl(chainType uint, addr string) string {
	if chainType == param.CHAIN_BNB {
		return getInterUrlBnb(addr)
	} else if chainType == param.CHAIN_ETHEREUM {
		return getInternalEthUrl(addr)
	} else if chainType == param.CHAIN_ARBITRUM {
		return getInterUrlArb(addr)
	}
	return ""
}

// 根据地址获取ethereum中usdt交易请求url
func getUsdtTransUrl(addr string) string {
	return param.API_ETH_USDT + addr
}

// bnb根据地址获取内部交易请求url
func getNormalUrlBnb(addr string) string {
	return param.API_BNB_TRANS + addr
}

// bnb根据地址获取内部交易请求url
func getInterUrlBnb(addr string) string {
	return param.API_BNB_INTRANS + addr
}

// ethereum根据地址获取内部交易请求url
func getInternalEthUrl(addr string) string {
	return param.API_ETH_INTRANS + addr
}

// ethereum根据地址获取普通交易请求url
func getNormalUrlEth(addr string) string {
	return param.API_ETH_TRANS + addr
}

// arbitrum根据地址获取普通交易url
func getNormalUrlArb(addr string) string {
	return param.API_ARB_TRANS + addr
}

// arbitrum根据地址获得内部交易url
func getInterUrlArb(addr string) string {
	return param.API_ARB_INTRANS + addr
}
