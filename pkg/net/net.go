package net

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

// GenerateMask 根据CIDR表示的IP地址生成对应的子网掩码
func GenerateMask(cidr string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	ones, _ := ipNet.Mask.Size()
	onesBinary := strings.Repeat("1", ones)
	zerosBinary := strings.Repeat("0", 32-ones)

	maskBinary := onesBinary + zerosBinary
	maskInt, success := new(big.Int).SetString(maskBinary, 2)
	if !success {
		return "", fmt.Errorf("failed to convert binary string to integer")
	}

	// 将整数形式的子网掩码转换为四个字节的形式
	maskBytes := maskInt.Bytes()
	if len(maskBytes) < 4 {
		prefix := make([]byte, 4-len(maskBytes))
		maskBytes = append(prefix, maskBytes...)
	}

	mask := net.IP(maskBytes)
	return mask.String(), nil
}
