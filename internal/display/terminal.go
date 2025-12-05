package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/es-monitor/internal/config"
	"github.com/yourusername/es-monitor/internal/model"
)

const (
	DisplayWidth = 120
)

// Terminal 终端显示
type Terminal struct {
	thresholds config.Thresholds
}

// NewTerminal 创建终端显示器
func NewTerminal() *Terminal {
	return &Terminal{
		thresholds: config.DefaultThresholds,
	}
}

// Clear 清屏
func (t *Terminal) Clear() {
	fmt.Print("\033[H\033[2J")
}

// DisplayHeader 显示标题
func (t *Terminal) DisplayHeader() {
	t.Clear()
	TitleColor.Println(DrawSeparator(DisplayWidth, "="))
	TitleColor.Println(PadRight("  Elasticsearch 实时监控工具 v1.0", DisplayWidth))
	TitleColor.Println(DrawSeparator(DisplayWidth, "="))
	fmt.Println()
}

// DisplayClusterHealth 显示集群健康状态
func (t *Terminal) DisplayClusterHealth(health *model.ClusterHealth) {
	SectionColor.Println("[集群健康状态]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	fmt.Printf("集群名称: %s\n", ValueColor.Sprint(health.ClusterName))
	
	statusColor := GetStatusColor(health.Status)
	fmt.Printf("集群状态: %s\n", statusColor.Sprint(strings.ToUpper(health.Status)))

	fmt.Printf("节点总数: %s  (数据节点: %s)\n", 
		ValueColor.Sprint(health.NumberOfNodes), 
		ValueColor.Sprint(health.NumberOfDataNodes))

	fmt.Printf("活跃分片: %s  (主分片: %s)\n", 
		ValueColor.Sprint(health.ActiveShards), 
		ValueColor.Sprint(health.ActivePrimaryShards))

	// 显示异常分片
	if health.RelocatingShards > 0 {
		StatusYellow.Printf("迁移中分片: %d\n", health.RelocatingShards)
	}
	if health.InitializingShards > 0 {
		StatusYellow.Printf("初始化分片: %d\n", health.InitializingShards)
	}
	if health.UnassignedShards > 0 {
		StatusRed.Printf("未分配分片: %d [严重]\n", health.UnassignedShards)
	}
	if health.PendingTasks > 0 {
		StatusYellow.Printf("待处理任务: %d\n", health.PendingTasks)
	}

	fmt.Printf("分片活跃率: %s\n", ValueColor.Sprintf("%.2f%%", health.ActiveShardsPercent))
	fmt.Println()
}

// DisplayNodeStats 显示节点统计
func (t *Terminal) DisplayNodeStats(stats *model.NodeStats, prevData map[string]*PrevNodeMetrics) {
	SectionColor.Println("[节点详细统计]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	for nodeID, node := range stats.Nodes {
		fmt.Printf("\n节点: %s (IP: %s)\n", LabelColor.Sprint(node.Name), node.IP)
		fmt.Println(DrawSeparator(DisplayWidth, "."))

		// JVM 内存
		heapPercent := node.JVM.Mem.HeapUsedPercent
		heapUsedGB := float64(node.JVM.Mem.HeapUsedInBytes) / 1024 / 1024 / 1024
		heapMaxGB := float64(node.JVM.Mem.HeapMaxInBytes) / 1024 / 1024 / 1024

		heapColor := GetThresholdColor(heapPercent, t.thresholds.JVMHeapWarning, t.thresholds.JVMHeapCritical)
		fmt.Print("  JVM 堆内存: ")
		heapColor.Printf("%d%% (%.2fG/%.2fG)", heapPercent, heapUsedGB, heapMaxGB)
		if heapPercent >= t.thresholds.JVMHeapCritical {
			fmt.Print(" [严重: 内存压力大]")
		} else if heapPercent >= t.thresholds.JVMHeapWarning {
			fmt.Print(" [警告: 内存偏高]")
		}
		fmt.Println()

		// 非堆内存
		nonHeapMB := float64(node.JVM.Mem.NonHeapUsedInBytes) / 1024 / 1024
		fmt.Printf("  JVM 非堆内存: %.2f MB\n", nonHeapMB)

		// JVM 线程
		fmt.Printf("  JVM 线程数: %d (峰值: %d)\n", 
			node.JVM.Threads.Count, 
			node.JVM.Threads.PeakCount)

		// GC 统计
		youngGC := node.JVM.GC.Collectors.Young.CollectionCount
		youngGCTime := node.JVM.GC.Collectors.Young.CollectionTimeInMillis
		oldGC := node.JVM.GC.Collectors.Old.CollectionCount
		oldGCTime := node.JVM.GC.Collectors.Old.CollectionTimeInMillis

		fmt.Printf("  GC 统计: Young=%d (耗时: %s), Old=%d (耗时: %s)", 
			youngGC, FormatDuration(int64(youngGCTime)),
			oldGC, FormatDuration(int64(oldGCTime)))
		if oldGC > 10 {
			StatusYellow.Print(" [警告: Full GC 频繁]")
		}
		fmt.Println()

		// CPU 使用率
		cpuPercent := node.OS.CPU.Percent
		cpuColor := GetThresholdColor(cpuPercent, t.thresholds.CPUWarning, t.thresholds.CPUCritical)
		fmt.Print("  CPU 使用率: ")
		cpuColor.Printf("%d%%", cpuPercent)
		if cpuPercent >= t.thresholds.CPUCritical {
			fmt.Print(" [严重: CPU 过载]")
		} else if cpuPercent >= t.thresholds.CPUWarning {
			fmt.Print(" [警告: CPU 偏高]")
		}
		fmt.Println()

		// 系统负载
		if node.OS.CPU.LoadAverage.OneMinute > 0 {
			fmt.Printf("  系统负载: 1min=%.2f, 5min=%.2f, 15min=%.2f\n",
				node.OS.CPU.LoadAverage.OneMinute,
				node.OS.CPU.LoadAverage.FiveMinutes,
				node.OS.CPU.LoadAverage.FifteenMinutes)
		}

		// 系统内存
		memPercent := node.OS.Mem.UsedPercent
		memUsedGB := float64(node.OS.Mem.UsedInBytes) / 1024 / 1024 / 1024
		memTotalGB := float64(node.OS.Mem.TotalInBytes) / 1024 / 1024 / 1024

		memColor := GetThresholdColor(memPercent, t.thresholds.MemoryWarning, t.thresholds.MemoryCritical)
		fmt.Print("  系统内存: ")
		memColor.Printf("%d%% (%.2fG/%.2fG)", memPercent, memUsedGB, memTotalGB)
		if memPercent >= t.thresholds.MemoryCritical {
			fmt.Print(" [严重: 内存不足]")
		} else if memPercent >= t.thresholds.MemoryWarning {
			fmt.Print(" [警告: 内存偏高]")
		}
		fmt.Println()

		// Swap 内存
		if node.OS.Swap.TotalInBytes > 0 {
			swapUsedMB := float64(node.OS.Swap.UsedInBytes) / 1024 / 1024
			swapTotalMB := float64(node.OS.Swap.TotalInBytes) / 1024 / 1024
			fmt.Printf("  Swap 内存: %.2f MB / %.2f MB\n", swapUsedMB, swapTotalMB)
		}

		// 磁盘使用
		diskTotal := float64(node.FS.Total.TotalInBytes) / 1024 / 1024 / 1024
		diskAvail := float64(node.FS.Total.AvailableInBytes) / 1024 / 1024 / 1024
		diskUsedPercent := (diskTotal - diskAvail) / diskTotal * 100

		diskColor := GetThresholdColor(int(diskUsedPercent), t.thresholds.DiskWarning, t.thresholds.DiskCritical)
		fmt.Print("  磁盘使用: ")
		diskColor.Printf("%.1f%% (可用: %.2fG/%.2fG)", diskUsedPercent, diskAvail, diskTotal)
		if diskUsedPercent >= float64(t.thresholds.DiskCritical) {
			fmt.Print(" [严重: 磁盘空间不足]")
		} else if diskUsedPercent >= float64(t.thresholds.DiskWarning) {
			fmt.Print(" [警告: 磁盘空间偏低]")
		}
		fmt.Println()

		// 磁盘 IO
		if node.FS.IOStats.Total.Operations > 0 {
			fmt.Printf("  磁盘 IO: 读=%s, 写=%s (操作: 读=%d, 写=%d)\n",
				FormatBytes(node.FS.IOStats.Total.ReadKB*1024),
				FormatBytes(node.FS.IOStats.Total.WriteKB*1024),
				node.FS.IOStats.Total.ReadOps,
				node.FS.IOStats.Total.WriteOps)
		}

		// 文件描述符
		fdPercent := float64(node.Process.OpenFileDescriptors) / float64(node.Process.MaxFileDescriptors) * 100
		fmt.Printf("  文件描述符: %d / %d (%.1f%%)\n",
			node.Process.OpenFileDescriptors,
			node.Process.MaxFileDescriptors,
			fdPercent)

		// 网络统计
		fmt.Printf("  网络传输: 发送=%s, 接收=%s (连接: %d)\n",
			FormatBytes(node.Transport.TxSizeInBytes),
			FormatBytes(node.Transport.RxSizeInBytes),
			node.Transport.ServerOpen)

		// HTTP 连接
		fmt.Printf("  HTTP 连接: 当前=%d, 总计=%d\n",
			node.HTTP.CurrentOpen,
			node.HTTP.TotalOpened)

		// 索引和查询速率
		fmt.Printf("  文档统计: 总数=%s, 已删除=%s\n",
			ValueColor.Sprint(node.Indices.Docs.Count),
			ValueColor.Sprint(node.Indices.Docs.Deleted))

		fmt.Printf("  存储大小: %s\n",
			ValueColor.Sprint(FormatBytes(node.Indices.Store.SizeInBytes)))

		if prev, ok := prevData[nodeID]; ok {
			elapsed := time.Since(prev.Timestamp).Seconds()
			if elapsed > 0 {
				indexRate := float64(node.Indices.Indexing.IndexTotal-prev.IndexTotal) / elapsed
				queryRate := float64(node.Indices.Search.QueryTotal-prev.QueryTotal) / elapsed

				fmt.Printf("  写入速率: %s\n", ValueColor.Sprint(FormatRate(indexRate, "docs/s")))
				fmt.Printf("  查询速率: %s\n", ValueColor.Sprint(FormatRate(queryRate, "queries/s")))
			}
		} else {
			fmt.Println("  写入速率: 计算中...")
			fmt.Println("  查询速率: 计算中...")
		}

		// 索引操作详情
		fmt.Printf("  索引操作: 总计=%d, 当前=%d, 失败=%d (耗时: %s)\n",
			node.Indices.Indexing.IndexTotal,
			node.Indices.Indexing.IndexCurrent,
			node.Indices.Indexing.IndexFailed,
			FormatDuration(node.Indices.Indexing.IndexTimeInMillis))

		// 搜索操作详情
		fmt.Printf("  搜索操作: 总计=%d, 当前=%d, 上下文=%d (耗时: %s)\n",
			node.Indices.Search.QueryTotal,
			node.Indices.Search.QueryCurrent,
			node.Indices.Search.OpenContexts,
			FormatDuration(node.Indices.Search.QueryTimeInMillis))

		// Merge 操作
		if node.Indices.Merges.Total > 0 {
			fmt.Printf("  Merge 操作: 总计=%d, 当前=%d (耗时: %s, 大小: %s)\n",
				node.Indices.Merges.Total,
				node.Indices.Merges.Current,
				FormatDuration(node.Indices.Merges.TotalTimeInMillis),
				FormatBytes(node.Indices.Merges.TotalSizeInBytes))
		}

		// Refresh 和 Flush
		fmt.Printf("  Refresh: %d次 (耗时: %s), Flush: %d次 (耗时: %s)\n",
			node.Indices.Refresh.Total,
			FormatDuration(node.Indices.Refresh.TotalTimeInMillis),
			node.Indices.Flush.Total,
			FormatDuration(node.Indices.Flush.TotalTimeInMillis))
	}

	fmt.Println()
}

// DisplayIndexStats 显示索引统计（增强）
func (t *Terminal) DisplayIndexStats(indices []model.IndexInfo, stats *model.IndexStats, prevData map[string]*PrevIndexMetrics) {
	SectionColor.Println("[索引详细统计]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	if len(indices) == 0 {
		fmt.Println("  没有索引数据")
		fmt.Println()
		return
	}

	// 显示表头（增强）
	fmt.Printf("%-30s %-8s %-8s %-12s %-12s %-15s %-15s %-15s\n",
		"索引名称", "状态", "分片", "文档数", "大小", "写入速率", "查询速率", "分段数")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	for _, idx := range indices {
		// 限制只显示前20个索引
		if len(indices) > 20 {
			break
		}

		statusColor := GetStatusColor(idx.Health)
		indexName := TruncateString(idx.Index, 30)

		fmt.Printf("%-30s ", indexName)
		statusColor.Printf("%-8s ", strings.ToUpper(idx.Health))
		fmt.Printf("%-8s ", idx.Pri+"/"+idx.Rep)
		fmt.Printf("%-12s ", idx.DocsCount)
		fmt.Printf("%-12s ", idx.StoreSize)

		// 计算速率
		if stat, ok := stats.Indices[idx.Index]; ok {
			if prev, ok := prevData[idx.Index]; ok {
				elapsed := time.Since(prev.Timestamp).Seconds()
				if elapsed > 0 {
					indexRate := float64(stat.Total.Indexing.IndexTotal-prev.IndexTotal) / elapsed
					queryRate := float64(stat.Total.Search.QueryTotal-prev.QueryTotal) / elapsed

					fmt.Printf("%-15s ", FormatRate(indexRate, "docs/s"))
					fmt.Printf("%-15s ", FormatRate(queryRate, "q/s"))
				} else {
					fmt.Printf("%-15s %-15s ", "-", "-")
				}
			} else {
				fmt.Printf("%-15s %-15s ", "计算中", "计算中")
			}

			// 显示分段信息（新增）
			fmt.Printf("%-15s", stat.Segments.Count)
			if stat.Segments.MemoryInBytes > 0 {
				fmt.Printf(" (内存: %s)", FormatBytes(stat.Segments.MemoryInBytes))
			}
		} else {
			fmt.Printf("%-15s %-15s %-15s", "-", "-", "-")
		}

		fmt.Println()
	}

	if len(indices) > 20 {
		fmt.Printf("\n  ... 还有 %d 个索引未显示\n", len(indices)-20)
	}

	fmt.Println()
}

// DisplaySystemMetrics 显示系统指标（增强）
func (t *Terminal) DisplaySystemMetrics(metrics *model.SystemMetrics) {
	SectionColor.Println("[系统资源详细监控]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	// CPU 详细信息
	fmt.Println("\nCPU:")
	fmt.Printf("  核心数: %d\n", metrics.CPU.Cores)
	fmt.Printf("  总使用率: %.2f%%\n", metrics.CPU.UsagePercent)
	fmt.Printf("  用户态: %.2f%%, 内核态: %.2f%%, 空闲: %.2f%%, IO等待: %.2f%%\n",
		metrics.CPU.UserPercent,
		metrics.CPU.SystemPercent,
		metrics.CPU.IdlePercent,
		metrics.CPU.IOWaitPercent)

	// 内存详细信息
	fmt.Println("\n内存:")
	fmt.Printf("  总容量: %s\n", FormatBytes(int64(metrics.Memory.Total)))
	fmt.Printf("  已使用: %s (%.2f%%)\n", 
		FormatBytes(int64(metrics.Memory.Used)), 
		metrics.Memory.UsedPercent)
	fmt.Printf("  可用: %s\n", FormatBytes(int64(metrics.Memory.Available)))
	fmt.Printf("  缓冲区: %s, 缓存: %s\n",
		FormatBytes(int64(metrics.Memory.Buffers)),
		FormatBytes(int64(metrics.Memory.Cached)))
	if metrics.Memory.SwapTotal > 0 {
		swapPercent := float64(metrics.Memory.SwapUsed) / float64(metrics.Memory.SwapTotal) * 100
		fmt.Printf("  Swap: %s / %s (%.2f%%)\n",
			FormatBytes(int64(metrics.Memory.SwapUsed)),
			FormatBytes(int64(metrics.Memory.SwapTotal)),
			swapPercent)
	}

	// 磁盘详细信息（增强）
	fmt.Println("\n磁盘:")
	fmt.Printf("  IO 吞吐: 读=%s, 写=%s\n",
		FormatBytesPerSec(metrics.Disk.ReadBytesPerSec),
		FormatBytesPerSec(metrics.Disk.WriteBytesPerSec))
	fmt.Printf("  IO 操作: 读=%.1f ops/s, 写=%.1f ops/s\n",
		metrics.Disk.ReadOpsPerSec,
		metrics.Disk.WriteOpsPerSec)
	ioColor := GetThresholdColorFloat(metrics.Disk.IOUtilPercent, t.thresholds.CPUWarning, t.thresholds.CPUCritical)
	fmt.Printf("  IO 利用率: ")
	ioColor.Printf("%.2f%%", metrics.Disk.IOUtilPercent)
	fmt.Println()

	if len(metrics.Disk.Partitions) > 0 {
		fmt.Println("  分区使用:")
		for _, part := range metrics.Disk.Partitions {
			fmt.Printf("    %s (%s): %s / %s (%.2f%%)\n",
				part.Mountpoint,
				part.Device,
				FormatBytes(int64(part.Used)),
				FormatBytes(int64(part.Total)),
				part.UsedPercent)
			if part.IOStats.ReadBytesPerSec > 0 || part.IOStats.WriteBytesPerSec > 0 {
				fmt.Printf("      IO: 读=%s (%d ops/s), 写=%s (%d ops/s), 利用率=%.2f%%\n",
					FormatBytesPerSec(part.IOStats.ReadBytesPerSec), int(part.IOStats.ReadOpsPerSec),
					FormatBytesPerSec(part.IOStats.WriteBytesPerSec), int(part.IOStats.WriteOpsPerSec),
					part.IOStats.IOUtilPercent)
			}
		}
	}

	// 网络详细信息（增强）
	fmt.Println("\n网络:")
	fmt.Printf("  吞吐量: 发送=%s, 接收=%s\n",
		FormatBytesPerSec(metrics.Network.BytesSentPerSec),
		FormatBytesPerSec(metrics.Network.BytesRecvPerSec))
	fmt.Printf("  数据包: 发送=%.1f pkt/s, 接收=%.1f pkt/s\n",
		metrics.Network.PacketsSentPerSec,
		metrics.Network.PacketsRecvPerSec)

	if len(metrics.Network.Interfaces) > 0 {
		fmt.Println("  网卡统计:")
		for _, iface := range metrics.Network.Interfaces {
			if iface.BytesSent == 0 && iface.BytesRecv == 0 {
				continue
			}
			fmt.Printf("    %s: 发送=%s (%s), 接收=%s (%s)\n",
				iface.Name,
				FormatBytes(int64(iface.BytesSent)),
				FormatBytesPerSec(iface.BytesSentPerSec),
				FormatBytes(int64(iface.BytesRecv)),
				FormatBytesPerSec(iface.BytesRecvPerSec))
		}
	}

	fmt.Println()
}

// DisplayFooter 显示页脚
func (t *Terminal) DisplayFooter() {
	fmt.Println(DrawSeparator(DisplayWidth, "="))
	fmt.Printf("最后更新: %s | 按 Ctrl+C 退出\n", time.Now().Format("2006-01-02 15:04:05"))
}

// DisplayError 显示错误
func (t *Terminal) DisplayError(msg string, err error) {
	ErrorColor.Printf("[错误] %s: %v\n", msg, err)
}

// PrevNodeMetrics 节点历史指标
type PrevNodeMetrics struct {
	IndexTotal int
	QueryTotal int
	Timestamp  time.Time
}

// PrevIndexMetrics 索引历史指标
type PrevIndexMetrics struct {
	IndexTotal int
	QueryTotal int
	Timestamp  time.Time
}
