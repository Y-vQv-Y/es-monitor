package display

import (
	"fmt"
	"strings"

	"github.com/yourusername/es-monitor/pkg/util"
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
	return util.FormatDuration(time.Duration(millis * int64(time.Millisecond))).String()
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
