package risklist

//func TestGetAddrListByJSONOnBitcoin_Positive(t *testing.T) {
//	// 调用被测试函数
//	err := list.GetAddrListByJSONOnBitcoin("https://data.opensanctions.org/datasets/20230927/ransomwhere/statistics.json")
//	if err != nil {
//		t.Errorf("Expected no error:%v", err)
//	}
//}
//func TestGetAddrListByJSONOnBitcoin_Negative(t *testing.T) {
//	// 创建一个模拟的HTTP服务器
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// 返回一个错误的HTTP状态码
//		w.WriteHeader(http.StatusInternalServerError)
//	}))
//	defer server.Close()
//
//	// 调用被测试函数
//	err := list.GetAddrListByJSONOnBitcoin(server.URL)
//	if err == nil {
//		t.Error("Expected an error, got nil")
//	}
//}
//func TestGetAddrListOnCsv_Positive(t *testing.T) {
//	// 调用被测试函数
//	err := list.GetAddrListOnCsv("http://gist.githubusercontent.com/banteg/1657d4778eb86c460e03bc58b99970c0/raw/2b8e0b2c1074b995b992397f34ab2843cf6bdf72/uniswap-trm.csv")
//	if err != nil {
//		t.Errorf("GetAddrListOnCsv returned an error: %v", err)
//	}
//
//	// TODO: 添加断言来验证函数的输出或副作用
//}
//
//func TestGetAddrListOnCsv_Negative(t *testing.T) {
//	// 创建一个模拟的HTTP响应
//	handler := func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//	server := httptest.NewServer(http.HandlerFunc(handler))
//	defer server.Close()
//
//	// 调用被测试函数
//	err := list.GetAddrListOnCsv(server.URL)
//	if err == nil {
//		t.Errorf("GetAddrListOnCsv did not return an error")
//	}
//
//	// TODO: 添加断言来验证函数的输出或副作用
//}
//
//func TestGetAddrListOnXmlByElement(t *testing.T) {
//	// Positive test case
//	url := "https://www.treasury.gov/ofac/downloads/sdn.xml"
//	err := list.GetAddrListOnXmlByElement(url)
//	if err != nil {
//		t.Errorf("Expected no error, but got %v", err)
//	}
//}
