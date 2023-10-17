package constants

const (
	HTTP_GET  = "get"
	HTTP_POST = "POST"
)
const (
	CHAIN_DASH     = iota //所在链-DASH
	CHAIN_BITGOLD         //所在链-BITGOLD
	CHAIN_BNB             //所在链-BNB
	CHAIN_ETHEREUM        //所在链-ETHEREUM
	CHAIN_ARBITRUM        //所在链-ARBITRUM
	CHAIN_BITCOIN
)
const (
	CHAIN_ETH = "ETH"
)

// 风险名单来源类型
const (
	DS_OPENSANCTIONS = "opensanctions" //来源类型-OPENSANCTIONS
	DS_OFAC          = "ofac"          //来源类型-OFAC
	DS_UNISWAP       = "uniswap"       //来源类型-UNISWAP
	DS_CUSTOMIZE     = "customize"     //来源类型-平台自定义
	DS_WITH_RISK     = "with_risk"     //来源类型-涉及与风险地址交易
)
const (
	DSADDR_SDN = "https://www.treasury.gov/ofac/downloads/sdn.xml"
)

// 初始化数据
const (
	INIT_LEVEL = 1
)

// es相关常量定义信息
const (
	ES_ADDRESS     = "risk_address"     //风险名单索引名称
	ES_TRANSACTION = "risk_transaction" //风险交易索引名称
	ES_ENTITY      = "risk_entity"      //风险名单实体信息索引名称
)
const (
	RPC_ID           = 1
	RPC_VERSION      = "2.0"
	RPC_METHOD_TRACE = "trace_transaction"
)
const (
	ABI_NO = "Contract source code not verified"
)
