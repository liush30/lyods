package bitcoin

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"lyods-adsTool/pkg/constants"
	"net/http"
)

// GetLatestBlockNumber 获取最新区块高度
func GetLatestBlockNumber(client http.Client) (int64, error) {
	//发送请求
	res, err := client.Get(constants.REQUEST_BTC_BLOCK)
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
