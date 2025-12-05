package model

// SystemMetrics 系统指标（完整版）
type SystemMetrics struct {
	Timestamp int64
	CPU       CPUMetrics
	Memory    MemoryMetrics
	Disk      DiskMetrics
	Network   NetworkMetrics
}

// CPUMetrics CPU 详细指标
type CPUMetrics struct {
	// 总体使用率
	UsagePercent float64
	
	// 各状态百分比
	UserPercent   float64 // 用户态 CPU 使用率
	SystemPercent float64 // 系统态 CPU 使用率
	IdlePercent   float64 // 空闲 CPU 百分比
	IOWaitPercent float64 // IO 等待百分比
	IrqPercent    float64 // 硬中断百分比
	SoftIrqPercent float64 // 软中断百分比
	StealPercent  float64 // 虚拟化偷取百分比
	GuestPercent  float64 // 虚拟机 CPU 百分比
	
	// CPU 核心数
	Cores         int
	LogicalCores  int
	
	// 负载信息
	LoadAvg1  float64 // 1分钟负载
	LoadAvg5  float64 // 5分钟负载
	LoadAvg15 float64 // 15分钟负载
	
	// 每个核心的使用率
	PerCPUPercent []float64
}

// MemoryMetrics 内存详细指标
type MemoryMetrics struct {
	// 物理内存
	Total       uint64  // 总内存
	Available   uint64  // 可用内存
	Used        uint64  // 已使用内存
	UsedPercent float64 // 使用百分比
	Free        uint64  // 空闲内存
	
	// 缓存和缓冲区
	Buffers uint64 // 缓冲区大小
	Cached  uint64 // 缓存大小
	Shared  uint64 // 共享内存
	
	// Swap 内存
	SwapTotal       uint64  // Swap 总量
	SwapUsed        uint64  // Swap 已使用
	SwapFree        uint64  // Swap 空闲
	SwapUsedPercent float64 // Swap 使用百分比
	
	// 页面统计
	PageIn  uint64 // 页面换入
	PageOut uint64 // 页面换出
	
	// 内存压力
	Dirty        uint64 // 脏页
	Writeback    uint64 // 回写页
	Mapped       uint64 // 映射内存
	Slab         uint64 // Slab 内存
	
	// 内存活跃状态
	Active   uint64 // 活跃内存
	Inactive uint64 // 非活跃内存
}

// DiskMetrics 磁盘详细指标
type DiskMetrics struct {
	// 磁盘 IO 性能（实时）
	ReadBytesPerSec  float64 // 每秒读取字节数
	WriteBytesPerSec float64 // 每秒写入字节数
	ReadOpsPerSec    float64 // 每秒读操作数
	WriteOpsPerSec   float64 // 每秒写操作数
	
	// IO 延迟
	ReadLatencyMs  float64 // 读延迟（毫秒）
	WriteLatencyMs float64 // 写延迟（毫秒）
	
	// IO 利用率
	IOUtilPercent float64 // IO 使用率百分比
	
	// IO 队列
	IOQueueDepth float64 // IO 队列深度
	
	// 总计数器（累计值）
	TotalReadBytes  uint64 // 总读取字节数
	TotalWriteBytes uint64 // 总写入字节数
	TotalReadOps    uint64 // 总读操作数
	TotalWriteOps   uint64 // 总写操作数
	
	// 分区信息
	Partitions []PartitionMetrics
	
	// 每个磁盘设备的详细信息
	Devices []DiskDeviceMetrics
}

// PartitionMetrics 分区指标
type PartitionMetrics struct {
	Device      string  // 设备名
	Mountpoint  string  // 挂载点
	FSType      string  // 文件系统类型
	Total       uint64  // 总容量
	Used        uint64  // 已使用
	Free        uint64  // 空闲
	UsedPercent float64 // 使用百分比
	InodesTotal uint64  // Inode 总数
	InodesUsed  uint64  // 已使用 Inode
	InodesFree  uint64  // 空闲 Inode
}

// DiskDeviceMetrics 磁盘设备指标
type DiskDeviceMetrics struct {
	Device           string  // 设备名（如 sda, nvme0n1）
	ReadBytesPerSec  float64 // 读取速率
	WriteBytesPerSec float64 // 写入速率
	ReadOpsPerSec    float64 // 读操作速率
	WriteOpsPerSec   float64 // 写操作速率
	IOUtilPercent    float64 // IO 使用率
	AvgQueueSize     float64 // 平均队列长度
	AvgRequestSize   float64 // 平均请求大小
}

// NetworkMetrics 网络详细指标
type NetworkMetrics struct {
	// 总体吞吐量（实时）
	BytesSentPerSec float64 // 每秒发送字节数
	BytesRecvPerSec float64 // 每秒接收字节数
	PacketsSentPerSec float64 // 每秒发送包数
	PacketsRecvPerSec float64 // 每秒接收包数
	
	// 错误和丢包率
	ErrorsPerSec float64 // 每秒错误数
	DropsPerSec  float64 // 每秒丢包数
	
	// 总计数器（累计值）
	TotalBytesSent   uint64 // 总发送字节数
	TotalBytesRecv   uint64 // 总接收字节数
	TotalPacketsSent uint64 // 总发送包数
	TotalPacketsRecv uint64 // 总接收包数
	TotalErrors      uint64 // 总错误数
	TotalDrops       uint64 // 总丢包数
	
	// 连接统计
	TCPConnections    int // TCP 连接数
	TCPEstablished    int // 已建立的 TCP 连接
	TCPListening      int // 监听状态的 TCP 连接
	TCPTimeWait       int // TIME_WAIT 状态的连接
	UDPConnections    int // UDP 连接数
	
	// 每个网卡的详细信息
	Interfaces []InterfaceMetrics
}

// InterfaceMetrics 网卡详细指标
type InterfaceMetrics struct {
	Name string // 网卡名称（如 eth0, ens33）
	
	// 实时速率
	BytesSentPerSec   float64 // 发送速率
	BytesRecvPerSec   float64 // 接收速率
	PacketsSentPerSec float64 // 发送包速率
	PacketsRecvPerSec float64 // 接收包速率
	
	// 累计统计
	BytesSent   uint64 // 总发送字节数
	BytesRecv   uint64 // 总接收字节数
	PacketsSent uint64 // 总发送包数
	PacketsRecv uint64 // 总接收包数
	
	// 错误统计
	ErrorsIn  uint64 // 接收错误
	ErrorsOut uint64 // 发送错误
	DropsIn   uint64 // 接收丢包
	DropsOut  uint64 // 发送丢包
	
	// 网卡状态
	IsUp    bool   // 是否启用
	MTU     int    // MTU 大小
	Speed   uint64 // 网卡速度（Mbps）
	Duplex  string // 双工模式
	
	// IP 地址
	IPv4Addresses []string
	IPv6Addresses []string
}
