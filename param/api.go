package param

// 根据指定地址查询区块链上的交易信息
const (
	API_BTC_TRANS       = `https://blockchain.info/rawaddr/`                                                                  //BTC根据指定地址查询交易信息api
	API_ETH_TRANS       = `http://api.etherscan.io/api?module=account&action=txlist&apikey=` + ETH_KEY + `&address=`          //ETH根据指定地址查询普通交易信息api
	API_ETH_INTRANS     = `https://api.etherscan.io/api?module=account&action=txlistinternal&apikey=` + ETH_KEY + `&address=` //ETh根据指定地址查询内部交易信息
	API_ETH_USDT        = `https://api.etherscan.io/api?module=account&action=tokentx&apikey=` + ETH_KEY + `&contractaddress=` + ETH_USDT + `&address=`
	API_DASH_TRANS      = `https://explorer.dash.org/insight-api/txs?address=`                                               //Dash根据指定地址查询交易信息
	API_BTG_TRANS       = `https://explorer.bitcoingold.org/insight-api/txs?address=`                                        // bitglod根据指定地址查询交易信息
	API_BNB_TRANS       = `https://api.bscscan.com/api?module=account&action=txlist&apikey=` + BNB_KEY + `&address=`         //Bnb根据指定地址查询普通交易信息
	API_BNB_INTRANS     = `https://api.bscscan.com/api?module=account&action=txlistinternal&apikey=` + BNB_KEY + `&address=` //Bnb根据指定地址查询内部交易
	API_VERGE_ADDRTRANS = `https://verge-blockchain.info/api/address/txs/`                                                   //verge获得指定地址的交易记录
	API_VERGE_TRANS     = `https://verge-blockchain.info/api/tx/`                                                            //verge获取交易具体信息
	API_ARB_TRANS       = `https://api.arbiscan.io/api?module=account&action=txlist&apikey=` + ARB_KEY + `&address=`         //Arbiturm 根据指定地址查询普通交易信息
	API_ARB_INTRANS     = `https://api.arbiscan.io/api?module=account&action=txlistinternal&apikey=YourApiKeyToken&address=` //根据指定地址查询内部交易信息
)
