package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// EncryptPassword 密码加密
func EncryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
} 