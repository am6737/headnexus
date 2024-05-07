package code

import (
	"encoding/hex"
	"math/rand"
)

var bytes []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

func GenEnrollCode() string {
	return randStr1(43)
}

func randStr1(n int) string {
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = bytes[rand.Int31()%62]
	}
	return string(result)
}

func randStr2(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}
