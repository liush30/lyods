// Package domain Dash、Bitcoin Gold交易信息结构
package domain

import "math/big"

type TransactionDashOrBGold struct {
	PagesTotal uint             `json:"pagesTotal"` //总页数，索引值从0开始，若大于1则需要多次查询
	Txs        []TxsDashOrBGold `json:"txs"`        //交易信息
}
type TxsDashOrBGold struct {
	TxId          string            `json:"txid"`          //交易id
	Version       uint              `json:"version"`       //版本
	LockTime      uint              `json:"locktime"`      //锁定时间
	VIn           []VInDashOrBGold  `json:"vin"`           //交易输入信息
	VOut          []VOutDashOrBGold `json:"vout"`          //交易输出信息
	BlockHash     string            `json:"blockhash"`     //区块哈希
	BlockHeight   uint              `json:"blockheight"`   //区块高度
	Confirmations uint              `json:"confirmations"` //确认数量
	Time          big.Int           `json:"time"`          //交易时间
	BlockTime     big.Int           `json:"blocktime"`     //区块时间
	ValueOut      float64           `json:"valueOut"`      //发出的value总值
	Size          uint              `json:"size"`          //交易字节数
	ValueIn       float64           `json:"valueIn"`       //接收的value总值
	Fees          float64           `json:"fees"`          //交易费用
	TxBlock       bool              `json:"txlock"`
}
type VInDashOrBGold struct {
	TxId            string              `json:"txid"`      //交易id
	Vout            uint                `json:"vout"`      //前序交易输出对象
	Sequence        uint                `json:"sequence"`  //发送者定义的交易版本号
	N               int                 `json:"n"`         //交易索引
	ScriptSig       InScrSigDashOrBGold `json:"scriptSig"` //私钥签名
	Addr            string              `json:"addr"`      //交易地址
	ValueSat        big.Int             `json:"valueSat"`
	Value           float64             `json:"value"` //交易金额
	DoubleSpentTxID string              `json:"doubleSpentTxID"`
}
type VOutDashOrBGold struct {
	Value        string               `json:"value"`        //交易金额
	N            int                  `json:"n"`            //交易索引
	ScriptPubKey OutScrSigDashOrBGold `json:"scriptPubKey"` //收款人公钥信息
	SpentTxId    string               `json:"spentTxId"`    //交易id
	SpentIndex   uint                 `json:"spentIndex"`   //索引值
	SpentHeight  uint                 `json:"spentHeight"`  //区块高度
}
type InScrSigDashOrBGold struct {
	Hex string `json:"hex"` //hex格式-16进制的表达
	Asm string `json:"asm"` //asm格式
}
type OutScrSigDashOrBGold struct {
	Hex       string   `json:"hex"`       //hex格式-16进制的表达
	Asm       string   `json:"asm"`       //asm格式
	Addresses []string `json:"addresses"` //交易地址（收款地址）
	Type      string   `json:"type"`      //交易类型：如P2PKH、 Public Key hash
}
