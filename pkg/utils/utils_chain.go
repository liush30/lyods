package utils

import (
	"log"
	"lyods-adsTool/pkg"
)

// IsExitsAddrSource 判断地址中的地址来源是否已经存在
//func IsExitsAddrSource(sourcesList []domain.AdsDataSource, parentAddress string) (bool, int) {
//	for i, v := range sourcesList {
//		//判断来源地址是否与父地址一致，若一致返回索引值
//		if strings.EqualFold(v.DsAddr, parentAddress) {
//			return true, i
//		}
//	}
//	return false, -1
//}

// ErrorLog 判断是否存在错误，若存在则打印日志
func ErrorLog(err error, str1 string) {
	if err != nil {
		log.Printf("%s:%s\n", str1, err.Error())
	}
}
func GetChain(name string) string {
	return pkg.CurrencyToChain[name]
}