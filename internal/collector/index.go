package collector

import (
	"context"

	"github.com/Y-vQv-Y/es-monitor/internal/client"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// IndexCollector 索引指标采集器
type IndexCollector struct {
	client *client.ElasticsearchClient
}

// NewIndexCollector 创建索引采集器
func NewIndexCollector(client *client.ElasticsearchClient) *IndexCollector {
	return &IndexCollector{
		client: client,
	}
}

// CollectStats 采集索引统计（只读操作）
func (c *IndexCollector) CollectStats(ctx context.Context) (*model.IndexStats, error) {
	return c.client.GetIndexStats(ctx)
}

// CollectList 采集索引列表（只读操作）
func (c *IndexCollector) CollectList(ctx context.Context) ([]model.IndexInfo, error) {
	return c.client.GetCatIndices(ctx)
}
