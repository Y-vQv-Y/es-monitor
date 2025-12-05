package util

// Min 返回最小值
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max 返回最大值
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt64 返回最小值
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// MaxInt64 返回最大值
func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// MinFloat64 返回最小值
func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// MaxFloat64 返回最大值
func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Round 四舍五入
func Round(val float64, precision int) float64 {
	ratio := 1.0
	for i := 0; i < precision; i++ {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}
