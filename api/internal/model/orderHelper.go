package model

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

// GenerateOrderNo 生成订单号
// 格式: PREFIX + YYYYMMDDHHMMSS + 6位随机数
// 使用 crypto/rand 生成安全随机数，防止订单号被预测
func GenerateOrderNo(prefix string) string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := secureRandomInt(1000000)
	return fmt.Sprintf("%s%s%06d", prefix, timestamp, random)
}

// secureRandomInt 生成 [0, max) 范围内的安全随机数
func secureRandomInt(max int) int {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		// 极端情况下回退到时间戳纳秒部分
		return int(time.Now().UnixNano() % int64(max))
	}
	return int(binary.BigEndian.Uint64(b[:]) % uint64(max))
}

// GenerateEscortOrderNo 生成护航订单号
func GenerateEscortOrderNo() string {
	return GenerateOrderNo("ESC")
}

// GenerateGiftOrderNo 生成礼物订单号
func GenerateGiftOrderNo() string {
	return GenerateOrderNo("GIFT")
}
