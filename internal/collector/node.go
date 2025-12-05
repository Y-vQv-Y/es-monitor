package collector

import (
	"context"

	"github.com/Y-vQv-Y/es-monitor/internal/client"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// NodeCollector 节点指标采集器
type NodeCollector struct {
	client *client.ElasticsearchClient
}

// NewNodeCollector 创建节点采集器
func NewNodeCollector(client *client.ElasticsearchClient) *NodeCollector {
	return &NodeCollector{
		client: client,
	}
}

// Collect 采集节点指标（只读操作）
func (c *NodeCollector) Collect(ctx context.Context) (*model.NodeStats, error) {
	return c.client.GetNodeStats(ctx)
}
