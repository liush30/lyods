package tool

var tickerMap map[string]string

func init() {
	tickerMap = make(map[string]string)
	tickerMap["ETH"] = "ethereum"
	tickerMap["BTC"] = "bitcoin"

}

// GetChainByTicker 根据货币代码查询所在链
func GetChainByTicker(ticker string) string {
	return tickerMap[ticker]
}
