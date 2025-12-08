package display

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// FormatBytes 格式化字节数（精确版本）
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatBytesUint64 格式化 uint64 字节数
func FormatBytesUint64(bytes uint64) string {
	return FormatBytes(int64(bytes))
}

// FormatBytesPerSec 格式化每秒字节数
func FormatBytesPerSec(bytesPerSec float64) string {
	return fmt.Sprintf("%s/s", FormatBytes(int64(bytesPerSec)))
}

// FormatBandwidth 格式化带宽（支持多种单位显示）
func FormatBandwidth(mbps float64) string {
	if mbps < 1 {
		return fmt.Sprintf("%.2f Kbps", mbps*1024)
	} else if mbps < 1024 {
		return fmt.Sprintf("%.2f Mbps", mbps)
	} else {
		return fmt.Sprintf("%.2f Gbps", mbps/1024)
	}
}

// ParseESSize 解析 ES 返回的大小字符串（如 "1.5mb", "256kb", "2048"）
func ParseESSize(sizeStr string) string {
	// 如果已经是格式化的字符串（包含单位），直接返回
	if matched, _ := regexp.MatchString(`(?i)\d+(\.\d+)?\s*(b|kb|mb|gb|tb)`, sizeStr); matched {
		return normalizeSize(sizeStr)
	}
	
	// 如果是纯数字，当作字节处理
	if bytes, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
		return FormatBytes(bytes)
	}
	
	return sizeStr
}

// normalizeSize 规范化大小格式
func normalizeSize(sizeStr string) string {
	sizeStr = strings.ToLower(strings.TrimSpace(sizeStr))
	
	// 提取数字和单位
	re := regexp.MustCompile(`([\d.]+)\s*([a-z]+)`)
	matches := re.FindStringSubmatch(sizeStr)
	if len(matches) != 3 {
		return sizeStr
	}
	
	value, _ := strconv.ParseFloat(matches[1], 64)
	unit := matches[2]
	
	// 转换为字节
	var bytes int64
	switch unit {
	case "b":
		bytes = int64(value)
	case "kb":
		bytes = int64(value * 1024)
	case "mb":
		bytes = int64(value * 1024 * 1024)
	case "gb":
		bytes = int64(value * 1024 * 1024 * 1024)
	case "tb":
		bytes = int64(value * 1024 * 1024 * 1024 * 1024)
	default:
		return sizeStr
	}
	
	return FormatBytes(bytes)
}

// FormatPercent 格式化百分比
func FormatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

// FormatRate 格式化速率
func FormatRate(rate float64, unit string) string {
	if rate < 0.01 {
		return fmt.Sprintf("0.0 %s", unit)
	} else if rate < 1000 {
		return fmt.Sprintf("%.1f %s", rate, unit)
	} else if rate < 1000000 {
		return fmt.Sprintf("%.1f K%s", rate/1000, unit)
	} else {
		return fmt.Sprintf("%.1f M%s", rate/1000000, unit)
	}
}

// FormatDuration 格式化时长
func FormatDuration(millis int64) string {
	seconds := millis / 1000
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
	} else {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
}

// DrawSeparator 绘制分隔线
func DrawSeparator(width int, char string) string {
	return strings.Repeat(char, width)
}

// PadRight 右侧填充到指定宽度
func PadRight(s string, width int) string {
	sLen := len(s)
	if sLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-sLen)
}

// PadLeft 左侧填充
func PadLeft(s string, width int) string {
	sLen := len(s)
	if sLen >= width {
		return s
	}
	return strings.Repeat(" ", width-sLen) + s
}

// TruncateString 截断字符串
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatFloat 格式化浮点数
func FormatFloat(value float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, value)
}

// AlignRight 右对齐字符串
func AlignRight(s string, width int) string {
	return fmt.Sprintf("%*s", width, s)
}

// AlignLeft 左对齐字符串
func AlignLeft(s string, width int) string {
	return fmt.Sprintf("%-*s", width, s)
}
