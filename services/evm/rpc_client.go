package evm

import (
	"bytes"
	"encoding/json"
	"log"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"net/http"
	"time"
)

type ChainBaseClient struct {
	RequestCount    int
	LastRequestTime time.Time
	//httpclient      *http.Client
}

// CheckRequestStatus 查询请求状态,如果超过限制,则等待。1秒内只允许请求10次
func (client *ChainBaseClient) CheckRequestStatus() {
	//如果请求计数超过限制，等待1分钟
	if client.RequestCount >= constants.CHAINBASE_MAX_REQUEST {
		//计算自上一次请求以来的时间间隔
		elapsed := time.Since(client.LastRequestTime)
		//若小于一分钟，我们在 time.Sleep 中等待余下的时间。
		if elapsed < constants.CHAINBASE_PER_MINUTE {
			sleepTime := constants.CHAINBASE_PER_MINUTE - elapsed
			time.Sleep(sleepTime)
		}
		//重置计数器和时间戳
		client.RequestCount = 0
		client.LastRequestTime = time.Now()
	}
}

// SendHTTPRequest 根据指定的url发送http请求，最多尝试三次
func (client *ChainBaseClient) SendHTTPRequest(hash, chain string) (*http.Response, error) {
	client.CheckRequestStatus()
	url := getTraceTransactionUrl(chain)
	log.Println("url:", url)
	payload := map[string]interface{}{
		"id":      constants.RPC_ID,
		"jsonrpc": constants.RPC_VERSION,
		"method":  constants.RPC_METHOD_TRACE,
		"params":  []string{hash},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("Marshal Error:", err.Error())
		return nil, err
	}

	// 创建http请求
	req, err := http.NewRequest(constants.HTTP_POST, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	// 最多尝试三次
	maxRetries := constants.CHAINBASE_RETRY_COUNT
	retryCount := 0
	var resp *http.Response
	var requestErr error

	for retryCount < maxRetries {
		if retryCount > 0 {
			time.Sleep(constants.CHAINBASE_PER_MINUTE)
		}
		resp, requestErr = utils.DoRequest(req)
		if requestErr == nil {
			break // 如果请求成功，跳出循环
		}

		log.Println("Request Error (Retry ", retryCount+1, "):", requestErr.Error())
		retryCount++
	}

	if requestErr != nil {
		log.Println("Max retries reached, returning error.")
		return nil, requestErr
	}

	client.AddReqCount()
	return resp, nil
}

func (client *ChainBaseClient) AddReqCount() {
	client.RequestCount++
}
