package constants

//api相关常量定义
import "lyods-adsTool/config"

// 根据指定地址查询区块链上的交易信息
const (
	API_ETH_TRANS       = `http://api.etherscan.io/api?module=account&action=txlist&apikey=` + config.ETH_KEY + `&address=` //ETH根据指定地址查询普通交易信息api
	API_ETH_ABI         = `https://api.etherscan.io/api?module=contract&action=getabi&apikey=` + config.ETH_KEY + `&address=`
	API_ETH_INTRANS     = `https://api.etherscan.io/api?module=account&action=txlistinternal&apikey=` + config.ETH_KEY + `&address=` //ETh根据指定地址查询内部交易信息
	API_ETH_USDT        = `https://api.etherscan.io/api?module=account&action=tokentx&apikey=` + config.ETH_KEY + `&contractaddress=` + config.ETH_KEY + `&address=`
	API_DASH_TRANS      = `https://explorer.dash.org/insight-api/txs?address=`                                                      //Dash根据指定地址查询交易信息
	API_BTG_TRANS       = `https://explorer.bitcoingold.org/insight-api/txs?address=`                                               // bitglod根据指定地址查询交易信息
	API_BNB_TRANS       = `https://api.bscscan.com/api?module=account&action=txlist&apikey=` + config.ETH_KEY + `&address=`         //Bnb根据指定地址查询普通交易信息
	API_BNB_INTRANS     = `https://api.bscscan.com/api?module=account&action=txlistinternal&apikey=` + config.ETH_KEY + `&address=` //Bnb根据指定地址查询内部交易
	API_VERGE_ADDRTRANS = `https://verge-blockchain.info/api/address/txs/`                                                          //verge获得指定地址的交易记录
	API_VERGE_TRANS     = `https://verge-blockchain.info/api/tx/`                                                                   //verge获取交易具体信息
	API_ARB_TRANS       = `https://api.arbiscan.io/api?module=account&action=txlist&apikey=` + config.ETH_KEY + `&address=`         //Arbiturm 根据指定地址查询普通交易信息
	API_ARB_INTRANS     = `https://api.arbiscan.io/api?module=account&action=txlistinternal&apikey=YourApiKeyToken&address=`        //根据指定地址查询内部交易信息
)
const (
	// URL_INFRUA = "https://fabled-radial-isle.discover.quiknode.pro/406bcd9a632fae87343b9baf3ef47d664161810b/" //quick Node
	HTTP_ALCHEMY_ETH    = "https://eth-mainnet.g.alchemy.com/v2/THcL8Tk-e_GH4Yjagko7-zOMqRGnk2SL"
	HTTP_OMNIATECH_BSC  = "https://endpoints.omniatech.io/v1/bsc/mainnet/95dfd04a356748118abb5cd0f7958a99" //50/s  10w/day
	HTTP_CHAINBASE_ETH  = "https://ethereum-mainnet.s.chainbase.online/v1/"                                //10s /20个
	HTTP_CHAINBASE_BSC  = "https://bsc-mainnet.s.chainbase.online/v1/"
	HTTP_CHAINNODES_BSC = "https://mainnet.chainnodes.org/dde52609-2f9a-4ea5-9ab2-555760c33468"
)

// bitcoin
const (
	BTC_BLOCK            = "https://chain.api.btc.com/v3/block/latest" //btc.com 获取最新区块
	BTC_ADDR             = "https://chain.api.btc.com/v3/address/"     //一分钟请求6次
	BTC_ADDR_BLOCKCHAIN  = `https://blockchain.info/rawaddr/`          //BTC根据指定地址查询交易信息api
	BTC_LIMIT            = "5000"
	BTC_PAGRSIZE         = "50"
	BTC_INIT_PAGE        = "1"
	BTC_ADDRESS          = "https://chain.api.btc.com/v3/address/" //获取btc指定地址信息
	BTC_BLOCK_BLOCKCHAIN = "https://blockchain.info/rawblock/"
)

// evm
const (
	ETH_KEY1        = "G3VRF1S7IMYSSP3D9J8KKUMJY61XT4EK2Q"
	ETH_KEY2        = "UGE82ZM6XAU15BR5C5JVKSBMRM8DA3TQ4A"
	ETH_MAX_SECOND  = 5
	ETH_MAX_DAY     = 100000
	ETH_START_BLOCK = "0"
	//ETH_ADDR_ETHSCAN = "https://api.etherscan.io/api?module=account&action=txlist&startblock="
	ETH_ADDR_ETHSCAN = "https://api.etherscan.io/api?module=account&action=txlist&page=1&offset=50&startblock="
	ETH_MAX_TRANS    = 10000
	ETH_ABI          = "https://api.etherscan.io/api?module=contract&action=getabi&apikey="
)
const (
	BSC_KEY1       = "VM4NA84WJVKVMMMITVY3VH8RWZB7217C8H"
	BSC_KEY2       = "3RJKPK988FKS13T4AE7RNZAE29AUG8WTCY"
	BSC_MAX_SECOND = 5
	BSC_MAX_DAY    = 100000
	BSC_ENDPOINTS  = "https://api.bscscan.com/api"
	BSC_ABI        = "?module=contract&action=getabi&apikey="
	BSC_TX_ADDR    = "?module=account&action=txlist&page=1&offset=50&startblock="
)
