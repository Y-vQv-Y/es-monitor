package display

import (
	"fmt"
	"strings"
	"time"
	"strconv"
	"os"
	"os/exec"
	"runtime"

	"github.com/Y-vQv-Y/es-monitor/internal/config"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

const (
	DisplayWidth = 120
)

// Terminal 终端显示
type Terminal struct {
	thresholds config.Thresholds
	lastLines  int // 记录上次输出的行数
}

// NewTerminal 创建终端显示器
func NewTerminal() *Terminal {
	return &Terminal{
		thresholds: config.DefaultThresholds,
		lastLines:  0,
	}
}

// Clear 清屏 - 完全清空屏幕并移动光标到左上角
func (t *Terminal) Clear() {
	fmt.Print("\033[2J\033[H")
}

// ClearAndReset 清空并重置光标到开始位置（类似top的刷新方式）
func (t *Terminal) ClearAndReset() {
	// 移动光标到屏幕左上角
	fmt.Print("\033[H")
	// 清除从光标到屏幕末尾的所有内容
	fmt.Print("\033[J")
}

// GetTerminalSize 获取终端尺寸
func (t *Terminal) GetTerminalSize() (width, height int) {
	width, height = 120, 40 // 默认值
	
	if runtime.GOOS == "windows" {
		// Windows 使用 mode con 命令
		cmd := exec.Command("cmd", "/c", "mode", "con")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Lines:") || strings.Contains(line, "行:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						if h, err := strconv.Atoi(fields[len(fields)-1]); err == nil {
							height = h
						}
					}
				}
				if strings.Contains(line, "Columns:") || strings.Contains(line, "列:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						if w, err := strconv.Atoi(fields[len(fields)-1]); err == nil {
							width = w
						}
					}
				}
			}
		}
	} else {
		// Unix-like 系统使用 stty size
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		output, err := cmd.Output()
		if err == nil {
			fields := strings.Fields(string(output))
			if len(fields) == 2 {
				if h, err := strconv.Atoi(fields[0]); err == nil {
					height = h
				}
				if w, err := strconv.Atoi(fields[1]); err == nil {
					width = w
				}
			}
		}
	}
	
	return width, height
}

// DisplayHeader 显示标题
func (t *Terminal) DisplayHeader() {
	t.ClearAndReset() // 使用新的清屏方式
	TitleColor.Println(DrawSeparator(DisplayWidth, "="))
	TitleColor.Println(PadRight("  Elasticsearch 生产环境监控工具 (只读安全模式)", DisplayWidth))
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

// DisplaySystemMetrics 显示系统资源详细信息（完整版）
func (t *Terminal) DisplaySystemMetrics(metrics *model.SystemMetrics) {
	SectionColor.Println("[系统资源详细监控 - 实时数据]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	// ===== CPU 详细信息 =====
	fmt.Println()
	LabelColor.Println("【CPU 详细信息】")
	fmt.Println(DrawSeparator(90, "."))

	fmt.Printf("  核心数: 物理=%d, 逻辑=%d\n",
		metrics.CPU.Cores,
		metrics.CPU.LogicalCores)

	// CPU 总体使用率
	cpuColor := GetPercentColor(metrics.CPU.UsagePercent, 60, 80)
	fmt.Print("  总使用率: ")
	cpuColor.Printf("%.2f%%", metrics.CPU.UsagePercent)
	if metrics.CPU.UsagePercent >= 80 {
		fmt.Print(" [严重: CPU 过载]")
	} else if metrics.CPU.UsagePercent >= 60 {
		fmt.Print(" [警告: CPU 偏高]")
	}
	fmt.Println()

	// CPU 详细时间分布
	fmt.Printf("  时间分布: 用户态=%.2f%%, 内核态=%.2f%%, 空闲=%.2f%%\n",
		metrics.CPU.UserPercent,
		metrics.CPU.SystemPercent,
		metrics.CPU.IdlePercent)

	// IO 等待
	if metrics.CPU.IOWaitPercent > 0 {
		ioWaitColor := GetPercentColor(metrics.CPU.IOWaitPercent, 20, 40)
		fmt.Print("  IO 等待: ")
		ioWaitColor.Printf("%.2f%%", metrics.CPU.IOWaitPercent)
		if metrics.CPU.IOWaitPercent >= 40 {
			fmt.Print(" [严重: IO 压力大]")
		} else if metrics.CPU.IOWaitPercent >= 20 {
			fmt.Print(" [警告: 存在 IO 瓶颈]")
		}
		fmt.Println()
	}

	// 中断信息
	if metrics.CPU.IrqPercent > 0 || metrics.CPU.SoftIrqPercent > 0 {
		fmt.Printf("  中断: 硬中断=%.2f%%, 软中断=%.2f%%\n",
			metrics.CPU.IrqPercent,
			metrics.CPU.SoftIrqPercent)
	}

	// 系统负载
	fmt.Printf("  系统负载: 1分钟=%.2f, 5分钟=%.2f, 15分钟=%.2f\n",
		metrics.CPU.LoadAvg1,
		metrics.CPU.LoadAvg5,
		metrics.CPU.LoadAvg15)

	// 负载评估
	avgLoad := metrics.CPU.LoadAvg1
	coreCount := float64(metrics.CPU.LogicalCores)
	if coreCount > 0 {
		loadPercent := (avgLoad / coreCount) * 100
		loadColor := GetPercentColor(loadPercent, 70, 90)
		fmt.Print("  负载压力: ")
		loadColor.Printf("%.2f%%", loadPercent)
		if loadPercent >= 90 {
			fmt.Print(" [严重: 系统压力大]")
		} else if loadPercent >= 70 {
			fmt.Print(" [警告: 负载较高]")
		}
		fmt.Println()
	}

	// 每个核心的使用率（如果核心数不多，显示详细信息）
	if len(metrics.CPU.PerCPUPercent) > 0 && len(metrics.CPU.PerCPUPercent) <= 16 {
		fmt.Print("  各核心使用率: ")
		for i, percent := range metrics.CPU.PerCPUPercent {
			if i > 0 && i%8 == 0 {
				fmt.Print("\n                  ")
			}
			cpuColor := GetPercentColor(percent, 70, 85)
			cpuColor.Printf("CPU%d=%.1f%% ", i, percent)
		}
		fmt.Println()
	}

	// ===== 内存详细信息 =====
	fmt.Println()
	LabelColor.Println("【内存详细信息】")
	fmt.Println(DrawSeparator(90, "."))

	// 物理内存
	memColor := GetPercentColor(metrics.Memory.UsedPercent, 80, 90)
	fmt.Printf("  物理内存: 总量=%s\n", FormatBytesUint64(metrics.Memory.Total))
	fmt.Print("            已使用=")
	memColor.Printf("%s (%.2f%%)", FormatBytesUint64(metrics.Memory.Used), metrics.Memory.UsedPercent)
	if metrics.Memory.UsedPercent >= 90 {
		fmt.Print(" [严重: 内存不足]")
	} else if metrics.Memory.UsedPercent >= 80 {
		fmt.Print(" [警告: 内存偏高]")
	}
	fmt.Println()
	fmt.Printf("            可用=%s, 空闲=%s\n",
		FormatBytesUint64(metrics.Memory.Available),
		FormatBytesUint64(metrics.Memory.Free))

	// 缓存和缓冲区
	fmt.Printf("  缓存: 缓冲区=%s, 页面缓存=%s, 共享内存=%s\n",
		FormatBytesUint64(metrics.Memory.Buffers),
		FormatBytesUint64(metrics.Memory.Cached),
		FormatBytesUint64(metrics.Memory.Shared))

	// 内存状态
	fmt.Printf("  内存状态: 活跃=%s, 非活跃=%s\n",
		FormatBytesUint64(metrics.Memory.Active),
		FormatBytesUint64(metrics.Memory.Inactive))

	// 脏页（Writeback 已移除，因为某些平台不支持）
	if metrics.Memory.Dirty > 0 {
		dirtyMB := float64(metrics.Memory.Dirty) / 1024 / 1024
		fmt.Printf("  脏页: %.2f MB\n", dirtyMB)
	}

	// Slab 内存
	if metrics.Memory.Slab > 0 {
		slabMB := float64(metrics.Memory.Slab) / 1024 / 1024
		fmt.Printf("  内核 Slab: %.2f MB\n", slabMB)
	}

	// Swap 内存
	if metrics.Memory.SwapTotal > 0 {
		swapColor := GetPercentColor(metrics.Memory.SwapUsedPercent, 50, 80)
		fmt.Print("  Swap 内存: ")
		swapColor.Printf("%s / %s (%.2f%%)",
			FormatBytesUint64(metrics.Memory.SwapUsed),
			FormatBytesUint64(metrics.Memory.SwapTotal),
			metrics.Memory.SwapUsedPercent)
		if metrics.Memory.SwapUsedPercent >= 80 {
			fmt.Print(" [严重: Swap 使用过高，可能影响性能]")
		} else if metrics.Memory.SwapUsedPercent >= 50 {
			fmt.Print(" [警告: Swap 使用较多]")
		}
		fmt.Println()

		// 页面换入换出
		if metrics.Memory.PageIn > 0 || metrics.Memory.PageOut > 0 {
			fmt.Printf("              页面换入=%s, 页面换出=%s\n",
				FormatBytesUint64(metrics.Memory.PageIn),
				FormatBytesUint64(metrics.Memory.PageOut))
		}
	}

	// 内存压力评估
	availablePercent := float64(metrics.Memory.Available) / float64(metrics.Memory.Total) * 100
	if availablePercent < 10 {
		StatusRed.Println("  [内存压力评估: 严重 - 可用内存极少，系统可能出现 OOM]")
	} else if availablePercent < 20 {
		StatusYellow.Println("  [内存压力评估: 警告 - 可用内存较少，建议释放内存]")
	} else {
		StatusGreen.Printf("  [内存压力评估: 正常 - 可用内存充足 (%.2f%%)]\n", availablePercent)
	}

	fmt.Println()
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


// DisplayDiskMetrics 显示磁盘详细信息（完整版）
func (t *Terminal) DisplayDiskMetrics(metrics *model.DiskMetrics) {
	LabelColor.Println("【磁盘详细信息 - 实时 IO 性能】")
	fmt.Println(DrawSeparator(90, "."))

	// 磁盘 IO 吞吐量
	fmt.Printf("  IO 吞吐量: 读取=%s, 写入=%s\n",
		FormatBytesPerSec(metrics.ReadBytesPerSec),
		FormatBytesPerSec(metrics.WriteBytesPerSec))

	// IO 操作速率
	fmt.Printf("  IO 操作数: 读操作=%.1f ops/s, 写操作=%.1f ops/s\n",
		metrics.ReadOpsPerSec,
		metrics.WriteOpsPerSec)

	// 总操作数
	totalOps := metrics.ReadOpsPerSec + metrics.WriteOpsPerSec
	if totalOps > 0 {
		fmt.Printf("              总操作=%.1f ops/s", totalOps)
		
		// 评估 IO 压力
		if totalOps > 5000 {
			StatusRed.Print(" [严重: IO 压力大]")
		} else if totalOps > 2000 {
			StatusYellow.Print(" [警告: IO 较繁忙]")
		}
		fmt.Println()
	}

	// IO 利用率
	if metrics.IOUtilPercent > 0 {
		ioColor := GetPercentColor(metrics.IOUtilPercent, 70, 90)
		fmt.Print("  IO 利用率: ")
		ioColor.Printf("%.2f%%", metrics.IOUtilPercent)
		if metrics.IOUtilPercent >= 90 {
			fmt.Print(" [严重: 磁盘接近饱和]")
		} else if metrics.IOUtilPercent >= 70 {
			fmt.Print(" [警告: 磁盘压力较大]")
		}
		fmt.Println()
	}

	// 累计 IO 统计
	if metrics.TotalReadBytes > 0 || metrics.TotalWriteBytes > 0 {
		fmt.Printf("  累计 IO: 读取=%s, 写入=%s\n",
			FormatBytesUint64(metrics.TotalReadBytes),
			FormatBytesUint64(metrics.TotalWriteBytes))
	}

	// 各磁盘设备详细信息
	if len(metrics.Devices) > 0 {
		fmt.Println()
		fmt.Println("  各设备 IO 详情:")
		fmt.Printf("  %-12s %15s %15s %12s %12s %10s\n",
			"设备", "读速率", "写速率", "读 ops/s", "写 ops/s", "IO 使用率")
		fmt.Println("  " + DrawSeparator(88, "-"))

		for _, dev := range metrics.Devices {
			// 跳过没有 IO 的设备
			if dev.ReadBytesPerSec == 0 && dev.WriteBytesPerSec == 0 {
				continue
			}

			ioColor := GetPercentColor(dev.IOUtilPercent, 70, 90)
			fmt.Printf("  %-12s %15s %15s %12.1f %12.1f ",
				dev.Device,
				FormatBytesPerSec(dev.ReadBytesPerSec),
				FormatBytesPerSec(dev.WriteBytesPerSec),
				dev.ReadOpsPerSec,
				dev.WriteOpsPerSec)
			ioColor.Printf("%9.2f%%", dev.IOUtilPercent)
			
			if dev.IOUtilPercent >= 90 {
				StatusRed.Print(" [严重]")
			} else if dev.IOUtilPercent >= 70 {
				StatusYellow.Print(" [警告]")
			}
			fmt.Println()
		}
	}

	// 分区使用情况
	if len(metrics.Partitions) > 0 {
		fmt.Println()
		fmt.Println("  磁盘分区使用情况:")
		fmt.Printf("  %-20s %-25s %12s %12s %12s %10s\n",
			"设备", "挂载点", "总容量", "已使用", "可用", "使用率")
		fmt.Println("  " + DrawSeparator(95, "-"))

		for _, part := range metrics.Partitions {
			// 跳过特殊文件系统
			if strings.HasPrefix(part.FSType, "tmpfs") ||
				strings.HasPrefix(part.FSType, "devtmpfs") ||
				strings.HasPrefix(part.Mountpoint, "/sys") ||
				strings.HasPrefix(part.Mountpoint, "/proc") {
				continue
			}

			diskColor := GetPercentColor(part.UsedPercent, 80, 90)
			mountpoint := TruncateString(part.Mountpoint, 25)
			device := TruncateString(part.Device, 20)

			fmt.Printf("  %-20s %-25s %12s %12s %12s ",
				device,
				mountpoint,
				FormatBytesUint64(part.Total),
				FormatBytesUint64(part.Used),
				FormatBytesUint64(part.Free))

			diskColor.Printf("%9.2f%%", part.UsedPercent)

			if part.UsedPercent >= 90 {
				StatusRed.Print(" [严重]")
			} else if part.UsedPercent >= 80 {
				StatusYellow.Print(" [警告]")
			}
			fmt.Println()
		}
	}

	fmt.Println()
}

// 网卡详情：
func (t *Terminal) DisplayNetworkMetrics(metrics *model.NetworkMetrics) {
	LabelColor.Println("【网络详细信息 - 实时吞吐量】")
	fmt.Println(DrawSeparator(90, "."))

	// 网络总体吞吐量（字节）
	fmt.Printf("  总体吞吐: 发送=%s, 接收=%s\n",
		FormatBytesPerSec(metrics.BytesSentPerSec),
		FormatBytesPerSec(metrics.BytesRecvPerSec))

	// 转换为 Mbps 显示
	sendMbps := float64(metrics.BytesSentPerSec) * 8.0 / 1024.0 / 1024.0
    recvMbps := float64(metrics.BytesRecvPerSec) * 8.0 / 1024.0 / 1024.0
    totalMbps := sendMbps + recvMbps
	
	
	fmt.Printf("              (发送=%s, 接收=%s)\n", 
		FormatBandwidth(sendMbps),
		FormatBandwidth(recvMbps))

	// 数据包速率
	fmt.Printf("  数据包率: 发送=%.1f pkt/s, 接收=%.1f pkt/s\n",
		metrics.PacketsSentPerSec,
		metrics.PacketsRecvPerSec)

	// 累计流量统计
	if metrics.TotalBytesSent > 0 || metrics.TotalBytesRecv > 0 {
		fmt.Printf("  累计流量: 发送=%s, 接收=%s\n",
			FormatBytesUint64(metrics.TotalBytesSent),
			FormatBytesUint64(metrics.TotalBytesRecv))
	}

	// 各网卡详细信息（改进对齐）
	if len(metrics.Interfaces) > 0 {
		fmt.Println()
		fmt.Println("  各网卡详情:")
		fmt.Printf("  %-16s %15s %15s %14s %14s %10s %10s %10s\n",
			"网卡", "发送速率", "接收速率", "发送pkt/s", "接收pkt/s", "错误", "丢包", "状态")
		fmt.Println("  " + DrawSeparator(110, "-"))

		for _, iface := range metrics.Interfaces {
			// 跳过回环和没有流量的网卡
			if strings.HasPrefix(iface.Name, "veth") ||
               strings.HasPrefix(iface.Name, "docker") ||
               strings.HasPrefix(iface.Name, "br-") ||
               strings.HasPrefix(iface.Name, "cni") ||
               strings.HasPrefix(iface.Name, "flannel") ||
			   strings.HasPrefix(iface.Name, "tunl") ||
			   strings.HasPrefix(iface.Name, "vlan") ||
			   strings.HasPrefix(iface.Name, "vxlan") ||
			   strings.HasPrefix(iface.Name, "virb") ||
			   strings.HasPrefix(iface.Name, "virbr0-") ||
			   strings.HasPrefix(iface.Name, "calico") ||
               strings.HasPrefix(iface.Name, "kube-ipvs") ||
               iface.Name == "lo" {
                 continue
            }

			ifaceName := TruncateString(iface.Name, 16)
			
			// 格式化速率（自动选择单位）
			sendRate := FormatBytesPerSec(iface.BytesSentPerSec)
			recvRate := FormatBytesPerSec(iface.BytesRecvPerSec)
			
			fmt.Printf("  %-16s %15s %15s %14.1f %14.1f %10d %10d",
				ifaceName,
				sendRate,
				recvRate,
				iface.PacketsSentPerSec,
				iface.PacketsRecvPerSec,
				iface.ErrorsIn+iface.ErrorsOut,
				iface.DropsIn+iface.DropsOut)

			// 判断状态（改进丢包判断）
			status := "正常"
			totalPackets := iface.PacketsSent + iface.PacketsRecv
			totalDrops := iface.DropsIn + iface.DropsOut
			var dropRate float64
			if totalPackets > 0 {
				dropRate = float64(totalDrops) / float64(totalPackets) * 100
			}

			if iface.ErrorsIn+iface.ErrorsOut > 100 {
				status = "有错误"
				StatusRed.Printf(" %10s", status)
			} else if dropRate > 0.1 { // 丢包率 > 0.1%
				status = "丢包"
				StatusYellow.Printf(" %10s", status)
			} else if totalDrops > 1000 { // 累计丢包 > 1000 但丢包率低
				status = "正常"
				StatusGreen.Printf(" %10s", status)
			} else {
				StatusGreen.Printf(" %10s", status)
			}
			fmt.Println()
		}
	}

	// 网络性能评估（改进带宽显示）
	fmt.Println()
	if totalMbps > 800 {
		StatusRed.Printf("  [网络负载评估: 繁忙 - 总带宽 %s (%.2f Mbps = 发送 %.2f + 接收 %.2f)]\n", 
			FormatBandwidth(totalMbps), totalMbps, sendMbps, recvMbps)
	} else if totalMbps > 500 {
		StatusYellow.Printf("  [网络负载评估: 正常偏高 - 总带宽 %s (%.2f Mbps)]\n",
			FormatBandwidth(totalMbps), totalMbps)
	} else {
		StatusGreen.Printf("  [网络负载评估: 正常 - 总带宽 %s (%.2f Mbps)]\n",
			FormatBandwidth(totalMbps), totalMbps)
	}

	fmt.Println()
}

// DisplayNodeStats 显示节点统计（简化版，避免重复）
func (t *Terminal) DisplayNodeStats(stats *model.NodeStats, prevData map[string]*PrevNodeMetrics) {
	SectionColor.Println("[Elasticsearch 节点统计]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	for nodeID, node := range stats.Nodes {
		fmt.Printf("\n节点: %s (IP: %s)\n", LabelColor.Sprint(node.Name), node.IP)
		fmt.Println(DrawSeparator(DisplayWidth, "."))

		// JVM 堆内存
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

		// GC 统计
		youngGC := node.JVM.GC.Collectors.Young.CollectionCount
		oldGC := node.JVM.GC.Collectors.Old.CollectionCount
		fmt.Printf("  GC 次数: Young=%d, Old=%d", youngGC, oldGC)
		if oldGC > 10 {
			StatusYellow.Print(" [警告: Full GC 频繁]")
		}
		fmt.Println()

		// 文档和存储
		fmt.Printf("  文档总数: %s, 存储大小: %s\n",
			ValueColor.Sprint(formatInt64WithCommas(int64(node.Indices.Docs.Count))),
			ValueColor.Sprint(FormatBytes(node.Indices.Store.SizeInBytes)))

		// 实时速率（添加计算说明）
		if prev, ok := prevData[nodeID]; ok {
			elapsed := time.Since(prev.Timestamp).Seconds()
			if elapsed > 0 {
				// 计算速率：(当前值 - 上次值) / 时间间隔
				indexDelta := node.Indices.Indexing.IndexTotal - prev.IndexTotal
				queryDelta := node.Indices.Search.QueryTotal - prev.QueryTotal
				
				indexRate := float64(indexDelta) / elapsed
				queryRate := float64(queryDelta) / elapsed

				fmt.Printf("  写入速率: %s (增量: %d 条 / %.1f 秒)\n", 
					ValueColor.Sprint(FormatRate(indexRate, "docs/s")),
					indexDelta, elapsed)
				fmt.Printf("  查询速率: %s (增量: %d 次 / %.1f 秒)\n", 
					ValueColor.Sprint(FormatRate(queryRate, "queries/s")),
					queryDelta, elapsed)
			} else {
				fmt.Println("  写入速率: 计算中... (等待下次采集)")
				fmt.Println("  查询速率: 计算中... (等待下次采集)")
			}
		} else {
			fmt.Println("  写入速率: 初始化中... (首次采集)")
			fmt.Println("  查询速率: 初始化中... (首次采集)")
		}
	}

	fmt.Println()
}

// DisplayIndexStats 显示索引统计
func (t *Terminal) DisplayIndexStats(indices []model.IndexInfo, stats *model.IndexStats, prevData map[string]*PrevIndexMetrics) {
	SectionColor.Println("[索引统计（前20个）]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	if len(indices) == 0 {
		fmt.Println("  没有索引数据")
		fmt.Println()
		return
	}

	// 显示表头（固定宽度）
	fmt.Printf("%-35s %-10s %-12s %-15s %-15s\n",
		"索引名称", "状态", "分片(主/副)", "文档数", "大小")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	count := 0
	for _, idx := range indices {
		if count >= 20 {
			break
		}
		count++

		statusColor := GetStatusColor(idx.Health)
		indexName := TruncateString(idx.Index, 35)
		
		// 格式化分片信息
		shardInfo := fmt.Sprintf("%s/%s", idx.Pri, idx.Rep)
		
		// 格式化文档数（添加千分位）
		docCount := formatNumberWithCommas(idx.DocsCount)
		
		// 格式化大小（统一转换）
		size := ParseESSize(idx.StoreSize)

		// 左对齐索引名，右对齐数字
		fmt.Printf("%-35s ", indexName)
		statusColor.Printf("%-10s ", strings.ToUpper(idx.Health))
		fmt.Printf("%-12s %15s %15s\n",
			shardInfo,
			docCount,
			size)
	}

	if len(indices) > 20 {
		fmt.Printf("\n  ... 还有 %d 个索引未显示\n", len(indices)-20)
	}

	fmt.Println()
}

// formatNumberWithCommas 为数字添加千分位分隔符
func formatNumberWithCommas(numStr string) string {
	// 如果是空或 "0"，直接返回
	if numStr == "" || numStr == "0" {
		return "0"
	}
	
	// 解析为整数
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return numStr
	}
	
	// 格式化为千分位
	return formatInt64WithCommas(num)
}

// formatInt64WithCommas 为 int64 添加千分位
func formatInt64WithCommas(n int64) string {
	if n < 0 {
		return "-" + formatInt64WithCommas(-n)
	}
	
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	
	return formatInt64WithCommas(n/1000) + "," + fmt.Sprintf("%03d", n%1000)
}

// DisplayFooter 显示页脚
func (t *Terminal) DisplayFooter() {
	fmt.Println(DrawSeparator(DisplayWidth, "="))
	fmt.Printf("最后更新: %s | 只读安全模式 | 按 Ctrl+C 安全退出 | 刷新间隔: 5秒\n",
		time.Now().Format("2006-01-02 15:04:05"))
}

// DisplayError 显示错误
func (t *Terminal) DisplayError(msg string, err error) {
	ErrorColor.Printf("[错误] %s: %v\n", msg, err)
}
