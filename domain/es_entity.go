package domain

import "time"

// es存储结构-实体信息 risk-domain
type DateOfBirth struct {
	DateOfBirth string `json:"dateOfBirth"`
	MainEntry   bool   `json:"mainEntry"`
}

type PlaceOfBirth struct {
	PlaceOfBirth string `json:"placeOfBirth"`
	MainEntry    bool   `json:"mainEntry"`
}

type Nationality struct {
	Country   string `json:"country"`
	MainEntry bool   `json:"mainEntry"`
}

type ID struct {
	IDType         string `json:"idType"`   //ID列表
	IDNumber       string `json:"idNumber"` //
	IDCountry      string `json:"idCountry"`
	ExpirationDate string `json:"expirationDate"`
	IssueDate      string `json:"issueDate"` //签发
}

type OtherInfo struct {
	Type string `json:"type"` //类型
	Info string `json:"info"` //信息
}

type AddressList struct {
	Country         string   `json:"country"`
	StateOrProvince string   `json:"stateOrProvince"`
	City            string   `json:"city"`
	Other           []string `json:"other"`
}

type Entity struct {
	IsIndividual     bool           `json:"isIndividual"`     //是否为个体
	Name             string         `json:"name"`             //名字
	AkaList          []string       `json:"akaList"`          //别名列表
	AddressList      []AddressList  `json:"addressList"`      //地址列表
	DateOfBirthList  []DateOfBirth  `json:"dateOfBirthList"`  //出生日期列表
	PlaceOfBirth     []PlaceOfBirth `json:"placeOfBirth"`     //出生地址列表
	Gender           string         `json:"gender"`           //性别
	Email            []string       `json:"emailList"`        //邮箱列表
	Website          []string       `json:"websiteList"`      //网站列表
	PhoneNumber      []string       `json:"phoneNumberList"`  //电话号码
	IDList           []ID           `json:"idList"`           //ID列表信息
	NationalityList  []Nationality  `json:"nationalityList"`  //国籍列表
	OrganizationType string         `json:"organizationType"` //机构类型
	CitizenshipList  []Nationality  `json:"citizenshipList"`  //公民身份列表
	OrgEstDate       string         `json:"orgEstDate"`       //机构成立日期
	OtherInfo        []OtherInfo    `json:"otherInfo"`        //其他信息
}
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
	Value             int64         `json:"value"`             //交易金额
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
	FromAddr        string `json:"fromAddr"`  //内部交易发起方
	ToAddr          string `json:"toAddr"`    //内部交易接收方
	InputTx         string `json:"inputTx"`   //内部交易输入
	OutputTx        string `json:"outputTx"`  //内部交易输出
	Value           int64  `json:"value"`     //转账金额
	SubTraces       int64  `json:"subtraces"` //子交易个数
	CallType        string `json:"callType"`  //调用类型:call、staticcall（静态调用是一种不会修改合约状态的调用方式，它仅用于查询合约状态而不会进行任何状态变更。）
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

// Erc20Txn 交易中的erc20信息
type Erc20Txn struct {
	//Token        string  `json:"token"`           //token名称 ，若为原生代币，则为空
	//TokenDecimal int     `json:"tokenDecimal"`    //token精确值
	FromAddr     string `json:"fromAddr"`        //内部交易发起方
	ToAddr       string `json:"toAddr"`          //内部交易接收方
	ContractAddr string `json:"contractAddress"` //合约地址
	Amount       int64  `json:"amount"`          //内部交易总金额
}
type InputsTrans struct {
	Sequence int64  `json:"sequence"` //发送者定义的交易版本好
	Witness  string `json:"witness"`  //交易输入内容
	Script   string `json:"script"`   //前序交易输出的目标公钥脚本
	Addr     string `json:"addr"`     //***转入地址（发送方）
	Spent    bool   `json:"spent"`    //这笔钱是否已经被花费
	TxIndex  string `json:"txIndex"`  //交易索引
	Value    int64  `json:"value"`    //转入金额
}
type OutTrans struct {
	Spent   bool   `json:"spent"`   //这笔钱是否花费
	Value   int64  `json:"value"`   //转出金额
	N       int64  `json:"n"`       //交易输出序号
	TxIndex string `json:"txIndex"` //交易索引
	Script  string `json:"script"`  //交易输出脚本
	Addr    string `json:"addr"`    //***转出地址（接收方）
}

// WalletAddr es存储结构-风险名单及风险来源 risk-address
type WalletAddr struct {
	WaAddr      string          `json:"waAddr"`        //风险钱包地址
	EntityId    string          `json:"entityId"`      //entityID
	WaRiskLevel uint            `json:"waRiskLevel"`   //最高风险层级
	WaChain     string          `json:"waChain"`       //所在链
	DsAddr      []AdsDataSource `json:"adsDataSource"` //来源地址
	LevelNumber []Level         `json:"levelNumber"`   //被标记层级信息
	Rules       []string        `json:"rules"`         //被中规则id列表
	IsTrace     bool            `json:"isTrace"`       //是否追查子交易
	IsNeedTrace bool            `json:"isNeedTrace"`   //是否需要追查子交易
}
type AdsDataSource struct {
	DsAddr      string    `json:"dsAddr"`      //来源地址
	DsTransHash []string  `json:"dsTransHash"` //涉及风险交易哈希列表
	DsType      string    `json:"dsType"`      //涉及风险交易哈希列表
	Illustrate  string    `json:"illustrate"`  //风险说明
	Time        time.Time `json:"time"`        //被标记时间
	DsRules     []string  `json:"dsRules"`     //规则id
}
type Level struct {
	Level  int16 `json:"level"`  //所在层级
	Number int16 `json:"number"` //被标记次数
}
