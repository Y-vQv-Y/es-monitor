package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

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

// DisplayNetworkMetrics 显示网络详细信息（完整版）
func (t *Terminal) DisplayNetworkMetrics(metrics *model.NetworkMetrics) {
	LabelColor.Println("【网络详细信息 - 实时吞吐量】")
	fmt.Println(DrawSeparator(90, "."))

	// 网络总体吞吐量
	fmt.Printf("  总体吞吐: 发送=%s, 接收=%s\n",
		FormatBytesPerSec(metrics.BytesSentPerSec),
		FormatBytesPerSec(metrics.BytesRecvPerSec))

	// 转换为 Mbps 显示
	sendMbps := metrics.BytesSentPerSec * 8 / 1024 / 1024
	recvMbps := metrics.BytesRecvPerSec * 8 / 1024 / 1024
	fmt.Printf("              (发送=%.2f Mbps, 接收=%.2f Mbps)\n", sendMbps, recvMbps)

	// 数据包速率
	fmt.Printf("  数据包率: 发送=%.1f pkt/s, 接收=%.1f pkt/s\n",
		metrics.PacketsSentPerSec,
		metrics.PacketsRecvPerSec)

	// 错误和丢包
	if metrics.ErrorsPerSec > 0 || metrics.DropsPerSec > 0 {
		if metrics.ErrorsPerSec > 0 {
			StatusRed.Printf("  错误率: %.2f errors/s [需要检查网络质量]\n", metrics.ErrorsPerSec)
		}
		if metrics.DropsPerSec > 0 {
			StatusYellow.Printf("  丢包率: %.2f drops/s [可能存在网络拥塞]\n", metrics.DropsPerSec)
		}
	}

	// 累计流量统计
	if metrics.TotalBytesSent > 0 || metrics.TotalBytesRecv > 0 {
		fmt.Printf("  累计流量: 发送=%s, 接收=%s\n",
			FormatBytesUint64(metrics.TotalBytesSent),
			FormatBytesUint64(metrics.TotalBytesRecv))
	}

	// TCP 连接统计
	if metrics.TCPConnections > 0 {
		fmt.Printf("  TCP 连接: 总计=%d, 已建立=%d, 监听=%d, TIME_WAIT=%d\n",
			metrics.TCPConnections,
			metrics.TCPEstablished,
			metrics.TCPListening,
			metrics.TCPTimeWait)

		// TIME_WAIT 连接过多警告
		if metrics.TCPTimeWait > 1000 {
			StatusYellow.Printf("  [警告: TIME_WAIT 连接数较多 (%d)，可能需要调整系统参数]\n",
				metrics.TCPTimeWait)
		}
	}

	// 各网卡详细信息
	if len(metrics.Interfaces) > 0 {
		fmt.Println()
		fmt.Println("  各网卡详情:")
		fmt.Printf("  %-15s %15s %15s %12s %12s %8s %8s\n",
			"网卡", "发送速率", "接收速率", "发送pkt/s", "接收pkt/s", "发送错误", "接收错误")
		fmt.Println("  " + DrawSeparator(95, "-"))

		for _, iface := range metrics.Interfaces {
			// 跳过回环和没有流量的网卡
			if iface.Name == "lo" ||
				(iface.BytesSentPerSec == 0 && iface.BytesRecvPerSec == 0) {
				continue
			}

			ifaceName := TruncateString(iface.Name, 15)
			fmt.Printf("  %-15s %15s %15s %12.1f %12.1f %8d %8d",
				ifaceName,
				FormatBytesPerSec(iface.BytesSentPerSec),
				FormatBytesPerSec(iface.BytesRecvPerSec),
				iface.PacketsSentPerSec,
				iface.PacketsRecvPerSec,
				iface.ErrorsOut,
				iface.ErrorsIn)

			// 显示错误警告
			if iface.ErrorsIn > 0 || iface.ErrorsOut > 0 {
				StatusRed.Print(" [有错误]")
			}
			if iface.DropsIn > 0 || iface.DropsOut > 0 {
				StatusYellow.Print(" [有丢包]")
			}
			fmt.Println()

			// 显示网卡状态（如果可用）
			if iface.Speed > 0 {
				fmt.Printf("              状态: %s, 速度: %d Mbps, MTU: %d\n",
					func() string {
						if iface.IsUp {
							return "UP"
						}
						return "DOWN"
					}(),
					iface.Speed,
					iface.MTU)
			}
		}
	}

	// 网络性能评估
	totalBandwidth := metrics.BytesSentPerSec + metrics.BytesRecvPerSec
	totalMbps := totalBandwidth * 8 / 1024 / 1024

	fmt.Println()
	if totalMbps > 800 {
		StatusRed.Printf("  [网络负载评估: 繁忙 - 总带宽 %.2f Mbps]\n", totalMbps)
	} else if totalMbps > 500 {
		StatusYellow.Printf("  [网络负载评估: 正常偏高 - 总带宽 %.2f Mbps]\n", totalMbps)
	} else {
		StatusGreen.Printf("  [网络负载评估: 正常 - 总带宽 %.2f Mbps]\n", totalMbps)
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
			ValueColor.Sprint(node.Indices.Docs.Count),
			ValueColor.Sprint(FormatBytes(node.Indices.Store.SizeInBytes)))

		// 实时速率
		if prev, ok := prevData[nodeID]; ok {
			elapsed := time.Since(prev.Timestamp).Seconds()
			if elapsed > 0 {
				indexRate := float64(node.Indices.Indexing.IndexTotal-prev.IndexTotal) / elapsed
				queryRate := float64(node.Indices.Search.QueryTotal-prev.QueryTotal) / elapsed

				fmt.Printf("  写入速率: %s\n", ValueColor.Sprint(FormatRate(indexRate, "docs/s")))
				fmt.Printf("  查询速率: %s\n", ValueColor.Sprint(FormatRate(queryRate, "queries/s")))
			}
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

	// 显示表头
	fmt.Printf("%-30s %-8s %-8s %-12s %-12s\n",
		"索引名称", "状态", "分片", "文档数", "大小")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	count := 0
	for _, idx := range indices {
		if count >= 20 {
			break
		}
		count++

		statusColor := GetStatusColor(idx.Health)
		indexName := TruncateString(idx.Index, 30)

		fmt.Printf("%-30s ", indexName)
		statusColor.Printf("%-8s ", strings.ToUpper(idx.Health))
		fmt.Printf("%-8s %-12s %-12s\n",
			idx.Pri+"/"+idx.Rep,
			idx.DocsCount,
			idx.StoreSize)
	}

	if len(indices) > 20 {
		fmt.Printf("\n  ... 还有 %d 个索引未显示\n", len(indices)-20)
	}

	fmt.Println()
}

// DisplayFooter 显示页脚
func (t *Terminal) DisplayFooter() {
	fmt.Println(DrawSeparator(DisplayWidth, "="))
	fmt.Printf("最后更新: %s | 只读安全模式 | 按 Ctrl+C 安全退出\n",
		time.Now().Format("2006-01-02 15:04:05"))
}

// DisplayError 显示错误
func (t *Terminal) DisplayError(msg string, err error) {
	ErrorColor.Printf("[错误] %s: %v\n", msg, err)
}
