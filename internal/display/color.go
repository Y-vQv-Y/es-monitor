package display

import "github.com/fatih/color"

var (
	// 标题颜色
	TitleColor = color.New(color.FgCyan, color.Bold)

	// 章节颜色
	SectionColor = color.New(color.FgYellow, color.Bold)

	// 状态颜色
	StatusGreen  = color.New(color.FgGreen)
	StatusYellow = color.New(color.FgYellow)
	StatusRed    = color.New(color.FgRed, color.Bold)

	// 标签颜色
	LabelColor = color.New(color.FgCyan)

	// 数值颜色
	ValueColor = color.New(color.FgWhite, color.Bold)

	// 错误颜色
	ErrorColor = color.New(color.FgRed, color.Bold)

	// 成功颜色
	SuccessColor = color.New(color.FgGreen)
	
	// 警告颜色
	WarningColor = color.New(color.FgYellow)
	
	// 信息颜色
	InfoColor = color.New(color.FgBlue)
)

// GetStatusColor 根据状态获取颜色
func GetStatusColor(status string) *color.Color {
	switch status {
	case "green":
		return StatusGreen
	case "yellow":
		return StatusYellow
	case "red":
		return StatusRed
	default:
		return color.New(color.FgWhite)
	}
}

// GetThresholdColor 根据阈值获取颜色
func GetThresholdColor(value, warning, critical int) *color.Color {
	if value >= critical {
		return StatusRed
	} else if value >= warning {
		return StatusYellow
	}
	return StatusGreen
}

// GetPercentColor 根据百分比获取颜色
func GetPercentColor(percent, warning, critical float64) *color.Color {
	if percent >= critical {
		return StatusRed
	} else if percent >= warning {
		return StatusYellow
	}
	return StatusGreen
}
