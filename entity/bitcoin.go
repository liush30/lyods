// Package entity bitcoin交易信息结构
package entity

import "math/big"

type TransactionBtc struct {
	Hash160       string   `json:"hash160"`        //BTC地址的hash160格式
	Address       string   `json:"address"`        //BTC地址的base58表示
	NTx           uint     `json:"n_tx"`           //交易数量
	NUnredeemed   big.Int  `json:"n_unredeemed"`   //未赎回数量
	TotalReceived big.Int  `json:"total_received"` //收到的BTC总量
	TotalSent     big.Int  `json:"total_sent"`     //发出的BTC总量
	FinalBalance  big.Int  `json:"final_balance"`  //最终余额列表
	Txs           []TxsBtc `json:"txs"`            //交易信息列表
}
type TxsBtc struct {
	Hash        string      `json:"hash"`         //交易哈希
	Ver         uint        `json:"ver"`          //版本号
	VinSz       uint        `json:"vin_sz"`       //交易输入数量
	VoutSz      uint        `json:"vout_sz"`      //交易输出数量
	Size        uint        `json:"size"`         //交易字节数
	Weight      uint        `json:"weight"`       //权重
	Fee         uint        `json:"fee"`          //交易费用
	RelayedBy   string      `json:"relayed_by"`   //中继节点IP
	LockTime    uint        `json:"lock_time"`    //锁定时间
	TxIndex     big.Int     `json:"tx_index"`     //交易索引
	DoubleSpend bool        `json:"double_spend"` //是否双花
	Time        big.Int     `json:"time"`         //当前交易发生时间
	BlockIndex  uint        `json:"block_index"`  //区块索引
	BlockHeight uint        `json:"block_height"` //区块高度
	Inputs      []InputsBtc `json:"inputs"`       //交易输入
	Out         []OutBtc    `json:"out"`          //交易输出
}
type InputsBtc struct {
	Sequence big.Int    `json:"sequence"` //发送者定义的交易版本号
	Witness  string     `json:"witness"`  //交易输入内容
	Script   string     `json:"script"`   //前序交易输出的目标公钥脚本
	Index    uint       `json:"index"`    //所处交易索引值
	PrevOut  PrevOutBtc `json:"prev_out"` //前序交易输出对象
}
type OutBtc struct {
	Type              int                    `json:"type"`
	Spent             bool                   `json:"spent"`              //这笔钱是否被花掉
	Value             big.Int                `json:"value"`              //交易金额
	SpendingOutpoints []SpendingOutpointsBtc `json:"spending_outpoints"` //被花费BTC输出列表
	N                 int                    `json:"n"`                  //交易输出序号
	TxIndex           big.Int                `json:"tx_index"`           //交易索引
	Script            string                 `json:"script"`             //目标公钥脚本
	Addr              string                 `json:"addr"`               //目标账户地址
}
type PrevOutBtc struct {
	Addr              string                 `json:"addr"`               //交易账户地址
	N                 int                    `json:"n"`                  //索引编号
	Script            string                 `json:"script"`             //交易输出脚本
	SpendingOutpoints []SpendingOutpointsBtc `json:"spending_outpoints"` //被花费BTC输出列表
	Spent             bool                   `json:"spent"`              //这笔钱是否被花过
	TxIndex           big.Int                `json:"tx_index"`           //交易索引号
	Type              int                    `json:"type"`
	Value             big.Int                `json:"value"` //交易金额-单位为聪
}
type SpendingOutpointsBtc struct {
	N       int     `json:"n"`        //输出索引
	TxIndex big.Int `json:"tx_index"` //交易id
}
