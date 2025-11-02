package model

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOrderNo 生成订单号
// 格式: PREFIX + YYYYMMDDHHMMSS + 6位随机数
func GenerateOrderNo(prefix string) string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := rand.Intn(1000000)
	return fmt.Sprintf("%s%s%06d", prefix, timestamp, random)
}

// GenerateEscortOrderNo 生成护航订单号
func GenerateEscortOrderNo() string {
	return GenerateOrderNo("ESC")
}

// GenerateGiftOrderNo 生成礼物订单号
func GenerateGiftOrderNo() string {
	return GenerateOrderNo("GIFT")
}

