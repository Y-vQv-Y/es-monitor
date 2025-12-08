package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/Y-vQv-Y/es-monitor/internal/client"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// EnhancedCollector 增强监控采集器
type EnhancedCollector struct {
	client *client.ElasticsearchClient
}

// NewEnhancedCollector 创建增强采集器
func NewEnhancedCollector(client *client.ElasticsearchClient) *EnhancedCollector {
	return &EnhancedCollector{
		client: client,
	}
}

// Collect 采集增强指标
func (c *EnhancedCollector) Collect(ctx context.Context, nodeStats *model.NodeStats) (*model.EnhancedMetrics, error) {
	metrics := &model.EnhancedMetrics{
		HealthIssues: make([]model.HealthIssue, 0),
	}

	// 1. 分析线程池状态
	metrics.NodeThreadPools = c.analyzeThreadPools(nodeStats)
	
	// 2. 分析断路器状态
	metrics.NodeCircuitBreakers = c.analyzeCircuitBreakers(nodeStats)
	
	// 3. 检查集群健康问题
	c.checkHealthIssues(nodeStats, metrics)
	
	return metrics, nil
}

// analyzeThreadPools 分析线程池
func (c *EnhancedCollector) analyzeThreadPools(nodeStats *model.NodeStats) map[string]model.ThreadPoolStats {
	// 这里需要从 ES API 获取线程池数据
	// 简化实现，实际需要调用 /_nodes/stats/thread_pool
	pools := make(map[string]model.ThreadPoolStats)
	
	// 示例数据结构（需要添加到 NodeStat 中）
	// 关键线程池：search, write, get, bulk, management
	
	return pools
}

// analyzeCircuitBreakers 分析断路器
func (c *EnhancedCollector) analyzeCircuitBreakers(nodeStats *model.NodeStats) map[string]model.CircuitBreakerStats {
	breakers := make(map[string]model.CircuitBreakerStats)
	
	// 需要从 /_nodes/stats 中的 breakers 字段获取
	// 关键断路器：parent, fielddata, request, in_flight_requests
	
	return breakers
}

// checkHealthIssues 检查健康问题
func (c *EnhancedCollector) checkHealthIssues(nodeStats *model.NodeStats, metrics *model.EnhancedMetrics) {
	now := time.Now().Unix()
	
	for nodeID, node := range nodeStats.Nodes {
		// 1. JVM 堆内存问题
		if node.JVM.Mem.HeapUsedPercent >= 85 {
			metrics.HealthIssues = append(metrics.HealthIssues, model.HealthIssue{
				Level:      "critical",
				Component:  "jvm",
				NodeName:   node.Name,
				Message:    fmt.Sprintf("JVM 堆内存使用率过高: %d%%", node.JVM.Mem.HeapUsedPercent),
				Value:      node.JVM.Mem.HeapUsedPercent,
				Threshold:  85,
				Timestamp:  now,
				Suggestion: "增加堆内存或优化查询，检查是否有内存泄漏",
			})
		} else if node.JVM.Mem.HeapUsedPercent >= 75 {
			metrics.HealthIssues = append(metrics.HealthIssues, model.HealthIssue{
				Level:      "warning",
				Component:  "jvm",
				NodeName:   node.Name,
				Message:    fmt.Sprintf("JVM 堆内存使用率偏高: %d%%", node.JVM.Mem.HeapUsedPercent),
				Value:      node.JVM.Mem.HeapUsedPercent,
				Threshold:  75,
				Timestamp:  now,
				Suggestion: "关注内存使用趋势，考虑优化查询或增加堆内存",
			})
		}

		// 2. Full GC 频繁
		if node.JVM.GC.Collectors.Old.CollectionCount > 10 {
			metrics.HealthIssues = append(metrics.HealthIssues, model.HealthIssue{
				Level:      "warning",
				Component:  "jvm",
				NodeName:   node.Name,
				Message:    fmt.Sprintf("Full GC 次数过多: %d", node.JVM.GC.Collectors.Old.CollectionCount),
				Value:      node.JVM.GC.Collectors.Old.CollectionCount,
				Threshold:  10,
				Timestamp:  now,
				Suggestion: "检查堆内存配置，优化查询和聚合，减少 fielddata 使用",
			})
		}

		// 3. 磁盘空间不足
		diskTotal := float64(node.FS.Total.TotalInBytes) / 1024 / 1024 / 1024
		diskAvail := float64(node.FS.Total.AvailableInBytes) / 1024 / 1024 / 1024
		diskUsedPercent := (diskTotal - diskAvail) / diskTotal * 100

		if diskUsedPercent >= 90 {
			metrics.HealthIssues = append(metrics.HealthIssues, model.HealthIssue{
				Level:      "critical",
				Component:  "disk",
				NodeName:   node.Name,
				Message:    fmt.Sprintf("磁盘空间严重不足: %.1f%%", diskUsedPercent),
				Value:      diskUsedPercent,
				Threshold:  90.0,
				Timestamp:  now,
				Suggestion: "立即清理旧索引或扩展磁盘容量，启用 ILM 策略",
			})
		} else if diskUsedPercent >= 85 {
			metrics.HealthIssues = append(metrics.HealthIssues, model.HealthIssue{
				Level:      "warning",
				Component:  "disk",
				NodeName:   node.Name,
				Message:    fmt.Sprintf("磁盘空间不足: %.1f%%", diskUsedPercent),
				Value:      diskUsedPercent,
				Threshold:  85.0,
				Timestamp:  now,
				Suggestion: "计划清理旧索引或扩展磁盘，检查索引增长速度",
			})
		}

		// 4. 文件描述符使用率
		fdPercent := float64(node.Process.OpenFileDescriptors) / float64(node.Process.MaxFileDescriptors) * 100
		if fdPercent >= 80 {
			metrics.HealthIssues = append(metrics.HealthIssues, model.HealthIssue{
				Level:      "warning",
				Component:  "system",
				NodeName:   node.Name,
				Message:    fmt.Sprintf("文件描述符使用率过高: %.1f%%", fdPercent),
				Value:      fdPercent,
				Threshold:  80.0,
				Timestamp:  now,
				Suggestion: "增加系统文件描述符限制: ulimit -n",
			})
		}

		// 5. 段数量过多（影响性能）
		// 需要从 indices stats 获取
		// if segments > 1000 per shard {
		//     warning: too many segments, need force merge
		// }
	}
}
