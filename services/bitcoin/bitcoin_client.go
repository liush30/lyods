package bitcoin

import (
	"fmt"
	"lyods-adsTool/pkg/constants"
	"net/http"
	"strconv"
	"time"
)

type BitClient struct {
	RequestCount    int
	LastRequestTime time.Time
	Httpclient      *http.Client
}

// CheckRequestStatus 查询请求状态,如果超过限制,则等待1分钟，一分钟内只允许请求6次
func (client *BitClient) CheckRequestStatus() {
	//如果请求计数超过限制，等待1分钟
	if client.RequestCount >= constants.BTC_MAX_REQUEST {
		//计算自上一次请求以来的时间间隔
		elapsed := time.Since(client.LastRequestTime)
		//若小于一分钟，我们在 time.Sleep 中等待余下的时间。
		if elapsed < constants.BTC_PER_MINUTE {
			sleepTime := constants.BTC_PER_MINUTE - elapsed
			time.Sleep(sleepTime)
		}
		//重置计数器和时间戳
		client.RequestCount = 0
		client.LastRequestTime = time.Now()
	}
}

// SendHTTPRequest 根据指定的url发送http请求
//
//	func (client *BitClient) SendHTTPRequest(url string) (*http.Response, error) {
//		client.CheckRequestStatus()
//		resp, err := client.Httpclient.Get(url)
//		if err != nil {
//			return nil, fmt.Errorf("send http request error:%v", err)
//		} else if resp.StatusCode != http.StatusOK {
//			return nil, fmt.Errorf("status code is %s", strconv.Itoa(resp.StatusCode))
//		}
//		client.AddReqCount()
//		return resp, nil
//	}
func (client *BitClient) SendHTTPRequest(url string) (*http.Response, error) {
	for retry := 0; retry < 3; retry++ {
		client.CheckRequestStatus()
		resp, err := client.Httpclient.Get(url)
		if err != nil {
			if retry < 2 {
				// 如果请求失败，等待1分钟后重试
				time.Sleep(1 * time.Minute)
			} else {
				return nil, fmt.Errorf("send http request error after 3 retries: %v", err)
			}
		} else if resp.StatusCode == http.StatusOK {
			client.AddReqCount()
			return resp, nil
		} else {
			if retry < 2 {
				// 如果状态码不是200，等待1分钟后重试
				time.Sleep(1 * time.Minute)
			} else {
				return nil, fmt.Errorf("status code is %s after 3 retries", strconv.Itoa(resp.StatusCode))
			}
		}
	}

	return nil, fmt.Errorf("maximum number of retries reached")
}

func (client *BitClient) AddReqCount() {
	client.RequestCount++
}
