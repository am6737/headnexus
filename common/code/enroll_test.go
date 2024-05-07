package code

import (
	"fmt"
	"testing"
)

func TestGenEnrollCode(t *testing.T) {
	enrollCode := GenEnrollCode()
	fmt.Println("Generated enroll code:", enrollCode)

	// 检查生成的代码长度是否符合预期
	expectedLength := 43
	if len(enrollCode) != expectedLength {
		t.Errorf("Expected enroll code length %d, but got %d", expectedLength, len(enrollCode))
	}

	// 检查生成的代码是否只包含指定的字符集合
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	for _, char := range enrollCode {
		if !contains(validChars, char) {
			t.Errorf("Generated enroll code contains invalid character: %c", char)
		}
	}
}

// 检查字符是否包含在字符串中
func contains(str string, char rune) bool {
	for _, c := range str {
		if c == char {
			return true
		}
	}
	return false
}
