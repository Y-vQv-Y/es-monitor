package collector

import (
	"context"
	"strings"
	"time"

	"github.com/Y-vQv-Y/es-monitor/internal/model"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// SystemCollector 系统指标采集器
type SystemCollector struct {
	prevDiskIO       map[string]disk.IOCountersStat
	prevNetIO        map[string]net.IOCountersStat
	prevTime         time.Time
	initialized      bool
	
	// 网络流量平滑处理
	netHistoryWindow []NetworkSnapshot
	maxHistorySize   int
}

type NetworkSnapshot struct {
	Timestamp       time.Time
	BytesSentPerSec float64
	BytesRecvPerSec float64
}

func NewSystemCollector() *SystemCollector {
	return &SystemCollector{
		prevDiskIO:       make(map[string]disk.IOCountersStat),
		prevNetIO:        make(map[string]net.IOCountersStat),
		prevTime:         time.Now(),
		initialized:      false,
		netHistoryWindow: make([]NetworkSnapshot, 0),
		maxHistorySize:   10,
	}
}

func (c *SystemCollector) Collect(ctx context.Context) (*model.SystemMetrics, error) {
	metrics := &model.SystemMetrics{
		Timestamp: time.Now().Unix(),
	}

	cpuMetrics, err := c.collectCPU()
	if err != nil {
		cpuMetrics = model.CPUMetrics{}
	}
	metrics.CPU = cpuMetrics

	memMetrics, err := c.collectMemory()
	if err != nil {
		memMetrics = model.MemoryMetrics{}
	}
	metrics.Memory = memMetrics

	diskMetrics, err := c.collectDisk()
	if err != nil {
		diskMetrics = model.DiskMetrics{}
	}
	metrics.Disk = diskMetrics

	netMetrics, err := c.collectNetwork()
	if err != nil {
		netMetrics = model.NetworkMetrics{}
	}
	metrics.Network = netMetrics

	c.initialized = true
	return metrics, nil
}

func (c *SystemCollector) collectCPU() (model.CPUMetrics, error) {
	cpuMetrics := model.CPUMetrics{}

	counts, err := cpu.Counts(true)
	if err == nil {
		cpuMetrics.LogicalCores = counts
	}

	physCounts, err := cpu.Counts(false)
	if err == nil {
		cpuMetrics.Cores = physCounts
	}

	percentages, err := cpu.Percent(100*time.Millisecond, false)
	if err == nil && len(percentages) > 0 {
		cpuMetrics.UsagePercent = percentages[0]
	}

	perCPU, err := cpu.Percent(100*time.Millisecond, true)
	if err == nil {
		cpuMetrics.PerCPUPercent = perCPU
	}

	times, err := cpu.Times(false)
	if err == nil && len(times) > 0 {
		total := times[0].User + times[0].System + times[0].Idle + times[0].Iowait +
			times[0].Irq + times[0].Softirq + times[0].Steal + times[0].Guest

		if total > 0 {
			cpuMetrics.UserPercent = times[0].User / total * 100
			cpuMetrics.SystemPercent = times[0].System / total * 100
			cpuMetrics.IdlePercent = times[0].Idle / total * 100
			cpuMetrics.IOWaitPercent = times[0].Iowait / total * 100
			cpuMetrics.IrqPercent = times[0].Irq / total * 100
			cpuMetrics.SoftIrqPercent = times[0].Softirq / total * 100
			cpuMetrics.StealPercent = times[0].Steal / total * 100
			cpuMetrics.GuestPercent = times[0].Guest / total * 100
		}
	}

	avgStat, err := load.Avg()
	if err == nil {
		cpuMetrics.LoadAvg1 = avgStat.Load1
		cpuMetrics.LoadAvg5 = avgStat.Load5
		cpuMetrics.LoadAvg15 = avgStat.Load15
	}

	return cpuMetrics, nil
}

func (c *SystemCollector) collectMemory() (model.MemoryMetrics, error) {
	memMetrics := model.MemoryMetrics{}

	vmStat, err := mem.VirtualMemory()
	if err == nil {
		memMetrics.Total = vmStat.Total
		memMetrics.Available = vmStat.Available
		memMetrics.Used = vmStat.Used
		memMetrics.UsedPercent = vmStat.UsedPercent
		memMetrics.Free = vmStat.Free
		memMetrics.Buffers = vmStat.Buffers
		memMetrics.Cached = vmStat.Cached
		memMetrics.Shared = vmStat.Shared
		memMetrics.Active = vmStat.Active
		memMetrics.Inactive = vmStat.Inactive

		if vmStat.Dirty > 0 {
			memMetrics.Dirty = vmStat.Dirty
		}
		
		memMetrics.Mapped = vmStat.Mapped
		memMetrics.Slab = vmStat.Slab
	}

	swapStat, err := mem.SwapMemory()
	if err == nil {
		memMetrics.SwapTotal = swapStat.Total
		memMetrics.SwapUsed = swapStat.Used
		memMetrics.SwapFree = swapStat.Free
		if swapStat.Total > 0 {
			memMetrics.SwapUsedPercent = float64(swapStat.Used) / float64(swapStat.Total) * 100
		}
		memMetrics.PageIn = swapStat.Sin
		memMetrics.PageOut = swapStat.Sout
	}

	return memMetrics, nil
}

func (c *SystemCollector) collectDisk() (model.DiskMetrics, error) {
	diskMetrics := model.DiskMetrics{}
	now := time.Now()
	elapsed := now.Sub(c.prevTime).Seconds()

	ioCounters, err := disk.IOCounters()
	if err == nil && elapsed > 0 && c.initialized {
		var totalReadBytes, totalWriteBytes float64
		var totalReadOps, totalWriteOps float64
		var totalReadKB, totalWriteKB uint64

		deviceMetrics := make([]model.DiskDeviceMetrics, 0)

		for name, counter := range ioCounters {
			if prev, ok := c.prevDiskIO[name]; ok {
				readBytes := float64(counter.ReadBytes - prev.ReadBytes)
				writeBytes := float64(counter.WriteBytes - prev.WriteBytes)
				readOps := float64(counter.ReadCount - prev.ReadCount)
				writeOps := float64(counter.WriteCount - prev.WriteCount)
				ioTime := float64(counter.IoTime - prev.IoTime)

				readBytesPerSec := readBytes / elapsed
				writeBytesPerSec := writeBytes / elapsed
				readOpsPerSec := readOps / elapsed
				writeOpsPerSec := writeOps / elapsed

				totalReadBytes += readBytesPerSec
				totalWriteBytes += writeBytesPerSec
				totalReadOps += readOpsPerSec
				totalWriteOps += writeOpsPerSec

				ioUtilPercent := 0.0
				if elapsed > 0 {
					ioUtilPercent = ioTime / (elapsed * 1000) * 100
					if ioUtilPercent > 100 {
						ioUtilPercent = 100
					}
				}

				deviceMetric := model.DiskDeviceMetrics{
					Device:           name,
					ReadBytesPerSec:  readBytesPerSec,
					WriteBytesPerSec: writeBytesPerSec,
					ReadOpsPerSec:    readOpsPerSec,
					WriteOpsPerSec:   writeOpsPerSec,
					IOUtilPercent:    ioUtilPercent,
				}

				totalOpsNow := readOps + writeOps
				if totalOpsNow > 0 {
					deviceMetric.AvgRequestSize = (readBytes + writeBytes) / totalOpsNow
				}

				deviceMetrics = append(deviceMetrics, deviceMetric)
			}

			c.prevDiskIO[name] = counter
			totalReadKB += counter.ReadBytes / 1024
			totalWriteKB += counter.WriteBytes / 1024
		}

		diskMetrics.ReadBytesPerSec = totalReadBytes
		diskMetrics.WriteBytesPerSec = totalWriteBytes
		diskMetrics.ReadOpsPerSec = totalReadOps
		diskMetrics.WriteOpsPerSec = totalWriteOps
		diskMetrics.TotalReadBytes = totalReadKB * 1024
		diskMetrics.TotalWriteBytes = totalWriteKB * 1024
		diskMetrics.TotalReadOps = uint64(totalReadOps)
		diskMetrics.TotalWriteOps = uint64(totalWriteOps)
		diskMetrics.Devices = deviceMetrics
	}

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
				FSType:      partition.Fstype,
				Total:       usage.Total,
				Used:        usage.Used,
				Free:        usage.Free,
				UsedPercent: usage.UsedPercent,
				InodesTotal: usage.InodesTotal,
				InodesUsed:  usage.InodesUsed,
				InodesFree:  usage.InodesFree,
			}
			diskMetrics.Partitions = append(diskMetrics.Partitions, partMetric)
		}
	}

	c.prevTime = now
	return diskMetrics, nil
}

// collectNetwork 采集网络详细指标（修复计算错误）
func (c *SystemCollector) collectNetwork() (model.NetworkMetrics, error) {
	netMetrics := model.NetworkMetrics{}
	now := time.Now()
	elapsed := now.Sub(c.prevTime).Seconds()

	// 【关键修复1】确保时间间隔合理（至少 1 秒）
	if elapsed < 1.0 {
		elapsed = 1.0
	}

	ioCounters, err := net.IOCounters(true)
	if err == nil && c.initialized {
		var totalBytesSent, totalBytesRecv float64
		var totalPacketsSent, totalPacketsRecv float64
		var totalErrors, totalDrops float64
		var totalBytesSentAcc, totalBytesRecvAcc uint64
		var totalPacketsSentAcc, totalPacketsRecvAcc uint64

		interfaceMetrics := make([]model.InterfaceMetrics, 0)

		for _, counter := range ioCounters {
			// 过滤 Loopback 和虚拟网卡
			if counter.Name == "lo" ||
			   strings.HasPrefix(counter.Name, "veth") ||
			   strings.HasPrefix(counter.Name, "calico") ||
			   strings.HasPrefix(counter.Name, "br-") ||
			   strings.HasPrefix(counter.Name, "cni") ||
			   strings.HasPrefix(counter.Name, "flannel") ||
			   strings.HasPrefix(counter.Name, "tunl") ||
			   strings.HasPrefix(counter.Name, "vlan") ||
			   strings.HasPrefix(counter.Name, "vxlan") ||
			   strings.HasPrefix(counter.Name, "virb") ||
			   strings.HasPrefix(counter.Name, "virbr0-") ||
			   strings.HasPrefix(counter.Name, "kube-ipvs") ||
			   strings.HasPrefix(counter.Name, "docker") {
				continue
			}

			if prev, ok := c.prevNetIO[counter.Name]; ok {
				// 【关键修复2】计算增量，防止计数器溢出
				bytesSent := float64(counter.BytesSent - prev.BytesSent)
				bytesRecv := float64(counter.BytesRecv - prev.BytesRecv)
				packetsSent := float64(counter.PacketsSent - prev.PacketsSent)
				packetsRecv := float64(counter.PacketsRecv - prev.PacketsRecv)
				errors := float64((counter.Errin + counter.Errout) - (prev.Errin + prev.Errout))
				drops := float64((counter.Dropin + counter.Dropout) - (prev.Dropin + prev.Dropout))

				// 【关键修复3】检测计数器重置（负数）
				if bytesSent < 0 { bytesSent = 0 }
				if bytesRecv < 0 { bytesRecv = 0 }
				if packetsSent < 0 { packetsSent = 0 }
				if packetsRecv < 0 { packetsRecv = 0 }
				if errors < 0 { errors = 0 }
				if drops < 0 { drops = 0 }

				// 计算速率
				bytesSentPerSec := bytesSent / elapsed
				bytesRecvPerSec := bytesRecv / elapsed
				packetsSentPerSec := packetsSent / elapsed
				packetsRecvPerSec := packetsRecv / elapsed
				errorsPerSec := errors / elapsed
				dropsPerSec := drops / elapsed

				// 【关键修复4】异常值检测（速率 > 1 GB/s 肯定有问题）
				maxReasonableRate := 1.0 * 1024 * 1024 * 1024 // 1 GB/s
				if bytesSentPerSec > maxReasonableRate {
					bytesSentPerSec = 0
				}
				if bytesRecvPerSec > maxReasonableRate {
					bytesRecvPerSec = 0
				}

				totalBytesSent += bytesSentPerSec
				totalBytesRecv += bytesRecvPerSec
				totalPacketsSent += packetsSentPerSec
				totalPacketsRecv += packetsRecvPerSec
				totalErrors += errorsPerSec
				totalDrops += dropsPerSec

				ifaceMetric := model.InterfaceMetrics{
					Name:              counter.Name,
					BytesSentPerSec:   bytesSentPerSec,
					BytesRecvPerSec:   bytesRecvPerSec,
					PacketsSentPerSec: packetsSentPerSec,
					PacketsRecvPerSec: packetsRecvPerSec,
					BytesSent:         counter.BytesSent,
					BytesRecv:         counter.BytesRecv,
					PacketsSent:       counter.PacketsSent,
					PacketsRecv:       counter.PacketsRecv,
					ErrorsIn:          counter.Errin,
					ErrorsOut:         counter.Errout,
					DropsIn:           counter.Dropin,
					DropsOut:          counter.Dropout,
				}

				interfaceMetrics = append(interfaceMetrics, ifaceMetric)
			}

			// 更新历史数据
			c.prevNetIO[counter.Name] = counter

			// 累计总量
			totalBytesSentAcc += counter.BytesSent
			totalBytesRecvAcc += counter.BytesRecv
			totalPacketsSentAcc += counter.PacketsSent
			totalPacketsRecvAcc += counter.PacketsRecv
		}

		// 【关键修复5】使用中位数平滑，排除异常峰值
		smoothedSent, smoothedRecv := c.smoothNetworkTraffic(
			totalBytesSent, 
			totalBytesRecv,
			now,
		)

		netMetrics.BytesSentPerSec = smoothedSent
		netMetrics.BytesRecvPerSec = smoothedRecv
		netMetrics.PacketsSentPerSec = totalPacketsSent
		netMetrics.PacketsRecvPerSec = totalPacketsRecv
		netMetrics.ErrorsPerSec = totalErrors
		netMetrics.DropsPerSec = totalDrops
		netMetrics.TotalBytesSent = totalBytesSentAcc
		netMetrics.TotalBytesRecv = totalBytesRecvAcc
		netMetrics.TotalPacketsSent = totalPacketsSentAcc
		netMetrics.TotalPacketsRecv = totalPacketsRecvAcc
		netMetrics.Interfaces = interfaceMetrics
	}

	return netMetrics, nil
}

// smoothNetworkTraffic 平滑网络流量数据
func (c *SystemCollector) smoothNetworkTraffic(
	currentSent, currentRecv float64,
	timestamp time.Time,
) (smoothedSent, smoothedRecv float64) {
	// 添加当前快照
	snapshot := NetworkSnapshot{
		Timestamp:       timestamp,
		BytesSentPerSec: currentSent,
		BytesRecvPerSec: currentRecv,
	}
	c.netHistoryWindow = append(c.netHistoryWindow, snapshot)

	// 保持窗口大小
	if len(c.netHistoryWindow) > c.maxHistorySize {
		c.netHistoryWindow = c.netHistoryWindow[1:]
	}

	// 历史数据不足，直接返回当前值
	if len(c.netHistoryWindow) < 3 {
		return currentSent, currentRecv
	}

	// 提取历史值
	sentValues := make([]float64, len(c.netHistoryWindow))
	recvValues := make([]float64, len(c.netHistoryWindow))
	for i, snap := range c.netHistoryWindow {
		sentValues[i] = snap.BytesSentPerSec
		recvValues[i] = snap.BytesRecvPerSec
	}

	// 使用中位数（抗异常值）
	smoothedSent = median(sentValues)
	smoothedRecv = median(recvValues)

	// 额外保护：如果中位数仍然异常大（> 100 MB/s），使用最小值
	if smoothedSent > 100*1024*1024 {
		smoothedSent = minimum(sentValues)
	}
	if smoothedRecv > 100*1024*1024 {
		smoothedRecv = minimum(recvValues)
	}

	return smoothedSent, smoothedRecv
}

// minimum 计算最小值
func minimum(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// median 计算中位数
func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

// average 计算平均值
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
