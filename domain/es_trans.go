package domain

//es存储结构-风险交易信息 risk-transaction
import (
	"math/big"
)

type EsTrans struct {
	Hash              string        `json:"hash"`              //***交易哈希
	Address           string        `json:"address"`           //***交易所属地址
	Size              int64         `json:"size"`              //交易字节数
	Weight            int64         `json:"weight"`            //权重
	GasUsed           string        `json:"gasUsed"`           //交易费用||gas使用量
	IsError           string        `json:"isError"`           //是否发生错误，0表示没有错误，1表示发生错误
	ErrCode           string        `json:"errCode"`           //错误代码
	ContractAddress   string        `json:"contractAddress"`   //合约地址
	FunctionName      string        `json:"functionName"`      //方法名称
	MethodId          string        `json:"methodId"`          //函数签名
	TraceId           string        `json:"traceId"`           //跟踪ID
	Confirmations     string        `json:"confirmations"`     //交易确认数
	CumulativeGasUsed string        `json:"cumulativeGasUsed"` //区块累计交易用量
	GasPrice          string        `json:"gasPrice"`          //gas价格
	LockTime          int64         `json:"lockTime"`          //锁定时间
	TxIndex           string        `json:"txIndex"`           //交易索引
	DoubleSpend       bool          `json:"doubleSpend"`       //是否双花
	Time              int64         `json:"time"`              //交易时间
	BlockHeight       string        `json:"blockHeight"`       //区块高度
	BlockHash         string        `json:"blockHash"`         //区块哈希
	Value             big.Int       `json:"value"`             //交易金额
	ValueUSD          float64       `json:"valueUSD"`          //交易转换的美元价值
	Inputs            []InputsTrans `json:"inputs"`            //输入信息（发送方）
	Out               []OutTrans    `json:"out"`               //输出信息（接收方）
	InternalTx        []InternalTxn `json:"internalTx"`        //交易中的内部交易信息
	Logs              []Logs        `json:"logs"`              //交易中的日志信息
	Erc20Txn          Erc20Txn      `json:"erc20Txn"`
}

// InternalTxn 交易内部信息
type InternalTxn struct {
	Id              string `json:"id"`           //ID
	TraceAddress    string `json:"traceAddress"` //路径
	TraceAddressInt []int64
	FromAddr        string  `json:"fromAddr"`  //内部交易发起方
	ToAddr          string  `json:"toAddr"`    //内部交易接收方
	InputTx         string  `json:"inputTx"`   //内部交易输入
	OutputTx        string  `json:"outputTx"`  //内部交易输出
	Value           big.Int `json:"value"`     //转账金额
	SubTraces       int64   `json:"subtraces"` //子交易个数
	CallType        string  `json:"callType"`  //调用类型:call、staticcall（静态调用是一种不会修改合约状态的调用方式，它仅用于查询合约状态而不会进行任何状态变更。）
	//ActualAmount big.Int `json:"actualAmount"` //实际金额
	//AmountUSD    float64 `json:"amountUSD"`       //内部交易转化的美元价值
	//ActualUSD    float64 `json:"actualUSD"`       //实际转化的美元价值
}

// Logs 交易中的日志信息
type Logs struct {
	Address   string            `json:"address"`
	EventInfo string            `json:"eventInfo"`
	Topics    []TopicsValStruct `json:"topics"`
}
type TopicsValStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// erc20
type Erc20Txn struct {
	//Token        string  `json:"token"`           //token名称 ，若为原生代币，则为空
	//TokenDecimal int     `json:"tokenDecimal"`    //token精确值
	FromAddr     string  `json:"fromAddr"`        //内部交易发起方
	ToAddr       string  `json:"toAddr"`          //内部交易接收方
	ContractAddr string  `json:"contractAddress"` //合约地址
	Amount       big.Int `json:"amount"`          //内部交易总金额
}
type InputsTrans struct {
	Sequence int64   `json:"sequence"` //发送者定义的交易版本好
	Witness  string  `json:"witness"`  //交易输入内容
	Script   string  `json:"script"`   //前序交易输出的目标公钥脚本
	Addr     string  `json:"addr"`     //***转入地址（发送方）
	Spent    bool    `json:"spent"`    //这笔钱是否已经被花费
	TxIndex  string  `json:"txIndex"`  //交易索引
	Value    big.Int `json:"value"`    //转入金额
}
type OutTrans struct {
	Spent   bool    `json:"spent"`   //这笔钱是否花费
	Value   big.Int `json:"value"`   //转出金额
	N       int64   `json:"n"`       //交易输出序号
	TxIndex string  `json:"txIndex"` //交易索引
	Script  string  `json:"script"`  //交易输出脚本
	Addr    string  `json:"addr"`    //***转出地址（接收方）
}
