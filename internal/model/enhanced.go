package model

// EnhancedMetrics ES 增强监控指标
type EnhancedMetrics struct {
	// 集群级别
	ClusterStats ClusterStats
	
	// 节点级别
	NodeThreadPools  map[string]ThreadPoolStats
	NodeCircuitBreakers map[string]CircuitBreakerStats
	
	// 索引级别
	SlowLogs       []SlowLog
	RejectedTasks  RejectedTaskStats
	
	// 健康检查
	HealthIssues   []HealthIssue
}

// ClusterStats 集群统计（额外指标）
type ClusterStats struct {
	// 查询性能
	SearchQueueSize      int     // 搜索队列大小
	SearchRejected       int64   // 搜索拒绝数
	IndexQueueSize       int     // 索引队列大小
	IndexRejected        int64   // 索引拒绝数
	
	// 集群吞吐
	TotalIndexingRate    float64 // 集群总索引速率
	TotalSearchRate      float64 // 集群总搜索速率
	
	// 分片状态
	UnassignedShardAge   int64   // 未分配分片持续时间（秒）
	RelocatingShardCount int     // 正在迁移的分片数
	
	// 段统计
	TotalSegments        int     // 总段数
	SegmentMemoryMB      float64 // 段内存占用
	
	// 缓存统计
	FieldDataMemoryMB    float64 // FieldData 内存
	QueryCacheMemoryMB   float64 // 查询缓存内存
	RequestCacheMemoryMB float64 // 请求缓存内存
}

// ThreadPoolStats 线程池统计
type ThreadPoolStats struct {
	PoolName   string
	Active     int   // 活跃线程
	Queue      int   // 队列中任务
	QueueSize  int   // 队列大小
	Rejected   int64 // 拒绝数
	Completed  int64 // 完成数
	Threads    int   // 线程数
}

// CircuitBreakerStats 断路器统计
type CircuitBreakerStats struct {
	Name          string
	LimitSizeMB   float64 // 限制大小
	EstimatedMB   float64 // 估计使用
	Overhead      float64 // 开销倍数
	Tripped       int64   // 触发次数
	UsedPercent   float64 // 使用百分比
}

// SlowLog 慢查询日志
type SlowLog struct {
	Index     string
	Type      string  // search 或 index
	TookMs    int64   // 耗时（毫秒）
	Timestamp int64
	Source    string
}

// RejectedTaskStats 拒绝任务统计
type RejectedTaskStats struct {
	SearchRejected  int64
	IndexRejected   int64
	BulkRejected    int64
	GetRejected     int64
}

// HealthIssue 健康问题
type HealthIssue struct {
	Level       string // critical, warning, info
	Component   string // cluster, node, index, jvm, disk, etc.
	NodeName    string
	IndexName   string
	Message     string
	Value       interface{}
	Threshold   interface{}
	Timestamp   int64
	Suggestion  string // 修复建议
}
