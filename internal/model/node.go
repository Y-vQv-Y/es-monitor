// internal/model/node.go
package model

// NodeStats 节点统计信息（完整版）
type NodeStats struct {
	ClusterName string              `json:"cluster_name"`
	Nodes       map[string]NodeStat `json:"nodes"`
}

// NodeStat 单个节点统计
type NodeStat struct {
	Timestamp         int64     `json:"timestamp"`
	Name              string    `json:"name"`
	TransportAddress  string    `json:"transport_address"`
	Host              string    `json:"host"`
	IP                string    `json:"ip"`
	Version           string    `json:"version"`
	BuildFlavor       string    `json:"build_flavor"`
	BuildType         string    `json:"build_type"`
	BuildHash         string    `json:"build_hash"`
	Roles             []string  `json:"roles"`
	Attributes        map[string]string `json:"attributes"`
	Indices           IndicesStats      `json:"indices"`
	OS                OSStats           `json:"os"`
	Process           ProcessStats      `json:"process"`
	JVM               JVMStats          `json:"jvm"`
	ThreadPool        map[string]ThreadPoolStats `json:"thread_pool"`
	FS                FSStats           `json:"fs"`
	Transport         TransportStats    `json:"transport"`
	HTTP              HTTPStats         `json:"http"`
	Breakers          map[string]BreakerStats `json:"breakers"`
	Script            ScriptStats       `json:"script"`
	Discovery         DiscoveryStats    `json:"discovery"`
	Ingest            IngestStats       `json:"ingest"`
}

// IndicesStats 索引统计
type IndicesStats struct {
	Docs struct {
		Count   int64 `json:"count"`
		Deleted int64 `json:"deleted"`
	} `json:"docs"`
	Store struct {
		SizeInBytes          int64 `json:"size_in_bytes"`
		ReservedInBytes      int64 `json:"reserved_in_bytes"`
		TotalDataSetSizeInBytes int64 `json:"total_data_set_size_in_bytes"`
	} `json:"store"`
	Indexing struct {
		IndexTotal           int64 `json:"index_total"`
		IndexTimeInMillis    int64 `json:"index_time_in_millis"`
		IndexCurrent         int64 `json:"index_current"`
		IndexFailed          int64 `json:"index_failed"`
		DeleteTotal          int64 `json:"delete_total"`
		DeleteTimeInMillis   int64 `json:"delete_time_in_millis"`
		DeleteCurrent        int64 `json:"delete_current"`
		NoopUpdateTotal      int64 `json:"noop_update_total"`
		IsThrottled          bool  `json:"is_throttled"`
		ThrottleTimeInMillis int64 `json:"throttle_time_in_millis"`
	} `json:"indexing"`
	Get struct {
		Total               int64 `json:"total"`
		TimeInMillis        int64 `json:"time_in_millis"`
		ExistsTotal         int64 `json:"exists_total"`
		ExistsTimeInMillis  int64 `json:"exists_time_in_millis"`
		MissingTotal        int64 `json:"missing_total"`
		MissingTimeInMillis int64 `json:"missing_time_in_millis"`
		Current             int64 `json:"current"`
	} `json:"get"`
	Search struct {
		OpenContexts        int64 `json:"open_contexts"`
		QueryTotal          int64 `json:"query_total"`
		QueryTimeInMillis   int64 `json:"query_time_in_millis"`
		QueryCurrent        int64 `json:"query_current"`
		FetchTotal          int64 `json:"fetch_total"`
		FetchTimeInMillis   int64 `json:"fetch_time_in_millis"`
		FetchCurrent        int64 `json:"fetch_current"`
		ScrollTotal         int64 `json:"scroll_total"`
		ScrollTimeInMillis  int64 `json:"scroll_time_in_millis"`
		ScrollCurrent       int64 `json:"scroll_current"`
		SuggestTotal        int64 `json:"suggest_total"`
		SuggestTimeInMillis int64 `json:"suggest_time_in_millis"`
		SuggestCurrent      int64 `json:"suggest_current"`
	} `json:"search"`
	Merges struct {
		Current                    int64 `json:"current"`
		CurrentDocs                int64 `json:"current_docs"`
		CurrentSizeInBytes         int64 `json:"current_size_in_bytes"`
		Total                      int64 `json:"total"`
		TotalTimeInMillis          int64 `json:"total_time_in_millis"`
		TotalDocs                  int64 `json:"total_docs"`
		TotalSizeInBytes           int64 `json:"total_size_in_bytes"`
		TotalStoppedTimeInMillis   int64 `json:"total_stopped_time_in_millis"`
		TotalThrottledTimeInMillis int64 `json:"total_throttled_time_in_millis"`
		TotalAutoThrottleInBytes   int64 `json:"total_auto_throttle_in_bytes"`
	} `json:"merges"`
	Refresh struct {
		Total             int64 `json:"total"`
		TotalTimeInMillis int64 `json:"total_time_in_millis"`
		ExternalTotal     int64 `json:"external_total"`
		ExternalTotalTimeInMillis int64 `json:"external_total_time_in_millis"`
		Listeners         int64 `json:"listeners"`
	} `json:"refresh"`
	Flush struct {
		Total             int64 `json:"total"`
		Periodic          int64 `json:"periodic"`
		TotalTimeInMillis int64 `json:"total_time_in_millis"`
	} `json:"flush"`
	Warmer struct {
		Current           int64 `json:"current"`
		Total             int64 `json:"total"`
		TotalTimeInMillis int64 `json:"total_time_in_millis"`
	} `json:"warmer"`
	QueryCache struct {
		MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
		TotalCount        int64 `json:"total_count"`
		HitCount          int64 `json:"hit_count"`
		MissCount         int64 `json:"miss_count"`
		CacheSize         int64 `json:"cache_size"`
		CacheCount        int64 `json:"cache_count"`
		Evictions         int64 `json:"evictions"`
	} `json:"query_cache"`
	Fielddata struct {
		MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
		Evictions         int64 `json:"evictions"`
	} `json:"fielddata"`
	Completion struct {
		SizeInBytes int64 `json:"size_in_bytes"`
	} `json:"completion"`
	Segments struct {
		Count                     int64 `json:"count"`
		MemoryInBytes             int64 `json:"memory_in_bytes"`
		TermsMemoryInBytes        int64 `json:"terms_memory_in_bytes"`
		StoredFieldsMemoryInBytes int64 `json:"stored_fields_memory_in_bytes"`
		TermVectorsMemoryInBytes  int64 `json:"term_vectors_memory_in_bytes"`
		NormsMemoryInBytes        int64 `json:"norms_memory_in_bytes"`
		PointsMemoryInBytes       int64 `json:"points_memory_in_bytes"`
		DocValuesMemoryInBytes    int64 `json:"doc_values_memory_in_bytes"`
		IndexWriterMemoryInBytes  int64 `json:"index_writer_memory_in_bytes"`
		VersionMapMemoryInBytes   int64 `json:"version_map_memory_in_bytes"`
		FixedBitSetMemoryInBytes  int64 `json:"fixed_bit_set_memory_in_bytes"`
		MaxUnsafeAutoIdTimestamp  int64 `json:"max_unsafe_auto_id_timestamp"`
	} `json:"segments"`
	Translog struct {
		Operations              int64 `json:"operations"`
		SizeInBytes             int64 `json:"size_in_bytes"`
		UncommittedOperations   int64 `json:"uncommitted_operations"`
		UncommittedSizeInBytes  int64 `json:"uncommitted_size_in_bytes"`
		EarliestLastModifiedAge int64 `json:"earliest_last_modified_age"`
	} `json:"translog"`
	RequestCache struct {
		MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
		Evictions         int64 `json:"evictions"`
		HitCount          int64 `json:"hit_count"`
		MissCount         int64 `json:"miss_count"`
	} `json:"request_cache"`
	Recovery struct {
		CurrentAsSource      int64 `json:"current_as_source"`
		CurrentAsTarget      int64 `json:"current_as_target"`
		ThrottleTimeInMillis int64 `json:"throttle_time_in_millis"`
	} `json:"recovery"`
}

// OSStats 操作系统统计
type OSStats struct {
	Timestamp int64 `json:"timestamp"`
	CPU       struct {
		Percent     int `json:"percent"`
		LoadAverage struct {
			OneM     float64 `json:"1m"`
			FiveM    float64 `json:"5m"`
			FifteenM float64 `json:"15m"`
		} `json:"load_average"`
	} `json:"cpu"`
	Mem struct {
		TotalInBytes     int64 `json:"total_in_bytes"`
		FreeInBytes      int64 `json:"free_in_bytes"`
		UsedInBytes      int64 `json:"used_in_bytes"`
		FreePercent      int   `json:"free_percent"`
		UsedPercent      int   `json:"used_percent"`
	} `json:"mem"`
	Swap struct {
		TotalInBytes int64 `json:"total_in_bytes"`
		FreeInBytes  int64 `json:"free_in_bytes"`
		UsedInBytes  int64 `json:"used_in_bytes"`
	} `json:"swap"`
	Cgroup struct {
		CPUAcct struct {
			ControlGroup       string `json:"control_group"`
			UsageNanos         int64  `json:"usage_nanos"`
		} `json:"cpuacct"`
		CPU struct {
			ControlGroup          string `json:"control_group"`
			CfsPeriodMicros       int64  `json:"cfs_period_micros"`
			CfsQuotaMicros        int64  `json:"cfs_quota_micros"`
			NumberOfElapsedPeriods int64  `json:"number_of_elapsed_periods"`
			NumberOfTimesThrottled int64  `json:"number_of_times_throttled"`
			TimeThrottledNanos     int64  `json:"time_throttled_nanos"`
		} `json:"cpu"`
		Memory struct {
			ControlGroup    string `json:"control_group"`
			LimitInBytes    string `json:"limit_in_bytes"`
			UsageInBytes    int64  `json:"usage_in_bytes"`
		} `json:"memory"`
	} `json:"cgroup"`
}

// ProcessStats 进程统计
type ProcessStats struct {
	Timestamp           int64 `json:"timestamp"`
	OpenFileDescriptors int64 `json:"open_file_descriptors"`
	MaxFileDescriptors  int64 `json:"max_file_descriptors"`
	CPU                 struct {
		Percent       int   `json:"percent"`
		TotalInMillis int64 `json:"total_in_millis"`
	} `json:"cpu"`
	Mem struct {
		TotalVirtualInBytes int64 `json:"total_virtual_in_bytes"`
	} `json:"mem"`
}

// JVMStats JVM 统计（完整版）
type JVMStats struct {
	Timestamp      int64 `json:"timestamp"`
	UptimeInMillis int64 `json:"uptime_in_millis"`
	Mem            struct {
		HeapUsedInBytes         int64 `json:"heap_used_in_bytes"`
		HeapUsedPercent         int   `json:"heap_used_percent"`
		HeapCommittedInBytes    int64 `json:"heap_committed_in_bytes"`
		HeapMaxInBytes          int64 `json:"heap_max_in_bytes"`
		NonHeapUsedInBytes      int64 `json:"non_heap_used_in_bytes"`
		NonHeapCommittedInBytes int64 `json:"non_heap_committed_in_bytes"`
		Pools                   map[string]struct {
			UsedInBytes      int64 `json:"used_in_bytes"`
			MaxInBytes       int64 `json:"max_in_bytes"`
			PeakUsedInBytes  int64 `json:"peak_used_in_bytes"`
			PeakMaxInBytes   int64 `json:"peak_max_in_bytes"`
		} `json:"pools"`
	} `json:"mem"`
	Threads struct {
		Count     int `json:"count"`
		PeakCount int `json:"peak_count"`
	} `json:"threads"`
	GC struct {
		Collectors map[string]struct {
			CollectionCount        int64 `json:"collection_count"`
			CollectionTimeInMillis int64 `json:"collection_time_in_millis"`
		} `json:"collectors"`
	} `json:"gc"`
	BufferPools map[string]struct {
		Count                int64 `json:"count"`
		UsedInBytes          int64 `json:"used_in_bytes"`
		TotalCapacityInBytes int64 `json:"total_capacity_in_bytes"`
	} `json:"buffer_pools"`
	Classes struct {
		CurrentLoadedCount int64 `json:"current_loaded_count"`
		TotalLoadedCount   int64 `json:"total_loaded_count"`
		TotalUnloadedCount int64 `json:"total_unloaded_count"`
	} `json:"classes"`
}

// ThreadPoolStats 线程池统计
type ThreadPoolStats struct {
	Threads   int64 `json:"threads"`
	Queue     int64 `json:"queue"`
	Active    int64 `json:"active"`
	Rejected  int64 `json:"rejected"`
	Largest   int64 `json:"largest"`
	Completed int64 `json:"completed"`
}

// FSStats 文件系统统计（完整版）
type FSStats struct {
	Timestamp int64 `json:"timestamp"`
	Total     struct {
		TotalInBytes     int64 `json:"total_in_bytes"`
		FreeInBytes      int64 `json:"free_in_bytes"`
		AvailableInBytes int64 `json:"available_in_bytes"`
	} `json:"total"`
	Data []struct {
		Path             string `json:"path"`
		Mount            string `json:"mount"`
		Type             string `json:"type"`
		TotalInBytes     int64  `json:"total_in_bytes"`
		FreeInBytes      int64  `json:"free_in_bytes"`
		AvailableInBytes int64  `json:"available_in_bytes"`
	} `json:"data"`
	IOStats struct {
		Devices []struct {
			DeviceName      string `json:"device_name"`
			Operations      int64  `json:"operations"`
			ReadOperations  int64  `json:"read_operations"`
			WriteOperations int64  `json:"write_operations"`
			ReadKilobytes   int64  `json:"read_kilobytes"`
			WriteKilobytes  int64  `json:"write_kilobytes"`
		} `json:"devices"`
		Total struct {
			Operations      int64 `json:"operations"`
			ReadOperations  int64 `json:"read_operations"`
			WriteOperations int64 `json:"write_operations"`
			ReadKilobytes   int64 `json:"read_kilobytes"`
			WriteKilobytes  int64 `json:"write_kilobytes"`
		} `json:"total"`
	} `json:"io_stats"`
}

// TransportStats 传输层统计
type TransportStats struct {
	ServerOpen          int64 `json:"server_open"`
	RxCount             int64 `json:"rx_count"`
	RxSizeInBytes       int64 `json:"rx_size_in_bytes"`
	TxCount             int64 `json:"tx_count"`
	TxSizeInBytes       int64 `json:"tx_size_in_bytes"`
	InboundHandlingTimeInMillis  int64 `json:"inbound_handling_time_in_millis"`
	OutboundHandlingTimeInMillis int64 `json:"outbound_handling_time_in_millis"`
}

// HTTPStats HTTP 统计
type HTTPStats struct {
	CurrentOpen int64 `json:"current_open"`
	TotalOpened int64 `json:"total_opened"`
}

// BreakerStats 断路器统计
type BreakerStats struct {
	LimitSizeInBytes     int64   `json:"limit_size_in_bytes"`
	LimitSize            string  `json:"limit_size"`
	EstimatedSizeInBytes int64   `json:"estimated_size_in_bytes"`
	EstimatedSize        string  `json:"estimated_size"`
	Overhead             float64 `json:"overhead"`
	Tripped              int64   `json:"tripped"`
}

// ScriptStats 脚本统计
type ScriptStats struct {
	Compilations         int64 `json:"compilations"`
	CacheEvictions       int64 `json:"cache_evictions"`
	CompilationLimitTriggered int64 `json:"compilation_limit_triggered"`
}

// DiscoveryStats 发现统计
type DiscoveryStats struct {
	ClusterStateQueue struct {
		Total     int64 `json:"total"`
		Pending   int64 `json:"pending"`
		Committed int64 `json:"committed"`
	} `json:"cluster_state_queue"`
}

// IngestStats 摄取统计
type IngestStats struct {
	Total struct {
		Count        int64 `json:"count"`
		TimeInMillis int64 `json:"time_in_millis"`
		Current      int64 `json:"current"`
		Failed       int64 `json:"failed"`
	} `json:"total"`
	Pipelines map[string]struct {
		Count        int64 `json:"count"`
		TimeInMillis int64 `json:"time_in_millis"`
		Current      int64 `json:"current"`
		Failed       int64 `json:"failed"`
	} `json:"pipelines"`
}
