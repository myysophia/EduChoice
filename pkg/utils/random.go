package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateRandomCode 生成指定长度的随机数字验证码
func GenerateRandomCode(length int) string {
	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// 生成随机数
	max := int64(1)
	for i := 0; i < length; i++ {
		max *= 10
	}
	
	// 格式化为固定长度的字符串
	code := strconv.FormatInt(r.Int63n(max), 10)
	
	// 补齐前导零
	for len(code) < length {
		code = "0" + code
	}
	
	return code
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// 生成随机字符串
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	
	return string(b)
} 