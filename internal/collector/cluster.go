package collector

import (
	"context"

	"github.com/Y-vQv-Y/es-monitor/internal/client"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// ClusterCollector 集群指标采集器
type ClusterCollector struct {
	client *client.ElasticsearchClient
}

// NewClusterCollector 创建集群采集器
func NewClusterCollector(client *client.ElasticsearchClient) *ClusterCollector {
	return &ClusterCollector{
		client: client,
	}
}

// Collect 采集集群指标（只读操作）
func (c *ClusterCollector) Collect(ctx context.Context) (*model.ClusterHealth, error) {
	return c.client.GetClusterHealth(ctx)
}
