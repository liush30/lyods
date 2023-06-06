package param

const (
	HTTP_GET   = `GET`
	MAX_RECODE = "100000"
	ETH_KEY    = `G3VRF1S7IMYSSP3D9J8KKUMJY61XT4EK2Q`         //Etherscan API KEY
	BNB_KEY    = `3RJKPK988FKS13T4AE7RNZAE29AUG8WTCY`         //BNB API KEY
	ARB_KEY    = `D8HNF3S19ECJD3T3M6IFMY57159I866A1J`         //ARB API KEY
	BSV_KEY    = `mainnet_02b6944ba19ef7638152d029dffe926f`   //BSV TAAL KEY
	URL_PROXY  = `http://127.0.0.1:7890`                      //Etherscan 代理
	ETH_USDT   = `0xdac17f958d2ee523a2206206994597c13d831ec7` //以太坊中的USDT合约地址
	ETH_USDC   = `0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48` //以太坊中的USDC合约地址
)
const (
	CHAIN_DASH     = iota //所在链-DASH
	CHAIN_BITGOLD         //所在链-BITGOLD
	CHAIN_BNB             //所在链-BNB
	CHAIN_ETHEREUM        //所在链-ETHEREUM
	CHAIN_ARBITRUM        //所在链-ARBITRUM
)
const (
	DS_TYPE_URL = "1" //来源类型-url
)
const (
	INIT_NUMBER = 0 //初始标记次数
)
