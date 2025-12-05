package collector

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// SystemCollector 系统指标采集器
type SystemCollector struct {
	prevDiskIO  map[string]disk.IOCountersStat
	prevNetIO   map[string]net.IOCountersStat
	prevTime    time.Time
}

// NewSystemCollector 创建系统采集器
func NewSystemCollector() *SystemCollector {
	return &SystemCollector{
		prevDiskIO: make(map[string]disk.IOCountersStat),
		prevNetIO:  make(map[string]net.IOCountersStat),
		prevTime:   time.Now(),
	}
}

// Collect 采集系统指标
func (c *SystemCollector) Collect(ctx context.Context) (*model.SystemMetrics, error) {
	metrics := &model.SystemMetrics{
		Timestamp: time.Now().Unix(),
	}

	// CPU 指标
	cpuMetrics, err := c.collectCPU()
	if err != nil {
		return nil, err
	}
	metrics.CPU = cpuMetrics

	// 内存指标
	memMetrics, err := c.collectMemory()
	if err != nil {
		return nil, err
	}
	metrics.Memory = memMetrics

	// 磁盘指标
	diskMetrics, err := c.collectDisk()
	if err != nil {
		return nil, err
	}
	metrics.Disk = diskMetrics

	// 网络指标
	netMetrics, err := c.collectNetwork()
	if err != nil {
		return nil, err
	}
	metrics.Network = netMetrics

	return metrics, nil
}

// collectCPU 采集 CPU 指标
func (c *SystemCollector) collectCPU() (model.CPUMetrics, error) {
	cpuMetrics := model.CPUMetrics{}

	// CPU 核心数
	counts, err := cpu.Counts(true)
	if err != nil {
		return cpuMetrics, err
	}
	cpuMetrics.Cores = counts

	// CPU 使用率
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return cpuMetrics, err
	}
	if len(percentages) > 0 {
		cpuMetrics.UsagePercent = percentages[0]
	}

	// 详细 CPU 时间
	times, err := cpu.Times(false)
	if err == nil && len(times) > 0 {
		total := times[0].Total()
		if total > 0 {
			cpuMetrics.UserPercent = times[0].User / total * 100
			cpuMetrics.SystemPercent = times[0].System / total * 100
			cpuMetrics.IdlePercent = times[0].Idle / total * 100
			cpuMetrics.IOWaitPercent = times[0].Iowait / total * 100
		}
	}

	return cpuMetrics, nil
}

// collectMemory 采集内存指标
func (c *SystemCollector) collectMemory() (model.MemoryMetrics, error) {
	memMetrics := model.MemoryMetrics{}

	// 虚拟内存
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return memMetrics, err
	}

	memMetrics.Total = vmStat.Total
	memMetrics.Available = vmStat.Available
	memMetrics.Used = vmStat.Used
	memMetrics.UsedPercent = vmStat.UsedPercent
	memMetrics.Free = vmStat.Free
	memMetrics.Buffers = vmStat.Buffers
	memMetrics.Cached = vmStat.Cached

	// Swap 内存
	swapStat, err := mem.SwapMemory()
	if err == nil {
		memMetrics.SwapTotal = swapStat.Total
		memMetrics.SwapUsed = swapStat.Used
		memMetrics.SwapFree = swapStat.Free
	}

	return memMetrics, nil
}

// collectDisk 采集磁盘指标
func (c *SystemCollector) collectDisk() (model.DiskMetrics, error) {
	diskMetrics := model.DiskMetrics{}

	// 磁盘 IO 统计
	ioCounters, err := disk.IOCounters()
	if err == nil {
		now := time.Now()
		elapsed := now.Sub(c.prevTime).Seconds()

		if elapsed > 0 {
			var totalReadBytes, totalWriteBytes, totalReadTime, totalWriteTime float64
			var totalReadOps, totalWriteOps float64

			for name, counter := range ioCounters {
				if prev, ok := c.prevDiskIO[name]; ok {
					readBytes := float64(counter.ReadBytes - prev.ReadBytes)
					writeBytes := float64(counter.WriteBytes - prev.WriteBytes)
					readOps := float64(counter.ReadCount - prev.ReadCount)
					writeOps := float64(counter.WriteCount - prev.WriteCount)
					readTime := float64(counter.ReadTime - prev.ReadTime)
					writeTime := float64(counter.WriteTime - prev.WriteTime)

					totalReadBytes += readBytes / elapsed
					totalWriteBytes += writeBytes / elapsed
					totalReadOps += readOps / elapsed
					totalWriteOps += writeOps / elapsed
					totalReadTime += readTime / elapsed
					totalWriteTime += writeTime / elapsed
				}
				c.prevDiskIO[name] = counter
			}

			diskMetrics.ReadBytesPerSec = totalReadBytes
			diskMetrics.WriteBytesPerSec = totalWriteBytes
			diskMetrics.ReadOpsPerSec = totalReadOps
			diskMetrics.WriteOpsPerSec = totalWriteOps

			// 计算 IO 利用率（基于时间）
			totalIOTime := totalReadTime + totalWriteTime
			diskMetrics.IOUtilPercent = (totalIOTime / (elapsed * 1000)) * 100 // ms to %
			if diskMetrics.IOUtilPercent > 100 {
				diskMetrics.IOUtilPercent = 100
			}
		}

		c.prevTime = now
	}

	// 磁盘分区使用情况
	partitions, err := disk.Partitions(false)
	if err == nil {
		for _, partition := range partitions {
			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				continue
			}

			partMetric := model.PartitionMetrics{
				Device:      partition.Device,
				Mountpoint:  partition.Mountpoint,
				Total:       usage.Total,
				Used:        usage.Used,
				Free:        usage.Free,
				UsedPercent: usage.UsedPercent,
			}

			// 分区 IO 详细（新增）
			if counter, ok := ioCounters[partition.Device]; ok {
				if prev, ok := c.prevDiskIO[partition.Device]; ok && elapsed > 0 {
					readBytes := float64(counter.ReadBytes - prev.ReadBytes)
					writeBytes := float64(counter.WriteBytes - prev.WriteBytes)
					readOps := float64(counter.ReadCount - prev.ReadCount)
					writeOps := float64(counter.WriteCount - prev.WriteCount)
					readTime := float64(counter.ReadTime - prev.ReadTime)
					writeTime := float64(counter.WriteTime - prev.WriteTime)

					partMetric.IOStats.ReadBytesPerSec = readBytes / elapsed
					partMetric.IOStats.WriteBytesPerSec = writeBytes / elapsed
					partMetric.IOStats.ReadOpsPerSec = readOps / elapsed
					partMetric.IOStats.WriteOpsPerSec = writeOps / elapsed
					totalIOTime := readTime + writeTime
					partMetric.IOStats.IOUtilPercent = (totalIOTime / (elapsed * 1000)) * 100
					if partMetric.IOStats.IOUtilPercent > 100 {
						partMetric.IOStats.IOUtilPercent = 100
					}
				}
			}

			diskMetrics.Partitions = append(diskMetrics.Partitions, partMetric)
		}
	}

	return diskMetrics, nil
}

// collectNetwork 采集网络指标
func (c *SystemCollector) collectNetwork() (model.NetworkMetrics, error) {
	netMetrics := model.NetworkMetrics{}

	// 网络 IO 统计
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return netMetrics, err
	}

	now := time.Now()
	elapsed := now.Sub(c.prevTime).Seconds()

	if elapsed > 0 {
		var totalBytesSent, totalBytesRecv float64
		var totalPacketsSent, totalPacketsRecv float64

		for _, counter := range ioCounters {
			if counter.Name == "lo" {
				continue // 跳过回环接口
			}

			if prev, ok := c.prevNetIO[counter.Name]; ok {
				bytesSent := float64(counter.BytesSent - prev.BytesSent)
				bytesRecv := float64(counter.BytesRecv - prev.BytesRecv)
				packetsSent := float64(counter.PacketsSent - prev.PacketsSent)
				packetsRecv := float64(counter.PacketsRecv - prev.PacketsRecv)

				totalBytesSent += bytesSent / elapsed
				totalBytesRecv += bytesRecv / elapsed
				totalPacketsSent += packetsSent / elapsed
				totalPacketsRecv += packetsRecv / elapsed

				// 记录每个网卡的指标（增强瞬时速率）
				ifMetric := model.InterfaceMetrics{
					Name:              counter.Name,
					BytesSent:         counter.BytesSent,
					BytesRecv:         counter.BytesRecv,
					PacketsSent:       counter.PacketsSent,
					PacketsRecv:       counter.PacketsRecv,
					BytesSentPerSec:   bytesSent / elapsed,
					BytesRecvPerSec:   bytesRecv / elapsed,
					PacketsSentPerSec: packetsSent / elapsed,
					PacketsRecvPerSec: packetsRecv / elapsed,
				}
				netMetrics.Interfaces = append(netMetrics.Interfaces, ifMetric)
			}
			c.prevNetIO[counter.Name] = counter
		}

		netMetrics.BytesSentPerSec = totalBytesSent
		netMetrics.BytesRecvPerSec = totalBytesRecv
		netMetrics.PacketsSentPerSec = totalPacketsSent
		netMetrics.PacketsRecvPerSec = totalPacketsRecv
	}

	c.prevTime = now
	return netMetrics, nil
}
