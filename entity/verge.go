// Package entity XVG交易信息
package entity

import (
	"math/big"
)

// TransRecordByAddrOnVerge 指定地址的交易记录
type TransRecordByAddrOnVerge struct {
	Data  []RecordData `json:"data"`
	Total uint         `json:"total"` //交易总量
}
type RecordData struct {
	TxId  string  `json:"txid"`  //交易id
	Time  big.Int `json:"time"`  //交易时间
	Type  string  `json:"type"`  //交易类型 vin:作为输出方，转出钱 ，vout:作为输入方，收到钱
	Value string  `json:"value"` //交易金额
}
type TransactionOnVerge struct {
	Data TransData `json:"data"`
}
type TransData struct {
	TxId          string             `json:"txid"` //交易ID
	Hash          string             `json:"hash"` //交易hash
	Version       uint               `json:"version"`
	Size          uint               `json:"size"`
	VSize         uint               `json:"vsize"`
	Weight        uint               `json:"weight"`
	LockTime      uint               `json:"locktime"`
	Time          uint               `json:"time"`
	VIn           []TransVIn         `json:"vin"`
	VOut          []VOutScriptPubKey `json:"vout"`
	Hex           string             `json:"hex"`
	BlockHash     string             `json:"blockhash"`
	Confirmations uint               `json:"confirmations"`
	BlockTime     uint               `json:"blocktime"`
}
type TransVIn struct {
	TxId      string       `json:"txid"` //交易id
	Vout      uint         `json:"vout"`
	ScriptSig VInScriptSig `json:"scriptSig"`
	Sequence  uint         `json:"sequence"`
}
type TransVOut struct {
	Value big.Int `json:"value"`
	N     uint    `json:"n"`
}
type VInScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}
type VOutScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   uint     `json:"reqSigs"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}
