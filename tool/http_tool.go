package tool

import (
	"log"
	"lyods-adsTool/param"
	"net/http"
	"net/url"
)

var MClient *http.Client

func init() {
	MClient = createClient()
}

// 获取client对象
func createClient() *http.Client {
	//创建代理URL
	proxyURL, err := url.Parse(param.URL_PROXY)
	if err != nil {
		log.Println("Error parsing proxy URL: ", err)
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
