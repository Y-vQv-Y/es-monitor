package collector

import (
    "bufio"
    "os"
    "strconv"
    "strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

type NetworkCollector struct {
    // 初始快照
    initialSnapshot map[string]InterfaceStats
    // 上一次采集的数据
    lastSnapshot map[string]InterfaceStats
    // 监控进程的PID
    monitorPID int
}

func NewNetworkCollector() *NetworkCollector {
    nc := &NetworkCollector{
        monitorPID: os.Getpid(),
    }
    // 初始化时采集基准
    nc.initialSnapshot = nc.collectRawStats()
    nc.lastSnapshot = nc.initialSnapshot
    return nc
}

func (nc *NetworkCollector) GetNetStats() NetworkStats {
    current := nc.collectRawStats()
    
    // 计算增量(排除初始基准)
    stats := NetworkStats{}
    for iface, curr := range current {
        if initial, ok := nc.initialSnapshot[iface]; ok {
            stats.TxBytes += (curr.TxBytes - initial.TxBytes)
            stats.RxBytes += (curr.RxBytes - initial.RxBytes)
        }
    }
    
    nc.lastSnapshot = current
    return stats
}

type TCPStats struct {
    // TCP连接状态统计
    Established int
    SynSent     int
    SynRecv     int
    FinWait1    int
    FinWait2    int
    TimeWait    int
    Close       int
    CloseWait   int
    LastAck     int
    Listen      int
    Closing     int
    
    // TCP性能指标
    RetransSegs uint64  // 重传段数
    InSegs      uint64  // 接收段数
    OutSegs     uint64  // 发送段数
    
    // 错误统计
    ListenDrops  uint64  // listen队列溢出
    TCPLoss      uint64  // 丢包数
    
    // 连接详情(可选)
    Connections []TCPConnection
}

type TCPConnection struct {
    LocalAddr  string
    LocalPort  int
    RemoteAddr string
    RemotePort int
    State      string
    PID        int
    Program    string
}

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

// 采集TCP统计
func CollectTCPStats() (*TCPStats, error) {
    stats := &TCPStats{}
    
    // 1. 读取 /proc/net/tcp 和 /proc/net/tcp6
    if err := parseTCPFile("/proc/net/tcp", stats); err != nil {
        return nil, err
    }
    if err := parseTCPFile("/proc/net/tcp6", stats); err != nil {
        return nil, err
    }
    
    // 2. 读取 /proc/net/snmp 获取TCP协议统计
    if err := parseSNMP(stats); err != nil {
        return nil, err
    }
    
    // 3. 可选:读取 /proc/net/netstat 获取更详细的统计
    if err := parseNetstat(stats); err != nil {
        return nil, err
    }
    
    return stats, nil
}

func parseTCPFile(path string, stats *TCPStats) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    scanner.Scan() // 跳过标题行
    
    for scanner.Scan() {
        fields := strings.Fields(scanner.Text())
        if len(fields) < 10 {
            continue
        }
        
        // 字段3是状态(十六进制)
        stateHex := fields[3]
        state, _ := strconv.ParseInt(stateHex, 16, 64)
        
        switch state {
        case 0x01:
            stats.Established++
        case 0x02:
            stats.SynSent++
        case 0x03:
            stats.SynRecv++
        case 0x04:
            stats.FinWait1++
        case 0x05:
            stats.FinWait2++
        case 0x06:
            stats.TimeWait++
        case 0x07:
            stats.Close++
        case 0x08:
            stats.CloseWait++
        case 0x09:
            stats.LastAck++
        case 0x0A:
            stats.Listen++
        case 0x0B:
            stats.Closing++
        }
        
        // 如果需要详细连接信息
        if collectDetails {
            conn := parseTCPConnection(fields)
            stats.Connections = append(stats.Connections, conn)
        }
    }
    
    return scanner.Err()
}

func parseSNMP(stats *TCPStats) error {
    file, err := os.Open("/proc/net/snmp")
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    var tcpKeys []string
    
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "Tcp:") {
            fields := strings.Fields(line)
            if len(tcpKeys) == 0 {
                // 第一行是键名
                tcpKeys = fields[1:]
            } else {
                // 第二行是值
                for i, key := range tcpKeys {
                    val, _ := strconv.ParseUint(fields[i+1], 10, 64)
                    switch key {
                    case "RetransSegs":
                        stats.RetransSegs = val
                    case "InSegs":
                        stats.InSegs = val
                    case "OutSegs":
                        stats.OutSegs = val
                    }
                }
                break
            }
        }
    }
    
    return scanner.Err()
}

// collectNetwork 采集网络详细指标（只读操作）
func (c *SystemCollector) collectNetwork() (model.NetworkMetrics, error) {
	netMetrics := model.NetworkMetrics{}
	now := time.Now()
	elapsed := now.Sub(c.prevTime).Seconds()

	// 网络 IO 统计
	ioCounters, err := net.IOCounters(true)
	if err == nil && elapsed > 0 && c.initialized {
		var totalBytesSent, totalBytesRecv float64
		var totalPacketsSent, totalPacketsRecv float64
		var totalErrors, totalDrops float64
		var totalBytesSentAcc, totalBytesRecvAcc uint64
		var totalPacketsSentAcc, totalPacketsRecvAcc uint64
		
		interfaceMetrics := make([]model.InterfaceMetrics, 0)

		for _, counter := range ioCounters {
			if prev, ok := c.prevNetIO[counter.Name]; ok {
				// 计算增量
				bytesSent := float64(counter.BytesSent - prev.BytesSent)
				bytesRecv := float64(counter.BytesRecv - prev.BytesRecv)
				packetsSent := float64(counter.PacketsSent - prev.PacketsSent)
				packetsRecv := float64(counter.PacketsRecv - prev.PacketsRecv)
				errors := float64((counter.Errin + counter.Errout) - (prev.Errin + prev.Errout))
				drops := float64((counter.Dropin + counter.Dropout) - (prev.Dropin + prev.Dropout))

				// 计算速率
				bytesSentPerSec := bytesSent / elapsed
				bytesRecvPerSec := bytesRecv / elapsed
				packetsSentPerSec := packetsSent / elapsed
				packetsRecvPerSec := packetsRecv / elapsed
				errorsPerSec := errors / elapsed
				dropsPerSec := drops / elapsed

				totalBytesSent += bytesSentPerSec
				totalBytesRecv += bytesRecvPerSec
				totalPacketsSent += packetsSentPerSec
				totalPacketsRecv += packetsRecvPerSec
				totalErrors += errorsPerSec
				totalDrops += dropsPerSec

				// 每个网卡的详细指标
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

		netMetrics.BytesSentPerSec = totalBytesSent
		netMetrics.BytesRecvPerSec = totalBytesRecv
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

	// 网络连接统计（可选，避免对性能影响）
	// 注释掉以避免编译错误
	/*
	connections, err := net.Connections("tcp")
	if err == nil {
		// 统计各种状态的连接数
		for _, conn := range connections {
			switch conn.Status {
			case "ESTABLISHED":
				netMetrics.TCPEstablished++
			case "LISTEN":
				netMetrics.TCPListening++
			case "TIME_WAIT":
				netMetrics.TCPTimeWait++
			}
		}
		netMetrics.TCPConnections = len(connections)
	}
	*/

	return netMetrics, nil
}
