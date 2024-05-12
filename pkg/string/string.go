package string

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// string进行md5加密
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}

// GenerateRandomCode 生成随机6位数
func GenerateRandomCode() string {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机数
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}
