package utils

import (
	"fmt"
	"github.com/google/uuid"
	"log"
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
	log.Println("Sleep", randomSleepSeconds, "s......")
}

// IsTimeAfter 判断是否晚于某个时间
func IsTimeAfter(targetHour, targetMinute int) bool {
	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, 0, 0, now.Location())

	return now.After(targetTime)
}
func DateChange(inputDate string) (string, error) {
	if inputDate == "" {
		return "", nil
	}
	// 使用指定格式解析日期字符串
	parsedDate, err := parseDate(inputDate)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %s", inputDate)
	}
	// 格式化日期为 "yyyy-MM-dd" 格式
	formattedDate := parsedDate.Format("2006-01-02")
	return formattedDate, nil
}
func parseDate(dateString string) (time.Time, error) {
	// 尝试不同的日期解析格式
	formats := []string{"2006", "Jan 2006", "02 Jan 2006", "02 Jan 2006 to 02 Jan 2006"}
	var parsedDate time.Time
	var err error
	for _, format := range formats {
		parsedDate, err = time.Parse(format, dateString)
		if err == nil {
			break
		}
	}
	return parsedDate, err
}

func UnixToTime(unixTime int64) string {
	// 将时间戳转换为时间对象
	utcTime := time.Unix(unixTime, 0)

	// 指定时区为 UTC
	utcLocation := time.UTC
	utcTime = utcTime.In(utcLocation)

	// 格式化时间为字符串
	utcFormattedTime := utcTime.Format(time.DateTime)
	return utcFormattedTime
}
