package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GetTransid 生成交易流水号
func GetTransid(prefix string) string {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 生成交易流水号，格式：前缀 + 时间戳 + 随机数
	timestamp := time.Now().Format("0102150405")
	randomNum := rand.Intn(10000)

	return fmt.Sprintf("%s%s%04d", prefix, timestamp, randomNum)
}

// GetTimeFormat 获取指定格式的时间字符串
// format: YmdHis - 年月日时分秒，Ymd - 年月日
func GetTimeFormat(format string) string {
	now := time.Now()

	switch format {
	case "YmdHis":
		return now.Format("20060102150405")
	case "Ymd":
		return now.Format("20060102")
	case "Y-m-d H:i:s":
		return now.Format("2006-01-02 15:04:05")
	case "Y-m-d":
		return now.Format("2006-01-02")
	default:
		return now.Format("20060102150405")
	}
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i := range bytes {
		bytes[i] = chars[rand.Intn(len(chars))]
	}

	return string(bytes)
}
