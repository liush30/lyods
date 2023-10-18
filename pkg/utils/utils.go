package utils

import (
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// HexToBigInt 将16进制字符串转换为big.Int
func HexToBigInt(hexString string) (*big.Int, error) {
	// 去除 "0x" 前缀
	hexString = strings.TrimPrefix(hexString, "0x")

	// 使用SetString将16进制字符串转换成big.Int
	var result big.Int
	_, success := result.SetString(hexString, 16) // 基数为16，表示16进制
	if !success {
		return nil, fmt.Errorf("invalid hex string: %s", hexString)
	}

	return &result, nil
}

// GenerateTransactionID 生成trace交易的唯一标识符
func GenerateTransactionID(idPrefix string, traceInts []int64) string {
	if len(traceInts) == 0 {
		return idPrefix
	}
	// 获取id
	var stringArray []string
	for _, v := range traceInts {
		stringArray = append(stringArray, strconv.FormatInt(v, 10))
	}
	str := strings.Join(stringArray, "_")
	return idPrefix + str
}
func GenerateUuid() string {
	return uuid.New().String()
}

func JoinInt64SliceToString(arr []int64, delimiter string) string {
	strArr := make([]string, len(arr))
	for i, num := range arr {
		strArr[i] = strconv.FormatInt(num, 10)
	}
	return strings.Join(strArr, delimiter)
}

func RandomSleep() {
	// 生成一个1到5之间的随机数，代表休眠的秒数
	randomSleepSeconds := rand.Intn(5) + 1
	time.Sleep(time.Duration(randomSleepSeconds) * time.Second)
}
