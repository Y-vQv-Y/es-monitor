package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Y-vQv-Y/es-monitor/internal/config"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// ElasticsearchClient 添加流量计数器
type ElasticsearchClient struct {
	baseURL      string
	client       *http.Client
	config       *config.Config
	safety       *config.SafetyConfig
	netTracker   *NetworkTracker // 新增
}

// NetworkTracker 网络流量跟踪器
type NetworkTracker struct {
	BytesSent uint64
	BytesRecv uint64
	mu        sync.Mutex
}

// trackingTransport 包装 http.Transport 以跟踪流量
type trackingTransport struct {
	base    http.RoundTripper
	tracker *NetworkTracker
}

func (t *trackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 记录发送的字节数（估算）
	if req.Body != nil {
		if req.ContentLength > 0 {
			t.tracker.addBytesSent(uint64(req.ContentLength))
		}
	}
	// 请求头大小（粗略估算）
	headerSize := uint64(len(req.Method) + len(req.URL.String()) + 100)
	t.tracker.addBytesSent(headerSize)

	// 执行请求
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// 记录接收的字节数
	if resp.ContentLength > 0 {
		t.tracker.addBytesRecv(uint64(resp.ContentLength))
	} else {
		// 如果没有 ContentLength，读取 body 并计数
		body, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		t.tracker.addBytesRecv(uint64(len(body)))
	}
	// 响应头大小（粗略估算）
	t.tracker.addBytesRecv(200)

	return resp, nil
}


// NewElasticsearchClient 创建 ES 客户端
func NewElasticsearchClient(cfg *config.Config) *ElasticsearchClient {
	client := &ElasticsearchClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: config.DefaultSafetyConfig.RequestTimeout,
			Transport: &trackingTransport{
				base:    http.DefaultTransport,
				tracker: &NetworkTracker{}, // 新增追踪器
			},
		},
		config:     cfg,
		safety:     &config.DefaultSafetyConfig,
		netTracker: &NetworkTracker{}, // 新增
	}
	
	return client
}

func (t *trackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 记录发送的字节数（估算）
	if req.Body != nil {
		if req.ContentLength > 0 {
			t.tracker.addBytesSent(uint64(req.ContentLength))
		}
	}
	// 请求头大小（粗略估算）
	headerSize := uint64(len(req.Method) + len(req.URL.String()) + 100)
	t.tracker.addBytesSent(headerSize)

	// 执行请求
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// 记录接收的字节数
	if resp.ContentLength > 0 {
		t.tracker.addBytesRecv(uint64(resp.ContentLength))
	} else {
		// 如果没有 ContentLength，读取 body 并计数
		body, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		t.tracker.addBytesRecv(uint64(len(body)))
	}
	// 响应头大小（粗略估算）
	t.tracker.addBytesRecv(200)

	return resp, nil
}

// addBytesSent 添加发送字节数
func (n *NetworkTracker) addBytesSent(bytes uint64) {
	n.mu.Lock()
	n.BytesSent += bytes
	n.mu.Unlock()
}

// addBytesRecv 添加接收字节数
func (n *NetworkTracker) addBytesRecv(bytes uint64) {
	n.mu.Lock()
	n.BytesRecv += bytes
	n.mu.Unlock()
}

// GetAndReset 获取并重置计数器
func (n *NetworkTracker) GetAndReset() (sent, recv uint64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	sent, recv = n.BytesSent, n.BytesRecv
	n.BytesSent, n.BytesRecv = 0, 0
	return
}

// isEndpointAllowed 检查端点是否允许访问（生产环境安全检查）
func (c *ElasticsearchClient) isEndpointAllowed(endpoint string) bool {
	if !c.config.ReadOnly {
		return false
	}
	
	for _, allowed := range c.safety.AllowedEndpoints {
		if endpoint == allowed || len(endpoint) > len(allowed) && endpoint[:len(allowed)] == allowed {
			return true
		}
	}
	return false
}

// request 发送 HTTP 请求（只读安全版本）
func (c *ElasticsearchClient) request(ctx context.Context, endpoint string) ([]byte, error) {
	// 生产环境安全检查
	if !c.isEndpointAllowed(endpoint) {
		return nil, fmt.Errorf("生产环境安全限制: 不允许访问端点 %s", endpoint)
	}
	
	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加认证
	if c.config.Username != "" && c.config.Password != "" {
		req.SetBasicAuth(c.config.Username, c.config.Password)
	}
	
	// 添加安全请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ES-Monitor/1.0 (ReadOnly)")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP 状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}

// Ping 测试连接（安全操作）
func (c *ElasticsearchClient) Ping(ctx context.Context) error {
	_, err := c.request(ctx, "/")
	return err
}

// GetClusterHealth 获取集群健康状态（只读操作）
func (c *ElasticsearchClient) GetClusterHealth(ctx context.Context) (*model.ClusterHealth, error) {
	data, err := c.request(ctx, "/_cluster/health")
	if err != nil {
		return nil, err
	}

	var health model.ClusterHealth
	if err := json.Unmarshal(data, &health); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return &health, nil
}

// GetNodeStats 获取节点统计（只读操作）
func (c *ElasticsearchClient) GetNodeStats(ctx context.Context) (*model.NodeStats, error) {
	data, err := c.request(ctx, "/_nodes/stats")
	if err != nil {
		return nil, err
	}

	var stats model.NodeStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return &stats, nil
}

// GetIndexStats 获取索引统计（只读操作）
func (c *ElasticsearchClient) GetIndexStats(ctx context.Context) (*model.IndexStats, error) {
	data, err := c.request(ctx, "/_stats")
	if err != nil {
		return nil, err
	}

	var stats model.IndexStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return &stats, nil
}

// GetCatIndices 获取索引列表（只读操作）
func (c *ElasticsearchClient) GetCatIndices(ctx context.Context) ([]model.IndexInfo, error) {
	data, err := c.request(ctx, "/_cat/indices?format=json&bytes=b")
	if err != nil {
		return nil, err
	}

	var indices []model.IndexInfo
	if err := json.Unmarshal(data, &indices); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return indices, nil
}
