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

// CollectStats 采集索引统计
func (c *IndexCollector) CollectStats(ctx context.Context) (*model.IndexStats, error) {
	stats, err := c.client.GetIndexStats(ctx)
	if err != nil {
		return nil, err
	}

	// 合并分段信息
	segments, err := c.client.GetIndexSegments(ctx)
	if err == nil {
		for name, seg := range segments.Indices {
			if s, ok := stats.Indices[name]; ok {
				s.Segments = seg.Segments
				stats.Indices[name] = s
			}
		}
	}

	return stats, nil
}

// CollectList 采集索引列表
func (c *IndexCollector) CollectList(ctx context.Context) ([]model.IndexInfo, error) {
	return c.client.GetCatIndices(ctx)
}
