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

// ElasticsearchClient ES 客户端
type ElasticsearchClient struct {
	baseURL string
	client  *http.Client
	config  *config.Config
}

// NewElasticsearchClient 创建 ES 客户端
func NewElasticsearchClient(cfg *config.Config) *ElasticsearchClient {
	baseURL := fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)
	return &ElasticsearchClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		config: cfg,
	}
}

// request 发送 HTTP 请求
func (c *ElasticsearchClient) request(ctx context.Context, endpoint string) ([]byte, error) {
	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加认证
	if c.config.Username != "" && c.config.Password != "" {
		req.SetBasicAuth(c.config.Username, c.config.Password)
	}

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

// Ping 测试连接
func (c *ElasticsearchClient) Ping(ctx context.Context) error {
	_, err := c.request(ctx, "/")
	return err
}

// GetClusterHealth 获取集群健康状态
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

// GetNodeStats 获取节点统计
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

// GetIndexStats 获取索引统计
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

// GetIndexSegments 获取索引分段信息（新增）
func (c *ElasticsearchClient) GetIndexSegments(ctx context.Context) (*model.IndexStats, error) {
	data, err := c.request(ctx, "/_segments")
	if err != nil {
		return nil, err
	}

	var stats model.IndexStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return &stats, nil
}

// GetCatIndices 获取索引列表
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
