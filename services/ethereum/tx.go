package ethereum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lyods-adsTool/pkg/constants"
	"net/http"
)

// GetTraceTransaction 获取指定交易hash的内部交易信息
func GetTraceTransaction(hash string) error {
	url := getTraceTransactionUrl()
	payload := map[string]interface{}{
		"id":      constants.RPC_ID,
		"jsonrpc": constants.RPC_VERSION,
		"method":  constants.RPC_METHOD_TRACE,
		"params":  []string{hash},
	}
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest(constants.HTTP_POST, url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// 处理响应结果
	fmt.Println(string(body))
	return nil
}
