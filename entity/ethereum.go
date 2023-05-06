// Package entity ethereum交易信息结构
package entity

import "math/big"

// TransactionEth 普通交易信息集合
type TransactionEth struct {
	Status  uint      `json:"status"`  //状态码，1为成功
	Message string    `json:"message"` //调用描述信息，OK为成功
	Result  ResultEth `json:"result"`  //交易信息
}

// TransactionInternalEth 内部交易信息集合
type TransactionInternalEth struct {
	Status  uint              `json:"status"`  //状态码，1为成功
	Message string            `json:"message"` //调用描述信息，OK为成功
	Result  ResultInternalEth `json:"result"`  //交易信息
}

// TransactionErc20Eth Erc20转账交易信息集合
type TransactionErc20Eth struct {
	Status  uint           `json:"status"`  //状态码，1为成功
	Message string         `json:"message"` //调用描述信息，OK为成功
	Result  ResultErc20Eth `json:"result"`  //交易信息
}

// ResultEth 普通交易信息结构
type ResultEth struct {
	BlockNumber       uint    `json:"blockNumber"`       //区块编号
	TimeStamp         big.Int `json:"timeStamp"`         //区块时间戳
	Hash              string  `json:"hash"`              //交易哈希
	Nonce             string  `json:"nonce"`             //nonce值
	BlockHash         string  `json:"blockHash"`         //区块哈希值
	TransactionIndex  string  `json:"transactionIndex"`  //交易索引
	From              string  `json:"from"`              //发起账号
	To                string  `json:"to"`                //接收账号
	Value             string  `json:"value"`             //交易金额
	Gas               string  `json:"gas"`               //gas最大使用量
	GasPrice          string  `json:"gasPrice"`          //gas价格
	IsError           bool    `json:"isError"`           //是否发生错误，0表示没有错误，1表示发生错误
	TxreceiptStatus   uint    `json:"txreceipt_status"`  //交易收据状态
	Input             string  `json:"input"`             //交易附加数据，16进制字符串编码
	ContractAddress   string  `json:"contractAddress"`   //合约地址
	CumulativeGasUsed string  `json:"cumulativeGasUsed"` //区块累计交易用量
	GasUsed           string  `json:"gasUsed"`           //本交易的gas用量
	Confirmations     string  `json:"confirmations"`     //交易确认数
	MethodId          string  `json:"methodId"`          //函数签名
	FunctionName      string  `json:"functionName"`      //函数名称
}

// ResultInternalEth 内部交易信息结构
type ResultInternalEth struct {
	BlockNumber     uint    `json:"blockNumber"`     //所在区块号
	TimeStamp       big.Int `json:"timeStamp"`       //时间戳
	Hash            string  `json:"hash"`            //交易哈希
	From            string  `json:"from"`            //发起账号
	To              string  `json:"to"`              //接收账号
	Value           string  `json:"value"`           //交易金额
	ContractAddress string  `json:"contractAddress"` //合约地址
	Input           string  `json:"input"`           //交易附加数据，16进制字符串编码
	Type            string  `json:"type"`            //内部交易类型
	Gas             string  `json:"gas"`             //gas最大使用量
	GasUsed         string  `json:"gasUsed"`         //本交易的gas用量
	TraceId         string  `json:"traceId"`         //跟踪ID
	IsError         bool    `json:"isError"`         //是否发生错误，0表示没有错误，1表示发生错误
	ErrCode         string  `json:"errCode"`         //错误代码
}

// ResultErc20Eth Erc20转账交易信息结构
type ResultErc20Eth struct {
	BlockNumber       uint    `json:"blockNumber"`       //区块编号
	TimeStamp         big.Int `json:"timeStamp"`         //区块时间戳
	Hash              string  `json:"hash"`              //交易哈希
	Nonce             string  `json:"nonce"`             //nonce值
	BlockHash         string  `json:"blockHash"`         //区块哈希值
	From              string  `json:"from"`              //发起账号
	ContractAddress   string  `json:"contractAddress"`   //合约地址
	To                string  `json:"to"`                //接收账号
	Value             string  `json:"value"`             //交易金额
	TokenName         string  `json:"tokenName"`         //代币名称
	TokenSymbol       string  `json:"tokenSymbol"`       //代币符号
	TokenDecimal      string  `json:"tokenDecimal"`      //代币精度值
	TransactionIndex  string  `json:"transactionIndex"`  //交易索引
	Gas               string  `json:"gas"`               //gas最大使用量
	GasPrice          string  `json:"gasPrice"`          //gas价格
	GasUsed           string  `json:"gasUsed"`           //本交易的gas用量
	CumulativeGasUsed string  `json:"cumulativeGasUsed"` //区块累计交易用量
	Input             string  `json:"input"`             //交易附加数据，16进制字符串编码
	Confirmations     string  `json:"confirmations"`     //交易确认数
}
