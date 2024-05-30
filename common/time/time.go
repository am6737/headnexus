package time

import (
	"fmt"
	"time"
)

// FormatTimestamp 格式化时间戳为字符串，如果为0则返回空字符串
func FormatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp/1000, 0).Format("2006-01-02 15:04:05")
}

// CurrentTimestampMillis 返回当前时间的毫秒级时间戳
func CurrentTimestampMillis() int64 {
	return time.Now().UnixNano() / 1e6
}

func FormatTimeSince(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}

	now := time.Now()
	t := time.Unix(timestamp/1000, 0)
	diff := now.Sub(t)

	// 根据时间差的大小选择合适的格式化方式
	switch {
	case diff < time.Minute:
		return "1分钟前"
	case diff < time.Hour:
		return fmt.Sprintf("%d分钟前", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d小时前", int(diff.Hours()))
	default:
		// 如果时间差大于一天，则返回具体时间
		return t.Format("2006-01-02 15:04:05")
	}
}
