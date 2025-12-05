package display

import (
	"fmt"
	"strings"
)

// FormatBytes 格式化字节数
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

// FormatPercent 格式化百分比
func FormatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

// FormatRate 格式化速率
func FormatRate(rate float64, unit string) string {
	if rate < 1000 {
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

// PadRight 右侧填充
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// PadLeft 左侧填充
func PadLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
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
