package utils

import (
	"errors"
	"fmt"
	"log"
	"lyods-adsTool/config"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// RateLimitedClient 是一个带有请求速率限制的 HTTP 客户端
type RateLimitedClient struct {
	client      *http.Client
	rateLimiter <-chan time.Time
}

// NewRateLimitedClient 创建一个带有请求速率限制的 HTTP 客户端
func NewRateLimitedClient(requestsPerMinute int) *RateLimitedClient {
	rate := time.Minute / time.Duration(requestsPerMinute)
	rateLimiter := time.Tick(rate)
	return &RateLimitedClient{
		client:      &http.Client{},
		rateLimiter: rateLimiter,
	}
}

// SendRequest 发送 HTTP 请求
func (rlc *RateLimitedClient) SendRequest(url, method string, body []byte) (*http.Response, error) {
	// 等待 rateLimiter 发送信号，以确保不超过请求速率
	<-rlc.rateLimiter

	// 创建 HTTP 请求
	//req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	//if err != nil {
	//	return nil, err
	//}

	// 发送 HTTP 请求
	resp, err := rlc.client.Get(url)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Received status code: " + strconv.Itoa(resp.StatusCode))
	}

	return resp, nil
}

// CreateClient 获取client对象
func CreateClient() *http.Client {
	//创建代理URL
	proxyURL, err := url.Parse(config.URL_PROXY)
	if err != nil {
		log.Fatal("Error parsing proxy URL: ", err)
		return nil
	}
	// 创建 HTTP Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: transport,
	}
	return client
}

// DoRequest 发送http请求
func DoRequest(req *http.Request) (*http.Response, error) {
	client := CreateClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do Error:", err.Error())
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is " + strconv.Itoa(resp.StatusCode))
	}
	return resp, nil
}

// SendHTTPRequest 根据指定的url发送http请求
func SendHTTPRequest(client *http.Client, url string) (*http.Response, error) {
	//req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	//if err != nil {
	//	log.Println("Create Request Error:", err.Error())
	//	return nil, err
	//}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("Request Error:", err.Error())
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is " + strconv.Itoa(resp.StatusCode))
	}
	return resp, nil
}

// SendHTTPRequestWithRateLimit 根据指定的url发送http请求，处理429状态码
func SendHTTPRequestWithRateLimit(url, method string, body []byte) (*http.Response, error) {
	//随机休眠几秒
	RandomSleep()

	client := CreateClient()

	resp, err := client.Get(url)
	if err != nil {
		log.Println("Do Error:", err.Error())
		return nil, err
	} else if resp.StatusCode == http.StatusTooManyRequests {
		// 429: Too Many Requests
		// 控制重试逻辑，可以增加等待时间，然后重试请求
		maxRetries := 3
		retryInterval := 1 * time.Minute

		for i := 0; i < maxRetries; i++ {
			log.Printf("Received 429: Too Many Requests. Retrying... (Retry %d)\n", i+1)
			time.Sleep(retryInterval)
			resp, err = client.Get(url)
			if err != nil {
				log.Println("Do Error:", err.Error())
				return nil, err
			}

			if resp.StatusCode == http.StatusTooManyRequests {
				// 429: Too Many Requests, retry again
				continue
			} else if resp.StatusCode != http.StatusOK {
				return nil, errors.New("Received status code: " + strconv.Itoa(resp.StatusCode))
			} else {
				// Request succeeded after retries
				return resp, nil
			}
		}
		return nil, fmt.Errorf("max retries reached for 429 status code")
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Received status code: " + strconv.Itoa(resp.StatusCode))
	}
	return resp, nil
}
