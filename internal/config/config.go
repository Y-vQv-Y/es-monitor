package config

import "time"

// Config 配置结构
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Interval time.Duration
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
