package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/Y-vQv-Y/es-monitor/internal/config"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
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

	// 脏页和回写
	if metrics.Memory.Dirty > 0 || metrics.Memory.Writeback > 0 {
		fmt.Printf("  写入状态: 脏页=%s, 回写=%s\n",
			FormatBytesUint64(metrics.Memory.Dirty),
			FormatBytesUint64(metrics.Memory.Writeback))
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
