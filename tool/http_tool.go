package tool

import (
	"log"
	"lyods-adsTool/config"
	"net/http"
	"net/url"
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
