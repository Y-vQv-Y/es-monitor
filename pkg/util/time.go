package util

import "time"

// FormatDuration 格式化时长
func FormatDuration(d time.Duration) time.Duration {
	if d < time.Minute {
		return d.Round(time.Second)
	} else if d < time.Hour {
		return d.Round(time.Minute)
	} else {
		return d.Round(time.Hour)
	}
}

// ParseDuration 解析时长字符串
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// Now 返回当前时间戳（秒）
func Now() int64 {
	return time.Now().Unix()
}

// NowMillis 返回当前时间戳（毫秒）
func NowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
