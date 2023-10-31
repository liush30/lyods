package bitcoin

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"lyods-adsTool/pkg/constants"
	"math/big"
	"net/http"
)

// GetLatestBlockNumber 获取最新区块高度
func GetLatestBlockNumber(client http.Client) (int64, error) {
	//发送请求
	res, err := client.Get(constants.BTC_BLOCK)
	if err != nil {
		return 0, err
	} else if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code is %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	number, err := jsonparser.GetInt(body, "data", "height")
	if err != nil {
		return 0, err
	}
	return number, nil
}

// ConvertSatoshiToBTC 将给定的 Satoshi 转换为 BTC
func ConvertSatoshiToBTC(satoshi *big.Int) (float64, string) {
	// 换算关系：1 BTC = 100,000,000 Satoshi
	btcFloat := new(big.Float).Quo(new(big.Float).SetInt(satoshi), new(big.Float).SetInt64(100000000))
	btcString := btcFloat.Text('f', 8) // 以字符串形式表示，精确到8位小数

	btcValue, _ := btcFloat.Float64() // 转换为 float64

	return btcValue, btcString
}
