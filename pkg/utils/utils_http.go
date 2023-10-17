package utils

import (
	"bytes"
	"errors"
	"log"
	"lyods-adsTool/config"
	"net/http"
	"net/url"
	"strconv"
)

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
func SendHTTPRequest(url, method string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Create Request Error:", err.Error())
		return nil, err
	}
	resp, err := DoRequest(req)
	if err != nil {
		log.Println("Request Error:", err.Error())
		return nil, err
	}
	return resp, nil
}
