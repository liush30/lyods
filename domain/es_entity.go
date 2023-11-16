package domain

import "reflect"

// DateOfBirth es存储结构-实体信息 risk-domain
type DateOfBirth struct {
	DateOfBirth string `json:"dateOfBirth"`
	MainEntry   bool   `json:"mainEntry"`
}
type Rule struct {
	RuKey     string `json:"ruKey,dsRuKey"`         //规则id
	RuType    string `json:"ruType,dsRuType"`       //规则类型
	RuCode    string `json:"ruCode,dsRuCode"`       //规则代码
	RuDesc    string `json:"ruDesc,dsRuDesc"`       //规则名称
	Status    string `json:"status,dsRuStatus"`     //规则状态
	RuExpress string `json:"ruExpress,dsRuExpress"` //规则表达式
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
	IDType         string      `json:"idType"`   //ID列表
	IDNumber       string      `json:"idNumber"` //
	IDCountry      string      `json:"idCountry"`
	ExpirationDate interface{} `json:"expirationDate"`
	IssueDate      interface{} `json:"issueDate"` //签发
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
	EntityId         string             `json:"entityId"`         //实体id
	IsIndividual     bool               `json:"isIndividual"`     //是否为个体
	RiskLevel        int64              `json:"riskLevel"`        //风险等级
	LevelTime        string             `json:"levelTime"`        //最新风险层级标记时间
	Name             string             `json:"name"`             //名字
	AkaList          []string           `json:"akaList"`          //别名列表
	AddressList      []AddressList      `json:"addressList"`      //地址列表
	DateOfBirthList  []DateOfBirth      `json:"dateOfBirthList"`  //出生日期列表
	PlaceOfBirthList []PlaceOfBirth     `json:"placeOfBirthList"` //出生地址列表
	Gender           string             `json:"gender"`           //性别
	Email            []string           `json:"emailList"`        //邮箱列表
	Website          []string           `json:"websiteList"`      //网站列表
	PhoneNumber      []string           `json:"phoneNumberList"`  //电话号码
	IDList           []ID               `json:"idList"`           //ID列表信息
	NationalityList  []Nationality      `json:"nationalityList"`  //国籍列表
	OrganizationType string             `json:"organizationType"` //机构类型
	CitizenshipList  []Nationality      `json:"citizenshipList"`  //公民身份列表
	OrgEstDate       interface{}        `json:"orgEstDate"`       //机构成立日期
	OtherInfo        []OtherInfo        `json:"otherInfo"`        //其他信息
	RiskChgHistory   []RiskChangeRecord `json:"riskChgHistory"`   //风险变更记录
	Rules            []Rule             `json:"rules"`            //被中规则id列表
}

type EsTrans struct {
	Hash              string             `json:"hash"`              //***交易哈希
	Chain             string             `json:"chain"`             //交易所属链
	InputCount        int64              `json:"inputCount"`        //交易输入数量
	InputValue        float64            `json:"inputValue"`        //交易输入金额
	OutputCount       int64              `json:"outputCount"`       //交易输出数量
	OutputValue       float64            `json:"outputValue"`       //交易输出金额
	RiskLevel         int64              `json:"riskLevel"`         //风险等级
	AddressList       []string           `json:"addressList"`       //***交易所属地址
	AddressListId     []string           `json:"addressIdList"`     //交易所属风险地址id
	Balance           float64            `json:"balance"`           //交易后账户的 余额
	Size              int64              `json:"size"`              //交易字节数
	Weight            int64              `json:"weight"`            //权重
	GasUsed           string             `json:"gasUsed"`           //交易费用||gas使用量
	IsError           string             `json:"isError"`           //是否发生错误，0表示没有错误，1表示发生错误
	ErrCode           string             `json:"errCode"`           //错误代码
	ContractAddress   string             `json:"contractAddress"`   //合约地址
	FunctionName      string             `json:"functionName"`      //方法名称
	MethodId          string             `json:"methodId"`          //函数签名
	TraceId           string             `json:"traceId"`           //跟踪ID
	Confirmations     string             `json:"confirmations"`     //交易确认数
	CumulativeGasUsed string             `json:"cumulativeGasUsed"` //区块累计交易用量
	GasPrice          string             `json:"gasPrice"`          //gas价格
	LockTime          int64              `json:"lockTime"`          //锁定时间
	TxIndex           string             `json:"txIndex"`           //交易索引
	DoubleSpend       bool               `json:"doubleSpend"`       //是否双花
	Time              string             `json:"time"`              //交易时间
	BlockHeight       string             `json:"blockHeight"`       //区块高度
	BlockHash         string             `json:"blockHash"`         //区块哈希
	Value             float64            `json:"value"`             //交易金额
	ValueText         string             `json:"valueText"`         //交易金额-text
	ValueUSD          float64            `json:"valueUSD"`          //交易转换的美元价值
	Inputs            []InputsTrans      `json:"inputs"`            //输入信息（发送方）
	Out               []OutTrans         `json:"out"`               //输出信息（接收方）
	InternalTx        []InternalTxn      `json:"internalTx"`        //交易中的内部交易信息
	Logs              []Logs             `json:"logs"`              //交易中的日志信息
	Erc20Txn          []Erc20Txn         `json:"erc20Txn"`          //交易中的erc20信息
	RiskChgHistory    []RiskChangeRecord `json:"riskChgHistory"`    //风险变更记录
	Rules             []Rule             `json:"rules"`             //被中规则id列表
}

// InternalTxn 交易内部信息
type InternalTxn struct {
	Id              string  `json:"id"`           //ID
	TraceAddress    string  `json:"traceAddress"` //路径
	TraceAddressInt []int64 //临时存储traceAddress
	Init            string  `json:"init"`      //合同初始化字节码
	Address         string  `json:"address"`   //如果这是一个合约初始化交易，则表示新合约的地址，否则为空
	Code            string  `json:"code"`      //合约字节码文件。如果这是一个合约初始化交易，则包含合约的字节码，否则为空
	Type            string  `json:"type"`      //交易类型 create：表示这是一个合约初始化交易  call表示为一个普通 合约调用交易
	FromAddr        string  `json:"fromAddr"`  //内部交易发起方
	ToAddr          string  `json:"toAddr"`    //内部交易接收方
	InputTx         string  `json:"inputTx"`   //内部交易输入
	OutputTx        string  `json:"outputTx"`  //内部交易输出
	Value           float64 `json:"value"`     //转账金额
	ValueText       string  `json:"valueText"` //转账金额-text
	SubTraces       int64   `json:"subtraces"` //子交易个数
	CallType        string  `json:"callType"`  //合约调用类型:call、staticcall（静态调用是一种不会修改合约状态的调用方式，它仅用于查询合约状态而不会进行任何状态变更。）
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
	FromAddr     string  `json:"fromAddr"`        //内部交易发起方
	ToAddr       string  `json:"toAddr"`          //内部交易接收方
	ContractAddr string  `json:"contractAddress"` //合约地址
	Amount       float64 `json:"amount"`          //内部交易总金额
	AmountText   string  `json:"amountText"`      //内部交易总金额-text
}
type InputsTrans struct {
	Sequence  int64   `json:"sequence"`  //发送者定义的交易版本好
	Witness   string  `json:"witness"`   //交易输入内容
	Script    string  `json:"script"`    //前序交易输出的目标公钥脚本
	Addr      string  `json:"addr"`      //***转入地址（发送方）
	Spent     bool    `json:"spent"`     //这笔钱是否已经被花费
	TxIndex   string  `json:"txIndex"`   //交易索引
	Value     float64 `json:"value"`     //转入金额
	ValueText string  `json:"valueText"` //转入金额-text
}
type OutTrans struct {
	Spent     bool    `json:"spent"`     //这笔钱是否花费
	Value     float64 `json:"value"`     //转出金额
	ValueText string  `json:"valueText"` //转出金额-text
	N         int64   `json:"n"`         //交易输出序号
	TxIndex   string  `json:"txIndex"`   //交易索引
	Script    string  `json:"script"`    //交易输出脚本  //上一个交易哈希值
	Addr      string  `json:"addr"`      //***转出地址（接收方）
}

// WalletAddr es存储结构-风险名单及风险来源 risk-address
type WalletAddr struct {
	AddressId      string             `json:"addressId"`      //风险地址id
	WaAddr         string             `json:"waAddr"`         //风险钱包地址
	EntityId       string             `json:"entityId"`       //entityID
	Balance        float64            `json:"balance"`        //地址余额
	WaRiskLevel    uint               `json:"waRiskLevel"`    //最高风险层级
	LevelTime      string             `json:"levelTime"`      //最新风险层级标记时间
	WaChain        string             `json:"waChain"`        //所在链
	DsAddr         []AdsDataSource    `json:"adsDataSource"`  //来源地址
	LevelNumber    []Level            `json:"levelNumber"`    //被标记层级信息
	IsTrace        bool               `json:"isTrace"`        //是否追查子交易
	IsNeedTrace    bool               `json:"isNeedTrace"`    //是否需要追查子交易
	IsContract     bool               `json:"isContract"`     //是否是合约地址
	RiskChgHistory []RiskChangeRecord `json:"riskChgHistory"` //风险变更记录
	Rules          []Rule             `json:"rules"`          //被中规则id列表
}

type AdsDataSource struct {
	DsAddr      string   `json:"dsAddr"`      //来源地址
	DsTransHash []string `json:"dsTransHash"` //涉及风险交易哈希列表
	DsType      string   `json:"dsType"`      //涉及风险交易哈希列表
	Illustrate  string   `json:"illustrate"`  //风险说明
	Time        string   `json:"time"`        //被标记时间
	DsRules     []Rule   `json:"dsRules"`     //规则id
}
type Level struct {
	Level  int16 `json:"level"`  //所在层级
	Number int16 `json:"number"` //被标记次数
}
type RiskChangeRecord struct {
	DateOfChange string `json:"dateOfChange"` //变更日期
	RiskLevel    uint   `json:"riskLevel"`    //变更风险等级
	Description  string `json:"description"`  //变更描述
}

// IsEsTransEmpty 判断 EsTrans 结构体是否为空
func IsEsTransEmpty(esTrans EsTrans) bool {
	// 获取结构体的反射值
	val := reflect.ValueOf(esTrans)

	// 遍历结构体的字段
	for i := 0; i < val.NumField(); i++ {
		// 获取字段的值
		fieldVal := val.Field(i)

		// 判断字段的类型
		switch fieldVal.Kind() {
		case reflect.String:
			// 如果是字符串类型，判断是否为空字符串
			if fieldVal.String() != "" {
				return false
			}
		case reflect.Slice, reflect.Array, reflect.Map:
			// 如果是切片、数组或映射类型，判断是否为空
			if !fieldVal.IsNil() && fieldVal.Len() > 0 {
				return false
			}
		case reflect.Int, reflect.Int64, reflect.Float64:
			// 如果是整数或浮点数类型，判断是否为零值
			if fieldVal.Interface() != reflect.Zero(fieldVal.Type()).Interface() {
				return false
			}
		case reflect.Bool:
			// 如果是布尔类型，判断是否为 false
			if fieldVal.Bool() {
				return false
			}
		default:
			// 暂时不处理其他类型
		}
	}

	// 所有字段都为空，则结构体为空
	return true
}
