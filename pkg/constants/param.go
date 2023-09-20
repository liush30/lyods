package constants

const (
	HTTP_GET = "get"
)
const (
	CHAIN_DASH     = iota //所在链-DASH
	CHAIN_BITGOLD         //所在链-BITGOLD
	CHAIN_BNB             //所在链-BNB
	CHAIN_ETHEREUM        //所在链-ETHEREUM
	CHAIN_ARBITRUM        //所在链-ARBITRUM
	CHAIN_BITCOIN
)

// 风险名单来源类型
const (
	DS_TYPE_URL = iota //来源类型-普通url
	DS_OFAC
	DS_TYPE_ADDRESS //来源类型-地址
)

// 初始化数据
const (
	INIT_NUMBER = 1 //初始标记次数
	INIT_LEVEL  = 1
)

// es相关常量定义信息
const (
	ES_ADDRESS     = "risk_address"     //风险名单索引名称
	ES_TRANSACTION = "risk_transaction" //风险交易索引名称
	ES_ENTITY      = "risk_entity"      //风险名单实体信息索引名称
)
