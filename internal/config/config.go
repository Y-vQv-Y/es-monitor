package config

import "time"

// Config 配置结构
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Interval time.Duration
	ReadOnly bool // 只读模式，生产环境必须为 true
}

// Thresholds 阈值配置
type Thresholds struct {
	JVMHeapWarning  int // JVM 堆内存警告阈值
	JVMHeapCritical int // JVM 堆内存严重阈值
	CPUWarning      int // CPU 警告阈值
	CPUCritical     int // CPU 严重阈值
	MemoryWarning   int // 内存警告阈值
	MemoryCritical  int // 内存严重阈值
	DiskWarning     int // 磁盘警告阈值
	DiskCritical    int // 磁盘严重阈值
}

// DefaultThresholds 默认阈值
var DefaultThresholds = Thresholds{
	JVMHeapWarning:  75,
	JVMHeapCritical: 85,
	CPUWarning:      60,
	CPUCritical:     80,
	MemoryWarning:   80,
	MemoryCritical:  90,
	DiskWarning:     80,
	DiskCritical:    90,
}

// SafetyConfig 生产环境安全配置
type SafetyConfig struct {
	// 只允许读取操作的 API 端点
	AllowedEndpoints []string
	// 请求超时时间（避免长时间占用连接）
	RequestTimeout time.Duration
	// 最大并发请求数（避免对 ES 造成压力）
	MaxConcurrency int
}

// DefaultSafetyConfig 默认安全配置
var DefaultSafetyConfig = SafetyConfig{
	AllowedEndpoints: []string{
		"/_cluster/health",
		"/_nodes/stats",
		"/_stats",
		"/_cat/indices",
		"/",
	},
	RequestTimeout: 10 * time.Second,
	MaxConcurrency: 5,
}
