package pre

import (
	"log"
	"os"
)

// GetContractAddressToDb 解析本地csv文件，获取合约地址
func GetContractAddressToDb(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Failed to open CSV file:", err)
		return err
	}
	defer file.Close()

}
