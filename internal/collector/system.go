package collector

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// SystemCollector 系统指标采集器（完整版，生产环境安全）
type SystemCollector struct {
	prevDiskIO  map[string]disk.IOCountersStat
	prevNetIO   map[string]net.IOCountersStat
	prevTime    time.Time
	initialized bool
}

// NewSystemCollector 创建系统采集器
func NewSystemCollector() *SystemCollector {
	return &SystemCollector{
		prevDiskIO:  make(map[string]disk.IOCountersStat),
		prevNetIO:   make(map[string]net.IOCountersStat),
		prevTime:    time.Now(),
		initialized: false,
	}
}

// Collect 采集系统指标（只读操作，不影响系统性能）
func (c *SystemCollector) Collect(ctx context.Context) (*model.SystemMetrics, error) {
	metrics := &model.SystemMetrics{
		Timestamp: time.Now().Unix(),
	}

	// CPU 指标采集
	cpuMetrics, err := c.collectCPU()
	if err != nil {
		// 采集失败不影响其他指标，记录错误但继续
		cpuMetrics = model.CPUMetrics{}
	}
	metrics.CPU = cpuMetrics

	// 内存指标采集
	memMetrics, err := c.collectMemory()
	if err != nil {
		memMetrics = model.MemoryMetrics{}
	}
	metrics.Memory = memMetrics

	// 磁盘指标采集
	diskMetrics, err := c.collectDisk()
	if err != nil {
		diskMetrics = model.DiskMetrics{}
	}
	metrics.Disk = diskMetrics

	// 网络指标采集
	netMetrics, err := c.collectNetwork()
	if err != nil {
		netMetrics = model.NetworkMetrics{}
	}
	metrics.Network = netMetrics

	c.initialized = true
	return metrics, nil
}

// collectCPU 采集 CPU 详细指标（生产环境安全）
func (c *SystemCollector) collectCPU() (model.CPUMetrics, error) {
	cpuMetrics := model.CPUMetrics{}

	// CPU 核心数（物理和逻辑）
	counts, err := cpu.Counts(true)
	if err == nil {
		cpuMetrics.LogicalCores = counts
	}
	
	physCounts, err := cpu.Counts(false)
	if err == nil {
		cpuMetrics.Cores = physCounts
	}

	// 总体 CPU 使用率（使用短时间间隔，避免阻塞）
	percentages, err := cpu.Percent(100*time.Millisecond, false)
	if err == nil && len(percentages) > 0 {
		cpuMetrics.UsagePercent = percentages[0]
	}

	// 每个 CPU 核心的使用率
	perCPU, err := cpu.Percent(100*time.Millisecond, true)
	if err == nil {
		cpuMetrics.PerCPUPercent = perCPU
	}

	// 详细 CPU 时间统计
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

	// 系统负载
	avgStat, err := load.Avg()
	if err == nil {
		cpuMetrics.LoadAvg1 = avgStat.Load1
		cpuMetrics.LoadAvg5 = avgStat.Load5
		cpuMetrics.LoadAvg15 = avgStat.Load15
	}

	return cpuMetrics, nil
}

// collectMemory 采集内存详细指标（只读操作）
func (c *SystemCollector) collectMemory() (model.MemoryMetrics, error) {
	memMetrics := model.MemoryMetrics{}

	// 虚拟内存统计
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
		
		// 某些平台可能没有这些字段，使用条件编译或检查
		// Dirty 和 Writeback 在 Linux 上可用
		// 其他平台会设置为 0
		if vmStat.Dirty > 0 {
			memMetrics.Dirty = vmStat.Dirty
		}
		// Writeback 字段在某些版本的 gopsutil 中不存在
		// 我们跳过它或使用反射检查
		
		memMetrics.Mapped = vmStat.Mapped
		memMetrics.Slab = vmStat.Slab
	}

	// Swap 内存统计
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

// collectDisk 采集磁盘详细指标（只读操作）
func (c *SystemCollector) collectDisk() (model.DiskMetrics, error) {
	diskMetrics := model.DiskMetrics{}
	now := time.Now()
	elapsed := now.Sub(c.prevTime).Seconds()

	// 磁盘 IO 统计
	ioCounters, err := disk.IOCounters()
	if err == nil && elapsed > 0 && c.initialized {
		var totalReadBytes, totalWriteBytes float64
		var totalReadOps, totalWriteOps float64
		var totalReadKB, totalWriteKB uint64
		
		deviceMetrics := make([]model.DiskDeviceMetrics, 0)

		for name, counter := range ioCounters {
			if prev, ok := c.prevDiskIO[name]; ok {
				// 计算增量
				readBytes := float64(counter.ReadBytes - prev.ReadBytes)
				writeBytes := float64(counter.WriteBytes - prev.WriteBytes)
				readOps := float64(counter.ReadCount - prev.ReadCount)
				writeOps := float64(counter.WriteCount - prev.WriteCount)
				ioTime := float64(counter.IoTime - prev.IoTime)

				// 计算速率
				readBytesPerSec := readBytes / elapsed
				writeBytesPerSec := writeBytes / elapsed
				readOpsPerSec := readOps / elapsed
				writeOpsPerSec := writeOps / elapsed

				totalReadBytes += readBytesPerSec
				totalWriteBytes += writeBytesPerSec
				totalReadOps += readOpsPerSec
				totalWriteOps += writeOpsPerSec

				// IO 使用率
				ioUtilPercent := 0.0
				if elapsed > 0 {
					ioUtilPercent = ioTime / (elapsed * 1000) * 100
					if ioUtilPercent > 100 {
						ioUtilPercent = 100
					}
				}

				// 每个设备的详细指标
				deviceMetric := model.DiskDeviceMetrics{
					Device:           name,
					ReadBytesPerSec:  readBytesPerSec,
					WriteBytesPerSec: writeBytesPerSec,
					ReadOpsPerSec:    readOpsPerSec,
					WriteOpsPerSec:   writeOpsPerSec,
					IOUtilPercent:    ioUtilPercent,
				}
				
				// 平均请求大小
				totalOpsNow := readOps + writeOps
				if totalOpsNow > 0 {
					deviceMetric.AvgRequestSize = (readBytes + writeBytes) / totalOpsNow
				}
				
				deviceMetrics = append(deviceMetrics, deviceMetric)
			}
			
			// 更新历史数据
			c.prevDiskIO[name] = counter
			
			// 累计总量
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

// collectNetwork 采集网络详细指标（只读操作）
// helper: 判断是否为虚拟/overlay 网卡（跳过）
func isVirtualNIC(name string) bool {
	if name == "lo" {
		return true
	}
	prefixes := []string{
		"veth", "docker", "br-", "cni", "flannel", "kube-ipvs", "tunl", "virbr",
		"vlan", "vxlan",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// helper: 判断是否像物理网卡（ens/eth/bond）
func isPhysicalNIC(name string) bool {
	prefixes := []string{"ens", "eth", "eno", "enp", "bond", "bond0", "em"}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// 读取 /proc/<pid>/net/dev 的简单解析（返回 map[iface] -> (bytesSent, bytesRecv)）
// 返回值是粗略估算（和 /proc/net/dev 格式一致）。
func readProcPIDNetDev(pid int) (map[string]struct{ BytesRecv, BytesSent uint64 }, error) {
	result := make(map[string]struct{ BytesRecv, BytesSent uint64 })
	path := fmt.Sprintf("/proc/%d/net/dev", pid)
	f, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNo++
		// 前两行是 header
		if lineNo <= 2 {
			continue
		}
		// 格式： iface: bytes    packets ... bytes ... 
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		iface := strings.TrimSpace(parts[0])
		fields := strings.Fields(strings.TrimSpace(parts[1]))
		if len(fields) < 10 {
			continue
		}
		recvBytes, _ := strconv.ParseUint(fields[0], 10, 64)
		sendBytes, _ := strconv.ParseUint(fields[8], 10, 64)
		result[iface] = struct{ BytesRecv, BytesSent uint64 }{BytesRecv: recvBytes, BytesSent: sendBytes}
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}

// collectNetwork 采集网络详细指标（改进版）
// 参数 deductSelf: 是否尝试扣除当前进程（或 providedPID）的网络增量（若无法获取会跳过扣除）
// 参数 providedPID: 若 >0 则用于指定要扣除的进程 PID（否则使用当前进程 PID）
func (c *SystemCollector) collectNetwork(deductSelf bool, providedPID int) (model.NetworkMetrics, error) {
	netMetrics := model.NetworkMetrics{}
	now := time.Now()
	elapsed := now.Sub(c.prevTime).Seconds()
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		// 更新 prevTime 为当前，避免下一次 elapsed 为很大值
		c.prevTime = now
		return netMetrics, err
	}

	// 如果需要扣除自身，先读取采集开始时的进程 net/dev（粗略）
	var pid int
	if providedPID > 0 {
		pid = providedPID
	} else {
		pid = os.Getpid()
	}

	var beforeSelf map[string]struct{ BytesRecv, BytesSent uint64 }
	var afterSelf map[string]struct{ BytesRecv, BytesSent uint64 }
	if deductSelf {
		// read baseline (best-effort)
		beforeSelf, _ = readProcPIDNetDev(pid)
	}

	// 构建临时 map 便于按接口查找 previous counters
	if c.prevNetIO == nil {
		c.prevNetIO = make(map[string]net.IOCountersStat)
	}

	var totalBytesSentPerSec, totalBytesRecvPerSec float64
	var totalPacketsSentPerSec, totalPacketsRecvPerSec float64
	var totalErrorsPerSec, totalDropsPerSec float64
	var totalBytesSentAcc, totalBytesRecvAcc uint64
	var totalPacketsSentAcc, totalPacketsRecvAcc uint64

	interfaceMetrics := make([]model.InterfaceMetrics, 0)

	// 遍历接口，先计算每个接口增量（但仅把“物理”接口计入 total）
	for _, counter := range ioCounters {
		prev, hasPrev := c.prevNetIO[counter.Name]
		// 计算累计值（始终更新）
		totalBytesSentAcc += counter.BytesSent
		totalBytesRecvAcc += counter.BytesRecv
		totalPacketsSentAcc += counter.PacketsSent
		totalPacketsRecvAcc += counter.PacketsRecv

		// 计算速率（只有在已初始化并且有 prev 值时）
		if hasPrev && c.initialized && elapsed > 0 {
			bytesSent := float64(0)
			bytesRecv := float64(0)
			packetsSent := float64(0)
			packetsRecv := float64(0)
			errors := float64(0)
			drops := float64(0)

			if counter.BytesSent >= prev.BytesSent {
				bytesSent = float64(counter.BytesSent - prev.BytesSent)
			} else {
				// 计数器回绕（rare），用当前值处理
				bytesSent = float64(counter.BytesSent)
			}
			if counter.BytesRecv >= prev.BytesRecv {
				bytesRecv = float64(counter.BytesRecv - prev.BytesRecv)
			} else {
				bytesRecv = float64(counter.BytesRecv)
			}
			if counter.PacketsSent >= prev.PacketsSent {
				packetsSent = float64(counter.PacketsSent - prev.PacketsSent)
			} else {
				packetsSent = float64(counter.PacketsSent)
			}
			if counter.PacketsRecv >= prev.PacketsRecv {
				packetsRecv = float64(counter.PacketsRecv - prev.PacketsRecv)
			} else {
				packetsRecv = float64(counter.PacketsRecv)
			}
			errors = float64((counter.Errin + counter.Errout) - (prev.Errin + prev.Errout))
			drops = float64((counter.Dropin + counter.Dropout) - (prev.Dropin + prev.Dropout))

			bytesSentPerSec := bytesSent / elapsed
			bytesRecvPerSec := bytesRecv / elapsed
			packetsSentPerSec := packetsSent / elapsed
			packetsRecvPerSec := packetsRecv / elapsed
			errorsPerSec := errors / elapsed
			dropsPerSec := drops / elapsed

			// 如果是物理网卡则计入 total（避免 overlay/veth 重复计算）
			if isPhysicalNIC(counter.Name) {
				totalBytesSentPerSec += bytesSentPerSec
				totalBytesRecvPerSec += bytesRecvPerSec
				totalPacketsSentPerSec += packetsSentPerSec
				totalPacketsRecvPerSec += packetsRecvPerSec
				totalErrorsPerSec += errorsPerSec
				totalDropsPerSec += dropsPerSec
			}

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

		// 更新历史（始终更新，供下一轮使用）
		c.prevNetIO[counter.Name] = counter
	}

	// 若开启 deductSelf，在收集结束后读一次 /proc/<pid>/net/dev，计算本进程的接口增量并从 total 中扣除（best-effort）
	if deductSelf {
		afterSelf, _ = readProcPIDNetDev(pid)
		// 计算 delta 并按物理网卡匹配扣除（仅当 beforeSelf/afterSelf 有值时）
		if beforeSelf != nil && afterSelf != nil {
			var selfSendDelta, selfRecvDelta uint64
			for iface, after := range afterSelf {
				if before, ok := beforeSelf[iface]; ok {
					// 如果该 iface 被识别为物理网卡（或你希望扣除的 iface），才从 total 扣除
					if isPhysicalNIC(iface) {
						// 注意：这些值是累计字节，可能与 net.IOCounters 的单位一致
						var sendDelta uint64
						var recvDelta uint64
						if after.BytesSent >= before.BytesSent {
							sendDelta = after.BytesSent - before.BytesSent
						} else {
							sendDelta = after.BytesSent
						}
						if after.BytesRecv >= before.BytesRecv {
							recvDelta = after.BytesRecv - before.BytesRecv
						} else {
							recvDelta = after.BytesRecv
						}
						selfSendDelta += sendDelta
						selfRecvDelta += recvDelta
					}
				}
			}
			// 将进程自身的增量按秒转换为速率并从 total 扣除
			if elapsed > 0 {
				totalBytesSentPerSec -= float64(selfSendDelta) / elapsed
				totalBytesRecvPerSec -= float64(selfRecvDelta) / elapsed
				// 防止负值
				if totalBytesSentPerSec < 0 {
					totalBytesSentPerSec = 0
				}
				if totalBytesRecvPerSec < 0 {
					totalBytesRecvPerSec = 0
				}
			}
		}
	}

	// 填充返回结构
	netMetrics.BytesSentPerSec = totalBytesSentPerSec
	netMetrics.BytesRecvPerSec = totalBytesRecvPerSec
	netMetrics.PacketsSentPerSec = totalPacketsSentPerSec
	netMetrics.PacketsRecvPerSec = totalPacketsRecvPerSec
	netMetrics.ErrorsPerSec = totalErrorsPerSec
	netMetrics.DropsPerSec = totalDropsPerSec
	netMetrics.TotalBytesSent = totalBytesSentAcc
	netMetrics.TotalBytesRecv = totalBytesRecvAcc
	netMetrics.TotalPacketsSent = totalPacketsSentAcc
	netMetrics.TotalPacketsRecv = totalPacketsRecvAcc
	netMetrics.Interfaces = interfaceMetrics

	// mark initialized and update prevTime
	c.initialized = true
	c.prevTime = now

	return netMetrics, nil
}
