package time

import (
	"testing"
	"time"
)

func TestFormatTimestamp(t *testing.T) {
	// 测试时间戳为0的情况
	result := FormatTimestamp(0)
	if result != "" {
		t.Errorf("Expected empty string, but got %s", result)
	}

	// 获取当前本地时间
	now := time.Now()
	// 构造一个时间对象
	expectedTime := time.Date(2024, time.May, 3, 10, 30, 0, 0, now.Location())
	// 获取时间戳
	expectedTimestamp := expectedTime.Unix() * 1000

	// 测试非0时间戳的情况
	expected := "2024-05-03 10:30:00"
	result = FormatTimestamp(expectedTimestamp)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestCurrentTimestampMillis(t *testing.T) {
	// 获取当前本地时间
	now := time.Now()

	// 获取函数返回的时间戳
	result := CurrentTimestampMillis()

	// 构造一个时间对象
	expectedTime := time.Unix(result/1000, 0)
	// 检查函数返回的时间戳与当前本地时间的差距是否在1秒内
	if diff := expectedTime.Sub(now); diff < -time.Second || diff > time.Second {
		t.Errorf("Expected timestamp within 1 second of current time, but got %d", result)
	}
}

func TestFormatTimeSince(t *testing.T) {
	// 获取当前本地时间
	now := time.Now()

	// 测试1分钟前的情况
	oneMinuteAgo := now.Add(-1 * time.Minute).Unix()
	result := FormatTimeSince(oneMinuteAgo)
	if result != "1分钟前" {
		t.Errorf("Expected 1分钟前, but got %s", result)
	}

	// 测试1小时前的情况
	oneHourAgo := now.Add(-1 * time.Hour).Unix()
	result = FormatTimeSince(oneHourAgo)
	if result != "1小时前" {
		t.Errorf("Expected 1小时前, but got %s", result)
	}

	// 测试1天前的情况
	oneDayAgo := now.Add(-24 * time.Hour).Unix()
	expected := now.Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	result = FormatTimeSince(oneDayAgo)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// 测试超过1天前的情况
	olderTime := time.Date(2024, time.May, 1, 10, 0, 0, 0, now.Location()).Unix()
	expected = "2024-05-01 10:00:00"
	result = FormatTimeSince(olderTime)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
