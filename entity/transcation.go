package entity

type Transaction struct {
	TxId    string `json:"tx_id"` //交易哈希
	Address string `json:"address"`
}
