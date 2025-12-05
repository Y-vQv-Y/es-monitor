package model

// ClusterHealth 集群健康状态
type ClusterHealth struct {
	ClusterName         string  `json:"cluster_name"`
	Status              string  `json:"status"`
	TimedOut            bool    `json:"timed_out"`
	NumberOfNodes       int     `json:"number_of_nodes"`
	NumberOfDataNodes   int     `json:"number_of_data_nodes"`
	ActivePrimaryShards int     `json:"active_primary_shards"`
	ActiveShards        int     `json:"active_shards"`
	RelocatingShards    int     `json:"relocating_shards"`
	InitializingShards  int     `json:"initializing_shards"`
	UnassignedShards    int     `json:"unassigned_shards"`
	DelayedUnassigned   int     `json:"delayed_unassigned_shards"`
	PendingTasks        int     `json:"number_of_pending_tasks"`
	InFlightFetch       int     `json:"number_of_in_flight_fetch"`
	ActiveShardsPercent float64 `json:"active_shards_percent_as_number"`
}
