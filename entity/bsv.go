// Package entity Bitcoin Satoshi Vision交易信息结构
package entity

type HistoryTransOnBsv struct {
	TxHash string `json:"tx_hash"`
	Height uint   `json:"height"`
}
type TransactionOnBsv struct {
	TxId     string `json:"txid"`
	Hash     string `json:"hash"`
	Version  uint   `json:"version"`
	Size     uint   `json:"size"`
	LockTime uint   `json:"locktime"`
}
