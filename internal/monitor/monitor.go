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
	
	// 缓存最新采集的数据
	latestSystemMetrics *model.SystemMetrics
	latestHealth        *model.ClusterHealth
	latestNodeStats     *model.NodeStats
	latestIndexList     []*model.IndexInfo
	latestIndexStats    *model.IndexStats
	dataMu              sync.RWMutex
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

	// 首次立即执行系统采集
	m.collectSystemMetrics(ctx)
	
	// 延迟 2 秒后采集 ES 指标
	time.Sleep(2 * time.Second)
	m.collectESMetrics(ctx)
	
	// 首次显示
	m.display(ctx)

	// 启动两个独立的采集循环
	go m.systemCollectionLoop(ctx)
	go m.esCollectionLoop(ctx)

	// 主循环负责定期显示
	displayTicker := time.NewTicker(m.config.Interval)
	defer displayTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		case <-displayTicker.C:
			m.display(ctx)
		}
	}
}

// Stop 安全停止监控
func (m *Monitor) Stop() {
	close(m.stopChan)
	m.wg.Wait()
}

// systemCollectionLoop 系统指标采集循环
func (m *Monitor) systemCollectionLoop(ctx context.Context) {
	ticker := time.NewTicker(m.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.collectSystemMetrics(ctx)
		}
	}
}

// esCollectionLoop ES 指标采集循环（错开 2 秒）
func (m *Monitor) esCollectionLoop(ctx context.Context) {
	ticker := time.NewTicker(m.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		case <-ticker.C:
			// 在系统采集后 2 秒再采集 ES 指标
			time.Sleep(2 * time.Second)
			m.collectESMetrics(ctx)
		}
	}
}

// collectSystemMetrics 采集系统指标（无网络请求）
func (m *Monitor) collectSystemMetrics(ctx context.Context) {
	sysMetrics, err := m.systemCollector.Collect(ctx)
	if err != nil {
		return
	}

	m.dataMu.Lock()
	m.latestSystemMetrics = sysMetrics
	m.dataMu.Unlock()
}

// collectESMetrics 采集 ES 指标（有网络请求）
func (m *Monitor) collectESMetrics(ctx context.Context) {
	// 并发采集所有 ES 指标
	var wg sync.WaitGroup
	wg.Add(4)

	var health *model.ClusterHealth
	var nodeStats *model.NodeStats
	var indexList []*model.IndexInfo
	var indexStats *model.IndexStats

	go func() {
		defer wg.Done()
		health, _ = m.clusterCollector.Collect(ctx)
	}()

	go func() {
		defer wg.Done()
		nodeStats, _ = m.nodeCollector.Collect(ctx)
	}()

	go func() {
		defer wg.Done()
		indexList, _ = m.indexCollector.CollectList(ctx)
	}()

	go func() {
		defer wg.Done()
		indexStats, _ = m.indexCollector.CollectStats(ctx)
	}()

	wg.Wait()

	// 更新缓存
	m.dataMu.Lock()
	m.latestHealth = health
	m.latestNodeStats = nodeStats
	m.latestIndexList = indexList
	m.latestIndexStats = indexStats
	m.dataMu.Unlock()
}

// display 显示所有指标（使用缓存的数据）
func (m *Monitor) display(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.dataMu.RLock()
	defer m.dataMu.RUnlock()

	// 显示标题
	m.terminal.DisplayHeader()

	// 1. 显示集群健康状态
	if m.latestHealth != nil {
		m.terminal.DisplayClusterHealth(m.latestHealth)
	}

	// 2. 显示系统指标
	if m.latestSystemMetrics != nil {
		m.terminal.DisplaySystemMetrics(m.latestSystemMetrics)
		m.terminal.DisplayDiskMetrics(&m.latestSystemMetrics.Disk)
		m.terminal.DisplayNetworkMetrics(&m.latestSystemMetrics.Network)
	}

	// 3. 显示节点统计
	if m.latestNodeStats != nil {
		m.terminal.DisplayNodeStats(m.latestNodeStats, m.prevNodeData)
		m.updateNodePrevData(m.latestNodeStats)
	}

	// 4. 显示索引统计
	if m.latestIndexList != nil && m.latestIndexStats != nil {
		m.terminal.DisplayIndexStats(m.latestIndexList, m.latestIndexStats, m.prevIndexData)
		m.updateIndexPrevData(m.latestIndexStats)
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
