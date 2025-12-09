package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/Y-vQv-Y/es-monitor/internal/client"
	"github.com/Y-vQv-Y/es-monitor/internal/collector"
	"github.com/Y-vQv-Y/es-monitor/internal/config"
	"github.com/Y-vQv-Y/es-monitor/internal/display"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// Monitor 监控器
type Monitor struct {
	client           *client.ElasticsearchClient
	config           *config.Config
	terminal         *display.Terminal
	clusterCollector *collector.ClusterCollector
	nodeCollector    *collector.NodeCollector
	indexCollector   *collector.IndexCollector
	systemCollector  *collector.SystemCollector
	prevNodeData     map[string]*display.PrevNodeMetrics
	prevIndexData    map[string]*display.PrevIndexMetrics
	ticker           *time.Ticker
	stopChan         chan struct{}
	wg               sync.WaitGroup
	mu               sync.Mutex
}

// NewMonitor 创建监控器
func NewMonitor(client *client.ElasticsearchClient, cfg *config.Config) *Monitor {
	return &Monitor{
		client:           client,
		config:           cfg,
		terminal:         display.NewTerminal(),
		clusterCollector: collector.NewClusterCollector(client),
		nodeCollector:    collector.NewNodeCollector(client),
		indexCollector:   collector.NewIndexCollector(client),
		systemCollector:  collector.NewSystemCollector(),
		prevNodeData:     make(map[string]*display.PrevNodeMetrics),
		prevIndexData:    make(map[string]*display.PrevIndexMetrics),
		stopChan:         make(chan struct{}),
	}
}

// Start 启动监控（生产环境安全）
func (m *Monitor) Start(ctx context.Context) {
	m.ticker = time.NewTicker(m.config.Interval)
	defer m.ticker.Stop()

	// 首次立即执行
	m.collect(ctx)

	// 定期采集
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		case <-m.ticker.C:
			m.collect(ctx)
		}
	}
}

// Stop 安全停止监控
func (m *Monitor) Stop() {
	close(m.stopChan)
	m.wg.Wait()
}

// collect 采集并显示所有指标（只读操作）
func (m *Monitor) collect(ctx context.Context) {
	m.wg.Add(1)
	defer m.wg.Done()

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.systemCollector != nil && m.client.netTracker != nil {
	   m.systemCollector.SetNetworkTracker(m.client.netTracker)
	}

	// 显示标题
	m.terminal.DisplayHeader()

	// 1. 采集集群健康状态
	health, err := m.clusterCollector.Collect(ctx)
	if err != nil {
		m.terminal.DisplayError("获取集群健康状态失败", err)
		m.terminal.DisplayFooter()
		return
	}
	m.terminal.DisplayClusterHealth(health)

	// 2. 采集系统指标（优先显示，最关心的指标）
	sysMetrics, err := m.systemCollector.Collect(ctx)
	if err != nil {
		m.terminal.DisplayError("获取系统指标失败", err)
	} else {
		// 显示完整的系统资源监控
		m.terminal.DisplaySystemMetrics(sysMetrics)
		m.terminal.DisplayDiskMetrics(&sysMetrics.Disk)
		m.terminal.DisplayNetworkMetrics(&sysMetrics.Network)
	}

	// 3. 采集节点统计
	nodeStats, err := m.nodeCollector.Collect(ctx)
	if err != nil {
		m.terminal.DisplayError("获取节点统计失败", err)
	} else {
		m.terminal.DisplayNodeStats(nodeStats, m.prevNodeData)
		m.updateNodePrevData(nodeStats)
	}

	// 4. 采集索引统计
	indexList, err := m.indexCollector.CollectList(ctx)
	if err != nil {
		m.terminal.DisplayError("获取索引列表失败", err)
	} else {
		indexStats, err := m.indexCollector.CollectStats(ctx)
		if err != nil {
			m.terminal.DisplayError("获取索引统计失败", err)
		} else {
			m.terminal.DisplayIndexStats(indexList, indexStats, m.prevIndexData)
			m.updateIndexPrevData(indexStats)
		}
	}

	// 显示页脚
	m.terminal.DisplayFooter()
}

// updateNodePrevData 更新节点历史数据
func (m *Monitor) updateNodePrevData(stats *model.NodeStats) {
	now := time.Now()
	for nodeID, node := range stats.Nodes {
		m.prevNodeData[nodeID] = &display.PrevNodeMetrics{
			IndexTotal: node.Indices.Indexing.IndexTotal,
			QueryTotal: node.Indices.Search.QueryTotal,
			Timestamp:  now,
		}
	}
}

// updateIndexPrevData 更新索引历史数据
func (m *Monitor) updateIndexPrevData(stats *model.IndexStats) {
	now := time.Now()
	for indexName, indexStat := range stats.Indices {
		m.prevIndexData[indexName] = &display.PrevIndexMetrics{
			IndexTotal: indexStat.Total.Indexing.IndexTotal,
			QueryTotal: indexStat.Total.Search.QueryTotal,
			Timestamp:  now,
		}
	}
}
