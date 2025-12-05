package model

// SystemMetrics 系统指标
type SystemMetrics struct {
	Timestamp   int64
	CPU         CPUMetrics
	Memory      MemoryMetrics
	Disk        DiskMetrics
	Network     NetworkMetrics
}

// CPUMetrics CPU 指标
type CPUMetrics struct {
	UsagePercent  float64
	UserPercent   float64
	SystemPercent float64
	IdlePercent   float64
	IOWaitPercent float64
	Cores         int
}

// MemoryMetrics 内存指标
type MemoryMetrics struct {
	Total       uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
	Free        uint64
	Buffers     uint64
	Cached      uint64
	SwapTotal   uint64
	SwapUsed    uint64
	SwapFree    uint64
}

// DiskMetrics 磁盘指标
type DiskMetrics struct {
	ReadBytesPerSec  float64
	WriteBytesPerSec float64
	ReadOpsPerSec    float64
	WriteOpsPerSec   float64
	IOUtilPercent    float64 // 新增：IO 利用率
	Partitions       []PartitionMetrics
}

// PartitionMetrics 分区指标
type PartitionMetrics struct {
	Device      string
	Mountpoint  string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
	IOStats     PartitionIOStats // 新增：分区 IO 详细
}

// PartitionIOStats 分区 IO 统计（新增）
type PartitionIOStats struct {
	ReadBytesPerSec  float64
	WriteBytesPerSec float64
	ReadOpsPerSec    float64
	WriteOpsPerSec   float64
	IOUtilPercent    float64
}

// NetworkMetrics 网络指标
type NetworkMetrics struct {
	BytesSentPerSec float64
	BytesRecvPerSec float64
	PacketsSentPerSec float64
	PacketsRecvPerSec float64
	Interfaces      []InterfaceMetrics
}

// InterfaceMetrics 网卡指标
type InterfaceMetrics struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	BytesSentPerSec float64 // 新增：瞬时速率
	BytesRecvPerSec float64
	PacketsSentPerSec float64
	PacketsRecvPerSec float64
}
