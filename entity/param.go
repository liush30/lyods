package entity

const (
	EthereumApiKey  = `G3VRF1S7IMYSSP3D9J8KKUMJY61XT4EK2Q` //以太坊API KEY
	ApiBitcoinTrans = `https://blockchain.info/rawaddr/`   //BTC根据指定地址查询交易信息api
	ApiEthTrans     = `https://api.etherscan.io/api?module=account&action=txlist&apikey=` + EthereumApiKey + `&address=`
)
